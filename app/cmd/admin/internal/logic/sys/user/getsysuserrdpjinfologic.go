package user

import (
	"context"
	"encoding/json"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/common/utils"
	"gozore-mall/service/gormpool"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/globalkey"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSysUserRdpjInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSysUserRdpjInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSysUserRdpjInfoLogic {
	return &GetSysUserRdpjInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSysUserRdpjInfoLogic) GetSysUserRdpjInfo(req *types.GetSysUserRdpjInfoReq) (resp *types.GetSysUserRdpjInfoResp, err error) {
	currentUserId := utils.GetAdminUserId(l.ctx)
	return &types.GetSysUserRdpjInfoResp{
		Role:       l.roleList(currentUserId, req.UserId),
		Dept:       l.deptList(),
		Profession: l.professionList(),
		Job:        l.jobList(),
	}, nil
}

func (l *GetSysUserRdpjInfoLogic) roleList(currentUserId int64, editUserId int64) []types.RoleTree {
	var currentUserRoleIds []int64
	var roleIds []int64
	var sysRoleList []sysadmin.SysRole
	if currentUserId == globalkey.SysSuperUserId {
		_ = l.svcCtx.GormQuery.QueryFindAll(l.ctx, &sysRoleList, gormpool.Conditions{
			OrderBy: "order_num DESC",
		})
		for _, role := range sysRoleList {
			currentUserRoleIds = append(currentUserRoleIds, role.Id)
		}

	} else {
		currentUser := &sysadmin.SysUser{}
		_ = l.svcCtx.GormQuery.FindOne(l.ctx, currentUser, currentUserId)
		err := json.Unmarshal([]byte(currentUser.RoleIds), &currentUserRoleIds)
		if err != nil {
			return nil
		}

		roleIds = append(roleIds, currentUserRoleIds...)
		if editUserId != 0 {
			editUser := &sysadmin.SysUser{}
			_ = l.svcCtx.GormQuery.FindOne(l.ctx, editUser, editUserId)
			var editUserRoleIds []int64
			err = json.Unmarshal([]byte(editUser.RoleIds), &editUserRoleIds)
			if err != nil {
				return nil
			}

			roleIds = append(roleIds, editUserRoleIds...)
		}

		_ = l.svcCtx.GormQuery.QueryFindAll(l.ctx, &sysRoleList, gormpool.Conditions{
			Where:  "id!=? AND status=1 AND id IN(?)",
			Values: []interface{}{globalkey.SysSuperRoleId, roleIds},
		})
	}

	var role types.RoleTree
	roleList := make([]types.RoleTree, 0)
	for _, v := range sysRoleList {
		err := copier.Copy(&role, &v)
		if err != nil {
			return nil
		}

		if utils.ArrayContainValue(currentUserRoleIds, v.Id) {
			role.Has = 1
		} else {
			role.Has = 0
		}

		roleList = append(roleList, role)
	}

	return roleList
}

func (l *GetSysUserRdpjInfoLogic) deptList() []types.DeptTree {
	sysDeptList := make([]sysadmin.SysDept, 0)
	_ = l.svcCtx.GormQuery.QueryFindAll(l.ctx, &sysDeptList, gormpool.Conditions{
		Where:   "status=1",
		OrderBy: "order_num DESC",
	})
	var dept types.DeptTree
	deptList := make([]types.DeptTree, 0)
	for _, v := range sysDeptList {
		err := copier.Copy(&dept, &v)
		if err != nil {
			return nil
		}
		deptList = append(deptList, dept)
	}

	return deptList
}

func (l *GetSysUserRdpjInfoLogic) professionList() []types.Rdpj {
	sysProfessionList := make([]sysadmin.SysProfession, 0)
	_ = l.svcCtx.GormQuery.QueryFindAll(l.ctx, &sysProfessionList, gormpool.Conditions{
		Where:   "status=1",
		OrderBy: "order_num DESC",
	})
	var profession types.Rdpj
	professionList := make([]types.Rdpj, 0)
	for _, v := range sysProfessionList {
		err := copier.Copy(&profession, &v)
		if err != nil {
			return nil
		}
		professionList = append(professionList, profession)
	}

	return professionList
}

func (l *GetSysUserRdpjInfoLogic) jobList() []types.Rdpj {
	sysJobList := make([]sysadmin.SysJob, 0)
	_ = l.svcCtx.GormQuery.QueryFindAll(l.ctx, &sysJobList, gormpool.Conditions{
		Where:   "status=1",
		OrderBy: "order_num DESC",
	})
	var job types.Rdpj
	jobList := make([]types.Rdpj, 0)
	for _, v := range sysJobList {
		err := copier.Copy(&job, &v)
		if err != nil {
			return nil
		}
		jobList = append(jobList, job)
	}

	return jobList
}
