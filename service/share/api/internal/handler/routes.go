// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"FileStore-System/service/share/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/share/detail",
				Handler: DetailShareHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/share/create",
				Handler: CreateShareHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/share/save",
				Handler: SaveFromShareHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}