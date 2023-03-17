package user

import (
	"context"
	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/app/model"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"
	"gozore-mall/service/gormpool"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSysUserPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSysUserPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSysUserPageLogic {
	return &GetSysUserPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSysUserPageLogic) GetSysUserPage(req *types.SysUserPageReq) (resp *types.SysUserPageResp, err error) {
	s := strconv.FormatInt(req.DeptId, 10)
	deptIds := l.getDeptIds(s, req.DeptId)

	users := make([]sysadmin.SysUserDetail, 0)
	querySql := "SELECT u.id,u.dept_id,u.job_id,u.profession_id,u.account,u.username,u.nickname,u.avatar,u.gender,IFNULL(p.name,'NULL') as profession,IFNULL(j.name,'NULL') as job,IFNULL(d.name,'NULL') as dept,IFNULL(GROUP_CONCAT(r.name),'NULL') as roles,IFNULL(GROUP_CONCAT(r.id),0) as role_ids,u.email,u.mobile,u.remark,u.order_num,u.status,u.created,u.updated FROM (SELECT * FROM sys_user WHERE id!=? AND dept_id IN(?) ORDER BY order_num DESC LIMIT ?,?) u LEFT JOIN sys_profession p ON u.profession_id=p.id LEFT JOIN sys_dept d ON u.dept_id=d.id LEFT JOIN sys_job j ON u.job_id=j.id LEFT JOIN sys_role r ON JSON_CONTAINS(u.role_ids,JSON_ARRAY(r.id)) GROUP BY u.id"
	err = l.svcCtx.GormQuery.Query(l.ctx, querySql, []interface{}{globalkey.SysSuperUserId, deptIds, (req.Page - 1) * req.Limit, req.Limit}, &users)
	if err != nil {
		return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	var user types.User
	var userProfession types.UserProfession
	var userJob types.UserJob
	var userDept types.UserDept
	userList := make([]types.User, 0)
	for _, v := range users {
		err := copier.Copy(&user, &v)
		if err != nil {
			return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
		}

		userProfession.Id = v.ProfessionId
		userProfession.Name = v.Profession

		userJob.Id = v.JobId
		userJob.Name = v.Job

		userDept.Id = v.DeptId
		userDept.Name = v.Dept

		var userRole types.UserRole
		var roles []types.UserRole
		var roleNameArr []string
		var roleIdArr []string
		roleNameArr = strings.Split(v.Roles, ",")
		roleIdArr = strings.Split(v.RoleIds, ",")
		for i, n := range roleNameArr {
			userRole.Name = n
			userRole.Id, _ = strconv.ParseInt(roleIdArr[i], 10, 64)
			roles = append(roles, userRole)
		}

		user.Profession = userProfession
		user.Job = userJob
		user.Dept = userDept
		user.Roles = roles

		userList = append(userList, user)
	}

	total, err := l.svcCtx.GormQuery.CountTotal(l.ctx, &sysadmin.SysUser{}, gormpool.Conditions{
		Where:  "id!=? AND dept_id IN(?)",
		Values: []interface{}{globalkey.SysSuperUserId, deptIds},
	})
	if err != nil {
		return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	pagination := types.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}

	return &types.SysUserPageResp{
		List:       userList,
		Pagination: pagination,
	}, nil
}

func (l *GetSysUserPageLogic) getDeptIds(deptId string, id int64) string {
	deptList := make([]sysadmin.SysDept, 0)
	err := l.svcCtx.GormQuery.QueryFindAll(l.ctx, &deptList, gormpool.Conditions{
		Where:  "parent_id=?",
		Values: []interface{}{id},
	})
	if err != nil && err != model.ErrNotFound {
		return deptId
	}

	for _, v := range deptList {
		deptId = deptId + "," + strconv.FormatInt(v.Id, 10)
		deptId = l.getDeptIds(deptId, v.Id)
	}

	return deptId
}
