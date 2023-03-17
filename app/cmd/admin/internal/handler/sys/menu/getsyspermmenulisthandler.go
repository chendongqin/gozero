package menu

import (
	"net/http"

	"gozore-mall/app/cmd/admin/internal/logic/sys/menu"
	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetSysPermMenuListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := menu.NewGetSysPermMenuListLogic(r.Context(), svcCtx)
		resp, err := l.GetSysPermMenuList()
		if err != nil {
			httpx.Error(w, err)
			return
		}

		response.Response(w, resp, err)
	}
}
