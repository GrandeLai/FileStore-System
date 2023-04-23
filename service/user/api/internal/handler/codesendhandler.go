package handler

import (
	"FileStore-System/common/response"
	"net/http"

	"FileStore-System/service/user/api/internal/logic"
	"FileStore-System/service/user/api/internal/svc"
	"FileStore-System/service/user/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CodeSendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CodeSendRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewCodeSendLogic(r.Context(), svcCtx)
		resp, err := l.CodeSend(&req)
		response.Response(w, resp, err)
	}
}
