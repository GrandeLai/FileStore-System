package logic

import (
	"FileStore-System/service/user/rpc/internal/svc"
	"FileStore-System/service/user/rpc/types/user"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVolumeByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVolumeByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVolumeByUserIdLogic {
	return &GetVolumeByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVolumeByUserIdLogic) GetVolumeByUserId(in *user.VolumeDetailReq) (*user.VolumeDetailResp, error) {
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &user.VolumeDetailResp{
		NowVolume:   userInfo.NowVolume,
		TotalVolume: userInfo.TotalVolume,
	}, nil
}
