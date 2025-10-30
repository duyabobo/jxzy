package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"

	knowledgepb "jxzy/bll/bll_knowledge/bll_knowledge"
	"jxzy/bll/bll_knowledge/internal/model"
	"jxzy/bll/bll_knowledge/internal/svc"
	bsllm "jxzy/bs/bs_llm/bs_llm"
	"jxzy/bs/bs_rag/bs_rag"
	consts "jxzy/common/const"
	"jxzy/common/logger"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddVectorKnowledgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddVectorKnowledgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddVectorKnowledgeLogic {
	// 使用自定义的 ServiceLogger，在日志中显示服务名
	serviceLogger := logger.NewServiceLogger("bll-knowledge").WithContext(ctx)

	return &AddVectorKnowledgeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: serviceLogger,
	}
}

// AddVectorKnowledge 添加知识库到向量数据库
func (l *AddVectorKnowledgeLogic) AddVectorKnowledge(in *knowledgepb.AddVectorKnowledgeRequest) (*knowledgepb.AddVectorKnowledgeResponse, error) {
	l.Logger.Infof("AddVectorKnowledge called with file_url: %s, user_id: %s",
		in.FileUrl, in.UserId)

	// 1. 验证输入参数
	if err := l.validateInput(in); err != nil {
		l.Logger.Errorf("Input validation failed: %v", err)
		return &knowledgepb.AddVectorKnowledgeResponse{
			Success: false,
			Message: fmt.Sprintf("输入参数验证失败: %v", err),
		}, nil
	}

	// 2. 下载文件内容
	fileBytes, fileName, fileType, fileSize, err := l.downloadFile(in.FileUrl)
	if err != nil {
		l.Logger.Errorf("Failed to download file: %v", err)
		return &knowledgepb.AddVectorKnowledgeResponse{Success: false, Message: fmt.Sprintf("下载文件失败: %v", err)}, nil
	}

	// 3. 计算文件MD5并进行去重
	fileMd5 := l.md5Bytes(fileBytes)
	l.Logger.Infof("Computed file md5: %s (name=%s, size=%d, type=%s)", fileMd5, fileName, fileSize, fileType)

	// 3.1 查询是否已存在
	existing, _ := l.svcCtx.KnowledgeFileModel.FindOneByMd5(l.ctx, fileMd5)
	if existing != nil {
		l.Logger.Infof("File already processed, knowledge_file.id=%d", existing.Id)
		return &knowledgepb.AddVectorKnowledgeResponse{VectorId: fileMd5, Success: true, Message: "文件已存在，跳过处理"}, nil
	}

	// 4. 入库 knowledge_file
	kf := &model.KnowledgeFile{
		OssPath:  in.FileUrl,
		FileName: fileName,
		FileSize: fileSize,
		FileType: fileType,
		FileMd5:  fileMd5,
	}
	res, err := l.svcCtx.KnowledgeFileModel.Insert(l.ctx, kf)
	if err != nil {
		l.Logger.Errorf("Insert knowledge_file failed: %v", err)
		return &knowledgepb.AddVectorKnowledgeResponse{Success: false, Message: fmt.Sprintf("保存文件记录失败: %v", err)}, nil
	}
	fileId, _ := res.LastInsertId()

	// 5. 利用LLM进行语义段拆分
	segments, err := l.segmentTextWithLLM(string(fileBytes), in.UserId)
	if err != nil {
		l.Logger.Errorf("LLM segmentation failed: %v", err)
		return &knowledgepb.AddVectorKnowledgeResponse{Success: false, Message: fmt.Sprintf("LLM拆分失败: %v", err)}, nil
	}

	// 6. 存储语义段，并为每个语义段生成摘要句后存储与入RAG
	var documents []*bs_rag.VectorDocument
	for _, seg := range segments {
		segMd5 := l.md5String(seg)
		segRec := &model.KnowledgeSegment{KnowledgeFileId: fileId, SegmentText: seg, SegmentMd5: segMd5}
		segRes, err := l.svcCtx.KnowledgeSegmentModel.Insert(l.ctx, segRec)
		if err != nil {
			l.Logger.Errorf("Insert knowledge_segment failed: %v", err)
			continue
		}
		segId, _ := segRes.LastInsertId()

		// 6.1 为语义段生成摘要句
		summary, err := l.summarizeSegmentWithLLM(seg, in.UserId)
		if err != nil || strings.TrimSpace(summary) == "" {
			l.Logger.Errorf("LLM summary failed for segment %d: %v", segId, err)
			continue
		}
		sumMd5 := l.md5String(summary)
		sumRec := &model.KnowledgeSummarySentence{
			KnowledgeFileId:     fileId,
			KnowledgeSegmentId:  segId,
			SummarySentenceText: summary,
			SummarySentenceMd5:  sumMd5,
		}
		sumRes, err := l.svcCtx.KnowledgeSummarySentenceModel.Insert(l.ctx, sumRec)
		if err != nil {
			l.Logger.Errorf("Insert knowledge_summary_sentence failed: %v", err)
			continue
		}
		summaryId, _ := sumRes.LastInsertId()

		// 6.2 构建RAG文档，Id使用摘要句入库id，内容使用摘要句
		documents = append(documents, &bs_rag.VectorDocument{
			Id:   fmt.Sprintf("%d", summaryId),
			Text: summary,
			Metadata: map[string]string{
				"knowledge_file_id":    fmt.Sprintf("%d", fileId),
				"knowledge_segment_id": fmt.Sprintf("%d", segId),
				"user_id":              in.UserId,
			},
			Content: summary,
		})
	}

	if len(documents) == 0 {
		return &knowledgepb.AddVectorKnowledgeResponse{VectorId: fileMd5, Success: false, Message: "没有可插入的摘要文档"}, nil
	}

	// 7. 调用RAG服务插入向量（使用摘要句作为content/text）
	if l.svcCtx.RagRpc == nil {
		l.Logger.Error("RAG service is not available")
		return &knowledgepb.AddVectorKnowledgeResponse{VectorId: fileMd5, Success: false, Message: "RAG服务不可用"}, nil
	}

	ragReq := &bs_rag.VectorInsertRequest{CollectionName: consts.DefaultCollectionName, Documents: documents, UserId: in.UserId}
	ragResp, err := l.svcCtx.RagRpc.VectorInsert(l.ctx, ragReq)
	if err != nil {
		l.Logger.Errorf("Failed to insert vectors to RAG service: %v", err)
		return &knowledgepb.AddVectorKnowledgeResponse{VectorId: fileMd5, Success: false, Message: fmt.Sprintf("插入RAG失败: %v", err)}, nil
	}
	l.Logger.Infof("Inserted %d documents to RAG", ragResp.InsertedCount)

	return &knowledgepb.AddVectorKnowledgeResponse{VectorId: fileMd5, Success: true, Message: "知识库添加成功"}, nil
}

