package user

import (
	"context"
	"encoding/json"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/service/gormpool"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/app/model"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"
	"gozore-mall/common/utils"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddSysUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSysUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSysUserLogic {
	return &AddSysUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSysUserLogic) AddSysUser(req *types.AddSysUserReq) error {
	user := &sysadmin.SysUser{}
	err := l.svcCtx.GormQuery.FindOneBy(l.ctx, user, "account", req.Account)
	if err == model.ErrNotFound {
		currentUserId := utils.GetAdminUserId(l.ctx)
		var currentUserRoleIds []int64
		var roleIds []int64
		if currentUserId == globalkey.SysSuperUserId {
			sysRoleList := make([]sysadmin.SysRole, 0)
			_ = l.svcCtx.GormQuery.QueryFindAll(l.ctx, &sysRoleList, gormpool.Conditions{
				OrderBy: "order_num DESC",
			})
			for _, role := range sysRoleList {
				currentUserRoleIds = append(currentUserRoleIds, role.Id)
				roleIds = append(roleIds, role.Id)
			}

		} else {
			currentUser := &sysadmin.SysUser{}
			_ = l.svcCtx.GormQuery.FindOne(l.ctx, currentUser, currentUserId)
			err := json.Unmarshal([]byte(currentUser.RoleIds), &currentUserRoleIds)
			if err != nil {
				return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
			}

			roleIds = append(roleIds, currentUserRoleIds...)
		}

		for _, id := range req.RoleIds {
			if !utils.ArrayContainValue(roleIds, id) {
				return errorx.NewDefaultError(errorx.AssigningRolesErrorCode)
			}
		}
		sysDept := &sysadmin.SysDept{}
		err = l.svcCtx.GormQuery.FindOne(l.ctx, sysDept, req.DeptId)
		if err != nil {
			return errorx.NewDefaultError(errorx.DeptIdErrorCode)
		}
		sysProfession := &sysadmin.SysProfession{}
		err = l.svcCtx.GormQuery.FindOne(l.ctx, sysProfession, req.ProfessionId)
		if err != nil {
			return errorx.NewDefaultError(errorx.ProfessionIdErrorCode)
		}

		sysJob := &sysadmin.SysJob{}
		err = l.svcCtx.GormQuery.FindOne(l.ctx, sysJob, req.JobId)
		if err != nil {
			return errorx.NewDefaultError(errorx.JobIdErrorCode)
		}

		var sysUser = new(sysadmin.SysUser)
		err = copier.Copy(sysUser, req)
		if err != nil {
			return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
		}

		bytes, err := json.Marshal(req.RoleIds)
		sysUser.RoleIds = string(bytes)
		dictionary := &sysadmin.SysDictionary{}
		err = l.svcCtx.GormQuery.FindOneBy(l.ctx, dictionary, "unique_key", "sys_pwd")
		var password string
		if dictionary.Status == globalkey.SysEnable {
			password = dictionary.Value
		} else {
			password = globalkey.SysNewUserDefaultPassword
		}

		salt := utils.GetRandom(4, 3)
		sysUser.Salt = salt
		sysUser.Password = utils.MD5(password + salt + l.svcCtx.Config.Salt)
		sysUser.Avatar = utils.AvatarUrl()
		err = l.svcCtx.GormQuery.Save(l.ctx, sysUser)
		if err != nil {
			return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
		}

		return nil
	} else {

		return errorx.NewDefaultError(errorx.AddUserErrorCode)
	}
}
