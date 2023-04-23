package handler

import (
	"FileStore-System/common/response"
	"net/http"

	"FileStore-System/service/store/api/internal/logic"
	"FileStore-System/service/store/api/internal/svc"
	"FileStore-System/service/store/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FolderCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FolderCreateRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewFolderCreateLogic(r.Context(), svcCtx)
		resp, err := l.FolderCreate(&req)
		response.Response(w, resp, err)
	}
}
