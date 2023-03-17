package storemenu

import (
	"context"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddStorePermMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddStorePermMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddStorePermMenuLogic {
	return &AddStorePermMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddStorePermMenuLogic) AddStorePermMenu(req *types.AddStorePermMenuReq) error {
	// todo: add your logic here and delete this line

	return nil
}
