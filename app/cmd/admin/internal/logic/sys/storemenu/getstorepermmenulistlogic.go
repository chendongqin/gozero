package storemenu

import (
	"context"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStorePermMenuListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStorePermMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStorePermMenuListLogic {
	return &GetStorePermMenuListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetStorePermMenuListLogic) GetStorePermMenuList() (resp *types.StorePermMenuListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
