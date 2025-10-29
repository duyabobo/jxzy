package chat

import (
	"net/http"
	"time"

	"jxzy/apis/api_chat/internal/logic/chat"
	"jxzy/apis/api_chat/internal/svc"
	"jxzy/apis/api_chat/internal/types"
	"jxzy/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func StreamChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		logger := logx.WithContext(r.Context())

		// 生成请求ID
		requestID := utils.GenerateUUID()
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		// 记录请求开始
		logger.Infof("[REQUEST] ID:%s User:%s %s %s", requestID, userID, r.Method, r.URL.Path)

		// 解析请求
		var req types.ChatRequest
		if err := httpx.Parse(r, &req); err != nil {
			logger.Errorf("Failed to parse request: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		logger.Infof("Received chat request - UserID:%s SceneCode:%s SessionID:%s Message:%s",
			req.UserId, req.SceneCode, req.SessionId, req.Message)

		// 设置SSE headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

		logger.Debug("Starting stream chat logic")
		l := chat.NewStreamChatLogic(r.Context(), svcCtx)
		err := l.StreamChat(&req, w)

		duration := time.Since(startTime)

		if err != nil {
			logger.Errorf("Stream chat failed: %v", err)
			// 发送错误事件
			w.Write([]byte("event: error\ndata: " + err.Error() + "\n\n"))
			w.(http.Flusher).Flush()
			logger.Infof("[RESPONSE] ID:%s User:%s %s %s Status:%d Duration:%v",
				requestID, userID, r.Method, r.URL.Path, 500, duration)
		} else {
			logger.Info("Stream chat completed successfully")
			logger.Infof("[RESPONSE] ID:%s User:%s %s %s Status:%d Duration:%v",
				requestID, userID, r.Method, r.URL.Path, 200, duration)
		}
	}
}
