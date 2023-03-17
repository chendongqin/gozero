package storemenu

import (
	"context"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStorePermMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateStorePermMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStorePermMenuLogic {
	return &UpdateStorePermMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateStorePermMenuLogic) UpdateStorePermMenu(req *types.UpdateStorePermMenuReq) error {
	// todo: add your logic here and delete this line

	return nil
}
