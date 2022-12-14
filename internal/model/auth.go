package model

// 表名小写复数 users，字段单词小写下划线，即蛇形Python风格
type Account struct {
	*Model
	Username string `json:"username"` // 用户名
	Password string `json:"-"`        // 密码
	Role     uint8  `json:"role"`     // 角色 0 普通用户、1 管理员
	State    uint8  `json:"state"`    // 状态 0 为禁用、1 为启用
}

func (u Account) TableName() string {
	return "auth_account"
}