// validateInput 验证输入参数
func (l *AddVectorKnowledgeLogic) validateInput(in *knowledgepb.AddVectorKnowledgeRequest) error {
	if strings.TrimSpace(in.FileUrl) == "" {
		return fmt.Errorf("file_url不能为空")
	}
	if strings.TrimSpace(in.UserId) == "" {
		return fmt.Errorf("user_id不能为空")
	}
	return nil
}

// md5String 计算字符串MD5
func (l *AddVectorKnowledgeLogic) md5String(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

// md5Bytes 计算字节MD5
func (l *AddVectorKnowledgeLogic) md5Bytes(b []byte) string {
	sum := md5.Sum(b)
	return hex.EncodeToString(sum[:])
}

// downloadFile 下载远端文件，返回内容及基础元信息
func (l *AddVectorKnowledgeLogic) downloadFile(rawUrl string) ([]byte, string, string, int64, error) {
	req, err := http.NewRequestWithContext(l.ctx, http.MethodGet, rawUrl, nil)
	if err != nil {
		return nil, "", "", 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", "", 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", "", 0, fmt.Errorf("http status: %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", "", 0, err
	}
	// 解析文件名
	parsed, _ := url.Parse(rawUrl)
	base := path.Base(parsed.Path)
	// 猜测类型
	ctype := resp.Header.Get("Content-Type")
	if ctype == "" {
		ctype = mime.TypeByExtension(path.Ext(base))
	}
	size := int64(len(data))
	return data, base, ctype, size, nil
}

// segmentTextWithLLM 使用LLM将文本拆分为语义段
func (l *AddVectorKnowledgeLogic) segmentTextWithLLM(text string, userId string) ([]string, error) {
	if l.svcCtx.LlmRpc == nil {
		return []string{strings.TrimSpace(text)}, nil
	}
	req := &bsllm.LLMRequest{
		SceneCode: "knowledge_segmentation",
		UserId:    userId,
		Messages: []*bsllm.ChatMessage{
			{Role: "system", Content: "请将以下内容按语义进行合理拆分，返回为每段一行。不要编号。"},
			{Role: "user", Content: text},
		},
	}
	resp, err := l.svcCtx.LlmRpc.LLM(l.ctx, req)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(resp.Completion, "\n")
	cleaned := make([]string, 0, len(lines))
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln != "" {
			cleaned = append(cleaned, ln)
		}
	}
	if len(cleaned) == 0 {
		cleaned = []string{strings.TrimSpace(text)}
	}
	return cleaned, nil
}

// summarizeSegmentWithLLM 使用LLM为语义段生成摘要句
func (l *AddVectorKnowledgeLogic) summarizeSegmentWithLLM(segment string, userId string) (string, error) {
	if l.svcCtx.LlmRpc == nil {
		return segment, nil
	}
	req := &bsllm.LLMRequest{
		SceneCode: "knowledge_segment_summary",
		UserId:    userId,
		Messages: []*bsllm.ChatMessage{
			{Role: "system", Content: "请为以下文本生成一句话摘要，简洁、完整且可用于检索。直接返回摘要句。"},
			{Role: "user", Content: segment},
		},
	}
	resp, err := l.svcCtx.LlmRpc.LLM(l.ctx, req)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp.Completion), nil
}
