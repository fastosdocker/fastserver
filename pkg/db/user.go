package db

// User 用户模型
type User struct {
	*BaseModel
	// 账号
	UserName string `gorm:"not null;unique_index" json:"userName" form:"userName"`
	// 密码
	Password string `gorm:"not null" json:"password" form:"password"`
}

func (b User) TableName() string {
	return "user"
}
