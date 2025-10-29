package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	contextpb "jxzy/bll/bll_context/bll_context"
	"jxzy/bll/bll_context/internal/common"
	"jxzy/bll/bll_context/internal/model"
	"jxzy/bll/bll_context/internal/svc"
	"jxzy/bs/bs_llm/bs_llm"
	"jxzy/bs/bs_rag/bs_rag"
	consts "jxzy/common/const"
	"jxzy/common/logger"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"golang.org/x/sync/errgroup"
)

type StreamChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	embeddingService *common.EmbeddingService
}

func NewStreamChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StreamChatLogic {
	// 使用自定义的 ServiceLogger，在日志中显示服务名
	serviceLogger := logger.NewServiceLogger("bll-context").WithContext(ctx)

	return &StreamChatLogic{
		ctx:              ctx,
		svcCtx:           svcCtx,
		Logger:           serviceLogger,
		embeddingService: common.NewEmbeddingService(svcCtx.Config.Bailian.APIKey),
	}
}

// 流式聊天
func (l *StreamChatLogic) StreamChat(in *contextpb.ChatRequest, stream contextpb.BllContextService_StreamChatServer) error {
	l.Logger.Infof("StreamChat called with scene_code: %s, session_id: %s, user_id: %s", in.SceneCode, in.SessionId, in.UserId)

	// 1. 验证scene_code参数
	if in.SceneCode == "" {
		l.Logger.Error("scene_code is required")
		return fmt.Errorf("scene_code is required")
	}

	// 2. 获取或创建session
	l.Logger.Debug("Starting to get or create session")
	sessionId, err := l.getOrCreateSession(in)
	if err != nil {
		l.Logger.Errorf("Failed to get or create session: %v", err)
		return fmt.Errorf("failed to get or create session: %w", err)
	}

	l.Logger.Infof("Using session_id: %s for scene_code: %s", sessionId, in.SceneCode)

	// 3. 构建LLM请求
	messages := l.getLLMInputMessages(in)
	l.Logger.Debug("Building LLM request")
	llmReq := &bs_llm.LLMRequest{
		SceneCode:   in.SceneCode,
		Messages:    messages,
		ExtraParams: make(map[string]string),
		UserId:      in.UserId,
	}

	l.Logger.Infof("LLM request built - SceneCode: %s, Message: %s", llmReq.SceneCode, in.Message)

	// 4. 调用LLM的StreamLLM RPC
	l.Logger.Debug("Calling LLM StreamLLM RPC")
	llmStream, err := l.svcCtx.LLMRpc.StreamLLM(l.ctx, llmReq)
	if err != nil {
		l.Logger.Errorf("Failed to call LLM StreamLLM: %v", err)
		return fmt.Errorf("failed to call LLM StreamLLM: %w", err)
	}

	l.Logger.Info("LLM stream established, starting to process responses")

	responseCount := 0
	var finalUsage *bs_llm.LLMUsage

	// 5. 处理流式响应并转发
	for {
		llmResp, err := llmStream.Recv()
		if err != nil {
			if err == io.EOF {
				l.Logger.Infof("LLM stream ended normally after %d responses", responseCount)
				break
			}
			l.Logger.Errorf("Failed to receive from LLM stream: %v", err)
			return fmt.Errorf("failed to receive from LLM stream: %w", err)
		}

		responseCount++
		l.Logger.Infof("Received LLM response %d - Delta: %s, Finished: %v", responseCount, llmResp.Delta, llmResp.Finished)

		// 构建Context层的流式响应
		contextResp := &contextpb.StreamChatResponse{
			SessionId: sessionId,
			SceneCode: in.SceneCode,
			Delta:     llmResp.Delta,
			Finished:  llmResp.Finished,
		}

		// 如果有token使用信息，转换并添加
		if llmResp.Usage != nil {
			contextResp.Usage = &contextpb.TokenUsage{
				PromptTokens: llmResp.Usage.PromptTokens,
				ReplyTokens:  llmResp.Usage.CompletionTokens,
				TotalTokens:  llmResp.Usage.TotalTokens,
			}
			finalUsage = llmResp.Usage
			l.Logger.Infof("Token usage - Prompt: %d, Reply: %d, Total: %d",
				llmResp.Usage.PromptTokens, llmResp.Usage.CompletionTokens, llmResp.Usage.TotalTokens)
		}

		// 如果对话完成，更新会话的最后交互时间并记录chat_session_qas
		if llmResp.Finished {
			l.Logger.Debug("Chat finished, updating session interact time and recording chat_session_qas")
			l.updateSessionInteractTime(sessionId)

			// 记录chat_session_qas
			if err := l.recordChatSessionQas(sessionId, in.Message, finalUsage); err != nil {
				l.Logger.Errorf("Failed to record chat_session_qas: %v", err)
				// 不中断流程，只记录错误
			}

			l.Logger.Infof("Stream completed for session: %s, scene: %s, total responses: %d", sessionId, in.SceneCode, responseCount)
		}

		// 发送响应
		if err := stream.Send(contextResp); err != nil {
			l.Logger.Errorf("Failed to send stream response: %v", err)
			return fmt.Errorf("failed to send stream response: %w", err)
		}
	}

	l.Logger.Info("StreamChat logic completed successfully")
	return nil
}

