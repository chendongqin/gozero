package role

import (
	"net/http"

	"gozore-mall/app/cmd/admin/internal/logic/sys/role"
	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetSysRoleListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := role.NewGetSysRoleListLogic(r.Context(), svcCtx)
		resp, err := l.GetSysRoleList()
		if err != nil {
			httpx.Error(w, err)
			return
		}

		response.Response(w, resp, err)
	}
}
