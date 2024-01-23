package rbac

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"time"
)

type Synced struct {
	Rbac
	enforcer *casbin.SyncedEnforcer
}

func NewSynced(db *gorm.DB) (*Synced, error) {
	var u Synced
	// Gorm 适配器
	adapter, err := gormadapter.NewAdapterByDBWithCustomTable(db, u.CasbinRule, u.CasbinRule.TableName())
	if err != nil {
		return &u, err
	}
	// 通过mysql适配器新建一个enforcer
	u.enforcer, err = casbin.NewSyncedEnforcer("configs/rbac_model.conf", adapter, false)
	if err != nil {
		return &u, err
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
	u.enforcer.StartAutoLoadPolicy(300 * time.Second)
	u.PermissionPrefix = "p:"
	u.RolePrefix = "r:"
	u.UserPrefix = "u:"
	return &u, nil
}

// CheckPolicy 检测策略规则
// sub: user_id|role_id|permission_id
func (u *Synced) CheckPolicy(sub string, path string, method string) bool {
	res, err := u.enforcer.Enforce(sub, path, method)
	if err != nil {
		return false
	}
	return res
}

// CheckPermission 检测权限
func (u *Synced) CheckPermission(userId string, permission string) bool {
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

// LoadPolicy 加载策略规则
func (u *Synced) LoadPolicy() error {
	err := u.enforcer.LoadPolicy()
	if err != nil {
		return err
	}
	return nil
}
