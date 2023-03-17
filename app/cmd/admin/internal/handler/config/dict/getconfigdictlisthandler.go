package dict

import (
	"net/http"

	"gozore-mall/app/cmd/admin/internal/logic/config/dict"
	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetConfigDictListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := dict.NewGetConfigDictListLogic(r.Context(), svcCtx)
		resp, err := l.GetConfigDictList()
		if err != nil {
			httpx.Error(w, err)
			return
		}

		response.Response(w, resp, err)
	}
}
