package knowledge

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"jxzy/apis/api_knowledge/internal/logic/knowledge"
	"jxzy/apis/api_knowledge/internal/svc"
	"jxzy/apis/api_knowledge/internal/types"
)

func AddKnowledgeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddKnowledgeRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := knowledge.NewAddKnowledgeLogic(r.Context(), svcCtx)
		resp, err := l.AddKnowledge(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
