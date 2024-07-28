package db

// Application 应用模型
type Application struct {
	*BaseModel
	// 项目名称
	Name string `gorm:"not null;unique_index" json:"name" form:"name"`
	// 图标
	Image string `gorm:"not null" json:"image" form:"image"`
	// 版本号
	V string `gorm:"not null" json:"v" form:"v"`
	// 描述
	Description string `gorm:"not null" json:"description" form:"description"`
	// json格式数据
	Compose string `gorm:"not null" json:"compose" form:"compose"`
	// 状态
	Status string `json:"status" form:"status"`
	Note   string `json:"note" form:"note"`
}

func (b *Application) TableName() string {
	return "application"
}
