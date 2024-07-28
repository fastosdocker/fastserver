package db

import (
	"time"
)

type BaseModel struct {
	ID         int       `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	CreateTime time.Time `gorm:"autoCreateTime;not null" json:"createTime" form:"createTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime;not null" json:"updateTime" form:"updateTime"`
	// 0：正常，1：删除
	Del int `gorm:"not null;default:0" json:"del" form:"del"`
}
