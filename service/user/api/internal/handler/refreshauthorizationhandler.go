package handler

import (
	"FileStore-System/common/response"
	"net/http"

	"FileStore-System/service/user/api/internal/logic"
	"FileStore-System/service/user/api/internal/svc"
	"FileStore-System/service/user/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RefreshAuthorizationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefreshAuthRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewRefreshAuthorizationLogic(r.Context(), svcCtx)
		resp, err := l.RefreshAuthorization(&req, r.Header.Get("Authorization"))
		response.Response(w, resp, err)
	}
}
