package storemenu

import (
	"context"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteStorePermMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteStorePermMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteStorePermMenuLogic {
	return &DeleteStorePermMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteStorePermMenuLogic) DeleteStorePermMenu(req *types.DeleteStorePermMenuReq) error {
	// todo: add your logic here and delete this line

	return nil
}
