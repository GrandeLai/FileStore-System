package logic

import (
	"FileStore-System/common/errorx"
	"FileStore-System/service/user/rpc/internal/svc"
	"FileStore-System/service/user/rpc/types/user"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddVolumeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddVolumeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddVolumeLogic {
	return &AddVolumeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddVolumeLogic) AddVolume(in *user.AddVolumeReq) (*user.AddVolumeResp, error) {
	res, err := l.svcCtx.UserModel.UpdateVolume(l.ctx, in.Id, in.Size)
	if err != nil {
		return nil, err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if num == 0 {
		return nil, errorx.NewCodeError(100, "文件过大！")
	}
	return &user.AddVolumeResp{}, nil
}
