// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	common "gozore-mall/app/cmd/admin/internal/handler/common"
	configdict "gozore-mall/app/cmd/admin/internal/handler/config/dict"
	loglogin "gozore-mall/app/cmd/admin/internal/handler/log/login"
	sysdept "gozore-mall/app/cmd/admin/internal/handler/sys/dept"
	sysjob "gozore-mall/app/cmd/admin/internal/handler/sys/job"
	sysmenu "gozore-mall/app/cmd/admin/internal/handler/sys/menu"
	sysprofession "gozore-mall/app/cmd/admin/internal/handler/sys/profession"
	sysrole "gozore-mall/app/cmd/admin/internal/handler/sys/role"
	sysstoremenu "gozore-mall/app/cmd/admin/internal/handler/sys/storemenu"
	sysuser "gozore-mall/app/cmd/admin/internal/handler/sys/user"
	user "gozore-mall/app/cmd/admin/internal/handler/user"
	"gozore-mall/app/cmd/admin/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: user.LoginHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/login/captcha",
				Handler: user.GetLoginCaptchaHandler(serverCtx),
			},
		},
		rest.WithPrefix("/admin/user"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.AdminActLog},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/logout",
					Handler: user.LogoutHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/info",
					Handler: user.GetUserInfoHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/permmenu",
					Handler: user.GetUserPermMenuHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/profile/info",
					Handler: user.GetUserProfileInfoHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/profile/update",
					Handler: user.UpdateUserProfileHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/password/update",
					Handler: user.UpdateUserPasswordHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/avatar/generate",
					Handler: user.GetGenerateAvatarHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/admin/user"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.PermMenuAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/list",
					Handler: sysmenu.GetSysPermMenuListHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/add",
					Handler: sysmenu.AddSysPermMenuHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/delete",
					Handler: sysmenu.DeleteSysPermMenuHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/update",
					Handler: sysmenu.UpdateSysPermMenuHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/admin/sys/perm/menu"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.PermMenuAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/list",
					Handler: sysstoremenu.GetStorePermMenuListHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/add",
					Handler: sysstoremenu.AddStorePermMenuHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/delete",
					Handler: sysstoremenu.DeleteStorePermMenuHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/update",
					Handler: sysstoremenu.UpdateStorePermMenuHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/admin/sys/perm/store/menu"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.PermMenuAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/list",
					Handler: sysrole.GetSysRoleListHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/add",
					Handler: sysrole.AddSysRoleHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/delete",
					Handler: sysrole.DeleteSysRoleHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/update",
					Handler: sysrole.UpdateSysRoleHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/admin/sys/role"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.PermMenuAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/list",
					Handler: sysdept.GetSysDeptListHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/add",
					Handler: sysdept.AddSysDeptHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/delete",
					Handler: sysdept.DeleteSysDeptHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/update",
					Handler: sysdept.UpdateSysDeptHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/admin/sys/dept"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.PermMenuAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/page",
					Handler: sysjob.GetSysJobPageHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/add",
					Handler: sysjob.AddSysJobHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/delete",
					Handler: sysjob.DeleteSysJobHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/update",
					Handler: sysjob.UpdateSysJobHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/admin/sys/job"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.PermMenuAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/page",
					Handler: sysprofession.GetSysProfessionPageHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/add",
					Handler: sysprofession.AddSysProfessionHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/delete",
					Handler: sysprofession.DeleteSysProfessionHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/update",
					Handler: sysprofession.UpdateSysProfessionHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/admin/sys/profession"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.PermMenuAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/page",
					Handler: sysuser.GetSysUserPageHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/add",
					Handler: sysuser.AddSysUserHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/delete",
					Handler: sysuser.DeleteSysUserHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/update",
					Handler: sysuser.UpdateSysUserHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/password/update",
					Handler: sysuser.UpdateSysUserPasswordHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/rdpj/info",
					Handler: sysuser.GetSysUserRdpjInfoHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/admin/sys/user"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.PermMenuAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/list",
					Handler: configdict.GetConfigDictListHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/data/page",
					Handler: configdict.GetConfigDictPageHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/add",
					Handler: configdict.AddConfigDictHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/delete",
					Handler: configdict.DeleteConfigDictHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/update",
					Handler: configdict.UpdateConfigDictHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/admin/config/dict"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.PermMenuAuth, serverCtx.AdminActLog},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/page",
					Handler: loglogin.GetLogLoginPageHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/admin/log/login"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.PermMenuAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/oss/token",
					Handler: common.GetOssTokenHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/admin/common"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.PermMenuAuth},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/oss/callback",
					Handler: common.OssCallbackHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/admin/common"),
	)
}
