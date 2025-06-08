package auth

import (
	"time"

	casbin "github.com/casbin/casbin/v2"
	model "github.com/casbin/casbin/v2/model"
	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/google/wire"
	"gorm.io/gorm"
)

const (
	// 默认的 Casbin 访问控制模型.
	defaultAclModel = `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
e = !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && r.act == p.act`
)

var ProviderSet = wire.NewSet(NewAuthz, DefaultOptions)

// Authz 定义了一个授权器，提供授权功能.
type Authz struct {
	// 使用 Casbin 的同步授权器
	*casbin.SyncedEnforcer
}

// Option 定义了一个函数选项类型，用于自定义 NewAuthz 的行为.
type Option func(*authzConfig)

// authzConfig 是授权器的配置结构.
type authzConfig struct {
	aclModel           string        // Casbin 的模型字符串
	autoLoadPolicyTime time.Duration // 自动加载策略的时间间隔
}

// defaultAuthzConfig 返回一个默认的配置.
func defaultAuthzConfig() *authzConfig {
	return &authzConfig{
		// 默认使用内置的 ACL 模型
		aclModel: defaultAclModel,
		// 默认的自动加载策略时间间隔
		autoLoadPolicyTime: 5 * time.Second,
	}
}

// DefaultOptions 提供默认的授权器选项配置.
func DefaultOptions() []Option {
	return []Option{
		// 使用默认的 ACL 模型
		WithAclModel(defaultAclModel),
		// 设置自动加载策略的时间间隔为 10 秒
		WithAutoLoadPolicyTime(10 * time.Second),
	}
}

// WithAclModel 允许通过选项自定义 ACL 模型.
func WithAclModel(model string) Option {
	return func(cfg *authzConfig) {
		cfg.aclModel = model
	}
}

// WithAutoLoadPolicyTime 允许通过选项自定义自动加载策略的时间间隔.
func WithAutoLoadPolicyTime(interval time.Duration) Option {
	return func(cfg *authzConfig) {
		cfg.autoLoadPolicyTime = interval
	}
}

func NewAuthz(db *gorm.DB) (*Authz, error) {
	// 初始化默认配置
	cfg := defaultAuthzConfig()

	// 初始化 Gorm 适配器并用于 Casbin 授权器
	adapter, err := adapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// 从字符串中创建 Casbin 模型
	m, _ := model.NewModelFromString(cfg.aclModel)

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
