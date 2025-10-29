package chat

import (
	"context"
	"encoding/json"
	"net/http"

	"jxzy/apis/api_chat/internal/svc"
	"jxzy/apis/api_chat/internal/types"
	"jxzy/bll/bll_context/bllcontextservice"
	"jxzy/common/logger"

	"github.com/zeromicro/go-zero/core/logx"
)

type StreamChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStreamChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StreamChatLogic {
	serviceLogger := logger.NewServiceLogger("chat-api").WithContext(ctx)

	return &StreamChatLogic{
		Logger: serviceLogger,
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StreamChatLogic) StreamChat(req *types.ChatRequest, w http.ResponseWriter) error {
	l.Logger.Infof("Starting StreamChat logic - UserID:%s SceneCode:%s SessionID:%s",
		req.UserId, req.SceneCode, req.SessionId)

	// 验证必填参数
	if req.Message == "" {
		l.Logger.Error("Message is required")
		http.Error(w, "message is required", http.StatusBadRequest)
		return http.ErrServerClosed
	}
	if req.SceneCode == "" {
		l.Logger.Error("SceneCode is required")
		http.Error(w, "scene_code is required", http.StatusBadRequest)
		return http.ErrServerClosed
	}

	l.Logger.Debug("Request validation passed")

	// 检查RPC服务是否可用
	if l.svcCtx.BllContextRpc == nil {
		l.Logger.Error("BLL Context RPC service not available")
		http.Error(w, "bll_context RPC service not available", http.StatusServiceUnavailable)
		return http.ErrServerClosed
	}

	// 构建RPC请求
	rpcReq := &bllcontextservice.ChatRequest{
		Message:   req.Message,
		SessionId: req.SessionId,
		SceneCode: req.SceneCode,
		UserId:    req.UserId,
	}

	l.Logger.Debug("Calling BLL Context RPC service")
	// 调用RPC服务
	stream, err := l.svcCtx.BllContextRpc.StreamChat(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to call bll_context StreamChat: %v", err)
		http.Error(w, "failed to start chat stream", http.StatusInternalServerError)
		return http.ErrServerClosed
	}

	l.Logger.Info("RPC stream established, starting to process responses")

	responseCount := 0
	// 处理流式响应
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				// 流结束
				l.Logger.Infof("Stream ended normally after %d responses", responseCount)
				break
			}
			l.Logger.Errorf("Failed to receive stream response: %v", err)
			http.Error(w, "failed to receive stream response", http.StatusInternalServerError)
			return http.ErrServerClosed
		}

		responseCount++
		l.Logger.Infof("Received response %d - Delta:%s Finished:%v", responseCount, resp.Delta, resp.Finished)

		// 转换为API响应格式
		apiResp := &types.StreamChatData{
			SessionId: resp.SessionId,
			SceneCode: resp.SceneCode,
			Delta:     resp.Delta,
			Finished:  resp.Finished,
		}

		// 如果有token使用信息，转换并添加
		if resp.Usage != nil {
			apiResp.Usage = &types.TokenUsage{
				PromptTokens: resp.Usage.PromptTokens,
				ReplyTokens:  resp.Usage.ReplyTokens,
				TotalTokens:  resp.Usage.TotalTokens,
			}
			l.Logger.Infof("Token usage - Prompt:%d Reply:%d Total:%d",
				resp.Usage.PromptTokens, resp.Usage.ReplyTokens, resp.Usage.TotalTokens)
		}

		// 发送SSE数据
		data, err := json.Marshal(apiResp)
		if err != nil {
			l.Logger.Errorf("Failed to marshal response: %v", err)
			continue
		}

		_, err = w.Write([]byte("data: " + string(data) + "\n\n"))
		if err != nil {
			l.Logger.Errorf("Failed to write response: %v", err)
			return err
		}

		// 刷新缓冲区
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}

		// 如果对话完成，结束流
		if resp.Finished {
			l.Logger.Info("Chat finished, ending stream")
			break
		}
	}

	l.Logger.Info("StreamChat logic completed successfully")
	return nil
}