func (l *StreamChatLogic) getLLMInputMessages(in *contextpb.ChatRequest) []*bs_llm.ChatMessage {
	messages := make([]*bs_llm.ChatMessage, 0)

	// 1. 获取system prompt
	systemPrompt := l.getSystemPrompt(in)
	// todo 这里要sse输出 systemPrompt 的内容
	messages = append(messages, &bs_llm.ChatMessage{
		Role:    "system",
		Content: systemPrompt,
	})

	// 2. 获取记忆
	memory := l.getMemory(in)
	// todo 这里要sse输出 memory 的内容
	messages = append(messages, &bs_llm.ChatMessage{
		Role:    "system",
		Content: memory,
	})

	// 3. 获取rag
	rag := l.getRag(in)
	// todo 这里要sse输出 rag 的内容
	messages = append(messages, &bs_llm.ChatMessage{
		Role:    "system",
		Content: rag,
	})

	// 4. 获取用户消息
	messages = append(messages, &bs_llm.ChatMessage{
		Role:    "user",
		Content: in.Message,
	})

	return messages
}

func (l *StreamChatLogic) getSystemPrompt(in *contextpb.ChatRequest) string {
	// 获取system prompt
	return ""
}

func (l *StreamChatLogic) getMemory(in *contextpb.ChatRequest) string {
	// 获取记忆
	return ""
}

