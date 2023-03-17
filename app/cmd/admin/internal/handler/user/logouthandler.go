package user

import (
	"net/http"

	"gozore-mall/app/cmd/admin/internal/logic/user"
	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewLogoutLogic(r.Context(), svcCtx)
		err := l.Logout()
		if err != nil {
			httpx.Error(w, err)
			return
		}

		response.Response(w, nil, err)
	}
}
