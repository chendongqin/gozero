package dept

import (
	"net/http"

	"gozore-mall/app/cmd/admin/internal/logic/sys/dept"
	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetSysDeptListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := dept.NewGetSysDeptListLogic(r.Context(), svcCtx)
		resp, err := l.GetSysDeptList()
		if err != nil {
			httpx.Error(w, err)
			return
		}

		response.Response(w, resp, err)
	}
}
