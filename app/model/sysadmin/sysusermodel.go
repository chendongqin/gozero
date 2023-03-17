package sysadmin

import (
	"context"
	"gozore-mall/common/globalkey"
	"gozore-mall/common/utils"
	"gozore-mall/service/gormpool"
	"strings"
	"time"
)

type GormSysUserModel struct {
	gormpool.GormBaseModel
}

type SysUserDetail struct {
	Id           int64     `json:"id"`            // 编号
	Account      string    `json:"account"`       // 账号
	Username     string    `json:"username"`      // 姓名
	Nickname     string    `json:"nickname"`      // 昵称
	Avatar       string    `json:"avatar"`        // 头像
	Gender       int64     `json:"gender"`        // 0=保密 1=女 2=男
	Profession   string    `json:"profession"`    // 职称
	ProfessionId int64     `json:"profession_id"` // 职称id
	Job          string    `json:"job"`           // 岗位
	JobId        int64     `json:"job_id"`        // 岗位id
	Dept         string    `json:"dept"`          // 部门
	DeptId       int64     `json:"dept_id"`       // 部门id
	Roles        string    `json:"roles"`         // 角色集
	RoleIds      string    `json:"role_ids"`      // 角色集id
	Email        string    `json:"email"`         // 邮件
	Mobile       string    `json:"mobile"`        // 手机号
	Remark       string    `json:"remark"`        // 备注
	OrderNum     int64     `json:"order_num"`     // 排序值
	Status       int64     `json:"status"`        // 0=禁用 1=开启
	Created      time.Time `json:"created"`       // 创建时间
	Updated      time.Time `json:"updated"`       // 更新时间
}

func NewGormSysUserModel(ctx context.Context, gormQuery *gormpool.CommonQuery) *GormSysUserModel {
	return &GormSysUserModel{
		GormBaseModel: gormpool.GormBaseModel{
			Ctx:       ctx,
			GormQuery: gormQuery,
		},
	}
}

func (m *GormSysUserModel) CheckSuperRole(userId int64) bool {
	sysUser := &SysUser{}
	_ = m.GormQuery.FindOne(m.Ctx, sysUser, userId)
	roleIds := strings.Split(sysUser.RoleIds, ",")
	check := false
	for _, roleId := range roleIds {
		if utils.ToInt64(roleId) == globalkey.SysSuperRoleId {
			check = true
			break
		}
	}

	return check
}
