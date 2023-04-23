package handler

import (
	"FileStore-System/common/response"
	"net/http"

	"FileStore-System/service/share/api/internal/logic"
	"FileStore-System/service/share/api/internal/svc"
	"FileStore-System/service/share/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SaveFromShareHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SaveFromShareRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewSaveFromShareLogic(r.Context(), svcCtx)
		resp, err := l.SaveFromShare(&req)
		response.Response(w, resp, err)
	}
}
