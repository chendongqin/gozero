package storemenu

import (
	"net/http"

	"gozore-mall/app/cmd/admin/internal/logic/sys/storemenu"
	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetStorePermMenuListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := storemenu.NewGetStorePermMenuListLogic(r.Context(), svcCtx)
		resp, err := l.GetStorePermMenuList()
		if err != nil {
			httpx.Error(w, err)
			return
		}

		response.Response(w, resp, err)
	}
}
