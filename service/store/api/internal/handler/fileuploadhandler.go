package handler

import (
	"FileStore-System/common/response"
	"net/http"
	"strconv"

	"FileStore-System/service/store/api/internal/logic"
	"FileStore-System/service/store/api/internal/svc"
	"FileStore-System/service/store/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileUploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			httpx.Error(w, err)
			return
		}

		id := r.FormValue("parent_id")
		parentId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		l := logic.NewFileUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileUpload(&req, parentId, file, fileHeader)
		response.Response(w, resp, err)
	}
}