func (l *StreamChatLogic) getRag(in *contextpb.ChatRequest) string {
	l.Logger.Info("Starting RAG processing")

	// 1. 非流式调用一次LLM，解析用户输入，扩展和拆解出可用于RAG检索的关键句
	// 内置system prompt，解析出关键句列表
	systemPrompt := `你是一个智能助手，负责从用户输入中提取和扩展可用于检索增强生成(RAG)的关键句。

你的任务是：
1. 分析用户的问题或需求
2. 提取核心概念和实体
3. 对输入进行语义扩展，生成相关的查询句
4. 将复杂问题拆解为多个子问题
5. 生成不同角度的检索句

请以JSON格式返回，格式为：
{
  "key_sentences": [
    "原始问题的核心表述",
    "扩展的相关概念查询",
    "拆解的子问题1",
    "拆解的子问题2",
    "同义词或近义词查询",
    "技术术语的完整表述"
  ]
}

示例：
用户输入："如何优化数据库性能？"
返回：
{
  "key_sentences": [
    "数据库性能优化",
    "SQL查询优化",
    "数据库索引优化",
    "数据库配置调优",
    "数据库性能监控",
    "数据库缓存策略",
    "数据库连接池优化",
    "数据库存储优化"
  ]
}

只返回JSON，不要其他解释。`

	userMessage := &bs_llm.ChatMessage{
		Role:    "user",
		Content: in.Message,
	}

	llmReq := &bs_llm.LLMRequest{
		SceneCode: "rag-sentence-extraction", // 使用专门的关键句提取场景码
		Messages: []*bs_llm.ChatMessage{
			{Role: "system", Content: systemPrompt},
			userMessage,
		},
		ExtraParams: make(map[string]string),
		UserId:      in.UserId,
	}

	l.Logger.Debug("Calling LLM for key sentence extraction")
	llmResp, err := l.svcCtx.LLMRpc.LLM(l.ctx, llmReq)
	if err != nil {
		l.Logger.Errorf("Failed to call LLM for key sentence extraction: %v", err)
		return ""
	}

	l.Logger.Infof("LLM response for key sentence extraction: %s", llmResp.Completion)

	// 2. 解析LLM响应，提取关键句列表
	keySentences, err := l.parseKeySentencesFromLLMResponse(llmResp.Completion)
	if err != nil {
		l.Logger.Errorf("Failed to parse key sentences from LLM response: %v", err)
		return ""
	}

	if len(keySentences) == 0 {
		l.Logger.Info("No key sentences extracted, skipping RAG")
		return ""
	}

	l.Logger.Infof("Extracted key sentences: %v", keySentences)

	// 3. 使用协程并发调用RAG，获取结果列表
	ragResults := l.searchRAGConcurrently(keySentences, in.UserId)

	// 4. 拼接RAG结果为system prompt的一部分
	if len(ragResults) == 0 {
		l.Logger.Info("No relevant RAG results found")
		return ""
	}

	result := strings.Join(ragResults, "\n\n")
	l.Logger.Infof("RAG processing completed, total results: %d", len(ragResults))

	return result
}

// getOrCreateSession 获取或创建session
func (l *StreamChatLogic) getOrCreateSession(req *contextpb.ChatRequest) (string, error) {
	l.Logger.Debug("getOrCreateSession called")

	// 检查session model是否已初始化
	if l.svcCtx.ChatSessionModel == nil {
		l.Logger.Error("chat session model not initialized")
		return "", fmt.Errorf("chat session model not initialized")
	}

	// 如果提供了session_id，先尝试获取
	if req.SessionId != "" {
		l.Logger.Infof("Trying to find existing session: %s", req.SessionId)
		session, err := l.svcCtx.ChatSessionModel.FindOne(l.ctx, l.parseSessionId(req.SessionId))
		if err != nil && err != sqlc.ErrNotFound {
			l.Logger.Errorf("Failed to find session %s: %v", req.SessionId, err)
			return "", fmt.Errorf("failed to find session: %w", err)
		}
		if session != nil {
			// 验证scene_code是否匹配
			if session.SceneCode != req.SceneCode {
				l.Logger.Errorf("Session scene_code mismatch - expected: %s, got: %s", session.SceneCode, req.SceneCode)
				return "", fmt.Errorf("session scene_code mismatch: expected %s, got %s", session.SceneCode, req.SceneCode)
			}
			l.Logger.Infof("Found existing session: %s", req.SessionId)
			return req.SessionId, nil
		}
		l.Logger.Infof("Session %s not found, will create new session", req.SessionId)
	}

	// 创建新session
	l.Logger.Debug("Creating new session")
	return l.createNewSession(req)
}

