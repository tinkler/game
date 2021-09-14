package status

// Status 服务,角色,用户,物品等状态的方法集定义
// 配合logger使用，统一日志格式
type Status interface {
	ID() string
	TypeName() string
}
