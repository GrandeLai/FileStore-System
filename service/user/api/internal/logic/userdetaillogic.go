package logic

import (
	"FileStore-System/common/errorx"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"

	"FileStore-System/service/user/api/internal/svc"
	"FileStore-System/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDetailLogic {
	return &UserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDetailLogic) UserDetail(req *types.DetailRequest, Authorization string) (resp *types.DetailResponse, err error) {
	restClaims := make(jwt.MapClaims)
	_, err = jwt.ParseWithClaims(Authorization, restClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(l.svcCtx.Config.Auth.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	//userId := restClaims["userId"]
	userId, err := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId"))).Int64()
	if err != nil {
		return nil, err
	}
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		return nil, errorx.NewDefaultError("id输入有误！")
	}
	return &types.DetailResponse{
		Name:  userInfo.UserName,
		Email: userInfo.Email,
	}, nil
}
