package auth

import (
	"time"

	casbin "github.com/casbin/casbin/v2"
	model "github.com/casbin/casbin/v2/model"
	adapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

const (
	// aclModel 定义了 casbin 访问控制模型.
	aclModel = `
		[request_definition]
		r = sub, obj, act

		[policy_definition]
		p = sub, obj, act, eft

		[role_definition]
		g = _, _

		[policy_effect]
		e = !some(where (p.eft == deny))

		[matchers]
		m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && r.act == p.act
		`
)

// Authz 定义了一个授权器，提供授权功能.
type Authz struct {
	// 使用 Casbin 的同步授权器
	*casbin.SyncedEnforcer
}

func NewAuthz(db *gorm.DB) (*Authz, error) {
	// 初始化 Gorm 适配器并用于 Casbin 授权器
	adapter, err := adapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// 从字符串中创建 Casbin 模型
	m, _ := model.NewModelFromString(aclModel)

	enforcer, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	// 从数据库加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	// 启动自动加载策略，间隔为 5 秒
	enforcer.StartAutoLoadPolicy(5 * time.Second)

	return &Authz{enforcer}, nil
}

func (a *Authz) Authorize(sub, obj, act string) (bool, error) {
	return a.Enforce(sub, obj, act)
}
