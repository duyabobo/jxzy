package knowledge

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"jxzy/apis/api_knowledge/internal/logic/knowledge"
	"jxzy/apis/api_knowledge/internal/svc"
	"jxzy/apis/api_knowledge/internal/types"
)

func DeleteKnowledgeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteKnowledgeRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := knowledge.NewDeleteKnowledgeLogic(r.Context(), svcCtx)
		resp, err := l.DeleteKnowledge(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
