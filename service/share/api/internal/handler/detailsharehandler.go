package handler

import (
	"FileStore-System/common/response"
	"net/http"

	"FileStore-System/service/share/api/internal/logic"
	"FileStore-System/service/share/api/internal/svc"
	"FileStore-System/service/share/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DetailShareHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShareDetailRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewDetailShareLogic(r.Context(), svcCtx)
		resp, err := l.DetailShare(&req)
		response.Response(w, resp, err)
	}
}
