package rbac

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"strings"
)

type CasbinRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:100;uniqueIndex:unique_index"`
	V0    string `gorm:"size:100;uniqueIndex:unique_index"`
	V1    string `gorm:"size:100;uniqueIndex:unique_index"`
	V2    string `gorm:"size:100;uniqueIndex:unique_index"`
	V3    string `gorm:"size:100;uniqueIndex:unique_index"`
	V4    string `gorm:"size:100;uniqueIndex:unique_index"`
	V5    string `gorm:"size:100;uniqueIndex:unique_index"`
}

func (CasbinRule) TableName() string {
	return "admin_casbin_rule"
}

type Rbac struct {
	CasbinRule       CasbinRule
	PermissionPrefix string
	RolePrefix       string
	UserPrefix       string
	enforcer         *casbin.Enforcer
	//enforcer *casbin.SyncedEnforcer
}

func (u *Rbac) New(db *gorm.DB) (*Rbac, error) {

	// Gorm 适配器
	adapter, err := gormadapter.NewAdapterByDBWithCustomTable(db, u.CasbinRule, u.CasbinRule.TableName())
	if err != nil {
		return u, err
	}
	// 通过mysql适配器新建一个enforcer
	u.enforcer, err = casbin.NewEnforcer("configs/rbac_model.conf", adapter, false)
	//u.enforcer, err = casbin.NewSyncedEnforcer("rbac_model.conf", adapter, false)
	if err != nil {
		return u, err
	}
	// 是否自动保存 默认开启
	//u.enforcer.EnableAutoSave(false)
	// 日志记录
	//u.enforcer.EnableLog(true)

	// 加载策略规则
	//err = u.LoadPolicy()
	//if err != nil {
	//	return u, err
	//}
	//u.enforcer.StartAutoLoadPolicy(1 * time.Second)
	u.PermissionPrefix = "p:"
	u.RolePrefix = "r:"
	u.UserPrefix = "u:"
	return u, nil
}

// AddPermission 批量添加权限
func (u *Rbac) AddPermission(permissionId string, paths, methods []string) error {
	var permissions [][]string
	method := "*"
	if len(methods) > 0 {
		method = strings.Join(methods, "|")
	}
	for _, path := range paths {
		permissions = append(permissions, []string{u.PermissionPrefix + permissionId, path, method, "allow"})
	}
	_, err := u.enforcer.AddPolicies(permissions)
	if err != nil {
		return err
	}
	return nil
}

// RemovePermission 批量删除权限
// permission 权限标识
// removeBinding 是否删除角色绑定的权限
func (u *Rbac) RemovePermission(permissionId string, removeBinding bool) (err error) {
	_, err = u.enforcer.RemoveFilteredPolicy(0, u.PermissionPrefix+permissionId)
	if err != nil {
		return err
	}
	if !removeBinding {
		return nil
	}
	_, err = u.enforcer.RemoveFilteredGroupingPolicy(1, u.PermissionPrefix+permissionId)
	if err != nil {
		return err
	}
	return nil
}

// CheckPermission 检测权限
func (u *Rbac) CheckPermission(userId string, permission string) bool {
	// 成功返回true, 已存在返回false
	// 子角色无法直接判断，需先获取包括子角色的所有角色
	roles, _ := u.enforcer.GetImplicitRolesForUser(u.UserPrefix + userId)
	for _, v := range roles {
		if v == u.PermissionPrefix+permission {
			return true
		}
	}
	return false
}

// AddRole 批量添加角色
func (u *Rbac) AddRole(roleId string, permissions []string) (err error) {
	var roles [][]string
	for _, permission := range permissions {
		roles = append(roles, []string{u.RolePrefix + roleId, u.PermissionPrefix + permission})
	}
	_, err = u.enforcer.AddGroupingPolicies(roles)
	if err != nil {
		return err
	}
	return nil
}

// AddUserRoles 批量添加角色
func (u *Rbac) AddUserRoles(userId string, roleIds []string) (err error) {
	var roles [][]string
	for _, roleId := range roleIds {
		roles = append(roles, []string{u.UserPrefix + userId, u.RolePrefix + roleId})
	}
	_, err = u.enforcer.AddGroupingPolicies(roles)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserGroupingPolicy 批量删除用户角色和权限
func (u *Rbac) DeleteUserGroupingPolicy(userId string) (err error) {
	_, err = u.enforcer.RemoveFilteredGroupingPolicy(0, u.UserPrefix+userId)
	if err != nil {
		return err
	}
	return nil
}

// AddUserPermissions 批量添加权限
func (u *Rbac) AddUserPermissions(roleId string, permissionIds []string) (err error) {
	var permissions [][]string
	for _, permission := range permissionIds {
		permissions = append(permissions, []string{u.RolePrefix + roleId, u.PermissionPrefix + permission})
	}
	_, err = u.enforcer.AddGroupingPolicies(permissions)
	if err != nil {
		return err
	}
	return nil
}

// RemoveRole 批量删除角色
func (u *Rbac) RemoveRole(roleId string, removeBinding bool) (err error) {
	_, err = u.enforcer.RemoveFilteredGroupingPolicy(0, u.RolePrefix+roleId)
	if err != nil {
		return err
	}
	if !removeBinding {
		return nil
	}
	_, err = u.enforcer.RemoveFilteredGroupingPolicy(1, u.RolePrefix+roleId)
	if err != nil {
		return err
	}
	return nil
}

// CheckPolicy 检测策略规则
// sub: user_id|role_id|permission_id
func (u *Rbac) CheckPolicy(sub string, path string, method string) bool {
	res, err := u.enforcer.Enforce(sub, path, method)
	if err != nil {
		return false
	}
	return res
}

func (u *Rbac) SavePolicy() error {
	err := u.enforcer.SavePolicy()
	if err != nil {
		return err
	}
	return nil
}

// LoadPolicy 加载策略规则
func (u *Rbac) LoadPolicy() error {
	err := u.enforcer.LoadPolicy()
	if err != nil {
		return err
	}
	return nil
}