// createNewSession 创建新session
func (l *StreamChatLogic) createNewSession(req *contextpb.ChatRequest) (string, error) {
	l.Logger.Debug("createNewSession called")

	if l.svcCtx.ChatSessionModel == nil {
		l.Logger.Error("chat session model not initialized")
		return "", fmt.Errorf("chat session model not initialized")
	}

	userId := req.UserId
	if userId == "" {
		userId = "anonymous"
		l.Logger.Debug("Using anonymous user ID")
	}

	sessionName := fmt.Sprintf("Chat-%s", time.Now().Format("2006-01-02 15:04:05"))
	l.Logger.Infof("Creating session with name: %s", sessionName)

	session := &model.ChatSession{
		Name:             sessionName,
		SceneCode:        req.SceneCode,
		UserId:           userId,
		IsActive:         1,
		LastInteractTime: time.Now(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	result, err := l.svcCtx.ChatSessionModel.Insert(l.ctx, session)
	if err != nil {
		l.Logger.Errorf("Failed to create session: %v", err)
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	sessionId, err := result.LastInsertId()
	if err != nil {
		l.Logger.Errorf("Failed to get session id: %v", err)
		return "", fmt.Errorf("failed to get session id: %w", err)
	}

	sessionIdStr := strconv.FormatInt(sessionId, 10)
	l.Logger.Infof("Created new session: %s for user: %s, scene: %s", sessionIdStr, userId, req.SceneCode)

	return sessionIdStr, nil
}

// updateSessionInteractTime 更新会话的最后交互时间
func (l *StreamChatLogic) updateSessionInteractTime(sessionId string) {
	// 检查context是否已取消
	select {
	case <-l.ctx.Done():
		l.Logger.Infof("Context canceled, skipping session update")
		return
	default:
	}

	if l.svcCtx.ChatSessionModel == nil {
		l.Logger.Error("chat session model not initialized")
		return
	}

	id := l.parseSessionId(sessionId)
	if id <= 0 {
		l.Logger.Errorf("Invalid session_id: %s", sessionId)
		return
	}

	session, err := l.svcCtx.ChatSessionModel.FindOne(l.ctx, id)
	if err != nil {
		l.Logger.Errorf("Failed to find session for update: %v", err)
		return
	}

	session.LastInteractTime = time.Now()
	session.UpdatedAt = time.Now()

	err = l.svcCtx.ChatSessionModel.Update(l.ctx, session)
	if err != nil {
		l.Logger.Errorf("Failed to update session interact time: %v", err)
	} else {
		l.Logger.Infof("Updated session %s interact time", sessionId)
	}
}

// parseSessionId 解析session_id字符串为int64
func (l *StreamChatLogic) parseSessionId(sessionId string) int64 {
	id, err := strconv.ParseInt(sessionId, 10, 64)
	if err != nil {
		l.Logger.Errorf("Failed to parse session_id %s: %v", sessionId, err)
		return 0
	}
	return id
}

// parseKeySentencesFromLLMResponse 解析LLM响应中的关键句列表
func (l *StreamChatLogic) parseKeySentencesFromLLMResponse(llmResponse string) ([]string, error) {
	// 尝试解析JSON响应
	var response struct {
		KeySentences []string `json:"key_sentences"`
	}

	// 清理响应文本，移除可能的markdown格式
	cleanResponse := strings.TrimSpace(llmResponse)
	cleanResponse = strings.TrimPrefix(cleanResponse, "```json")
	cleanResponse = strings.TrimSuffix(cleanResponse, "```")

	if err := json.Unmarshal([]byte(cleanResponse), &response); err != nil {
		l.Logger.Infof("Failed to parse JSON response: %v, trying fallback parsing", err)
		// 回退解析：简单提取引号内的内容
		return l.extractKeySentencesFallback(cleanResponse)
	}

	return response.KeySentences, nil
}

// extractKeySentencesFallback 回退方法：从文本中提取关键句
func (l *StreamChatLogic) extractKeySentencesFallback(text string) ([]string, error) {
	var keySentences []string

	// 简单的启发式提取：查找引号包围的内容
	start := -1
	for i, r := range text {
		if r == '"' || r == '\'' {
			if start == -1 {
				start = i + 1
			} else {
				sentence := strings.TrimSpace(text[start:i])
				if len(sentence) > 2 { // 过滤太短的句子
					keySentences = append(keySentences, sentence)
				}
				start = -1
			}
		}
	}

	// 如果没有找到引号包围的内容，尝试按逗号分割
	if len(keySentences) == 0 {
		parts := strings.Split(text, ",")
		for _, part := range parts {
			sentence := strings.TrimSpace(part)
			sentence = strings.Trim(sentence, "\"'")
			if len(sentence) > 2 {
				keySentences = append(keySentences, sentence)
			}
		}
	}

	l.Logger.Infof("Fallback parsing extracted key sentences: %v", keySentences)
	return keySentences, nil
}

// generateQueryVector 生成查询向量
// 使用阿里云百炼 Embedding API 生成真实的向量
func (l *StreamChatLogic) generateQueryVector(sentence string) []float32 {
	// 调用 EmbeddingService 生成向量
	vector, err := l.embeddingService.GenerateEmbedding(sentence)
	if err != nil {
		l.Logger.Errorf("Failed to generate embedding for sentence '%s': %v", sentence, err)
		// 返回空向量，让调用方处理错误
		return make([]float32, consts.DashVectorDefaultDimension) // DefaultEmbeddingDimension
	}

	l.Logger.Debugf("Generated embedding vector for sentence '%s': length=%d", sentence, len(vector))
	return vector
}

// searchRAGConcurrently 并发搜索RAG
func (l *StreamChatLogic) searchRAGConcurrently(keySentences []string, userId string) []string {
	var ragResults []string
	var mu sync.Mutex // 保护 ragResults 的并发访问

	// 创建 errgroup 用于并发处理
	g, ctx := errgroup.WithContext(l.ctx)

	// 为每个关键句启动一个协程
	for _, sentence := range keySentences {
		sentence := sentence // 创建局部变量避免闭包问题

		g.Go(func() error {
			// 检查上下文是否已取消
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			l.Logger.Infof("Searching RAG for key sentence: %s", sentence)

			// 构建向量查询请求
			ragReq := &bs_rag.VectorSearchRequest{
				QueryVector:    l.generateQueryVector(sentence),
				TopK:           3,                            // 返回前3个最相似的结果
				MinScore:       0.5,                          // 最小相似度阈值
				CollectionName: consts.DefaultCollectionName, // 集合名称
				Filters:        make(map[string]string),
				UserId:         userId,
			}

			// 调用RAG服务
			ragResp, err := l.svcCtx.RagRpc.VectorSearch(ctx, ragReq)
			if err != nil {
				l.Logger.Errorf("Failed to search RAG for key sentence %s: %v", sentence, err)
				return nil // 返回 nil 而不是错误，避免中断其他协程
			}

			l.Logger.Infof("RAG search returned %d results for key sentence: %s", len(ragResp.Results), sentence)

			// 处理RAG结果
			var localResults []string
			for _, result := range ragResp.Results {
				if result.Score >= 0.5 { // 只使用相似度足够高的结果
					localResults = append(localResults, fmt.Sprintf("相关内容 (相似度: %.2f): %s", result.Score, result.Content))
				}
			}

			// 线程安全地添加到全局结果
			if len(localResults) > 0 {
				mu.Lock()
				ragResults = append(ragResults, localResults...)
				mu.Unlock()
			}

			return nil
		})
	}

	// 等待所有协程完成
	if err := g.Wait(); err != nil {
		l.Logger.Errorf("Error in concurrent RAG search: %v", err)
	}

	l.Logger.Infof("Concurrent RAG search completed, total results: %d", len(ragResults))
	return ragResults
}

// recordChatSessionQas 记录聊天会话问答记录
func (l *StreamChatLogic) recordChatSessionQas(sessionId string, userMessage string, usage *bs_llm.LLMUsage) error {
	sessionIdInt := l.parseSessionId(sessionId)
	if sessionIdInt <= 0 {
		return fmt.Errorf("invalid session_id: %s", sessionId)
	}

	// 创建chat_session_qas记录
	chatSessionQas := &model.ChatSessionQas{
		SessionId:       sessionIdInt,
		LlmCompletionId: 0,
		SequenceNum:     time.Now().Unix(),
		CreatedAt:       time.Now(),
	}

	_, err := l.svcCtx.ChatSessionQasModel.Insert(l.ctx, chatSessionQas)
	if err != nil {
		l.Logger.Errorf("Failed to insert chat_session_qas record: %v", err)
		return err
	}

	l.Logger.Infof("Recorded chat_session_qas for session: %s", sessionId)
	return nil
}
