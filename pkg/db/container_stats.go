package db

import "time"

// ContainerStats 容器stats模型
type ContainerStats struct {
	ID int `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	//容器ID
	ContainerID string `gorm:"not null;index:idx_query" json:"containerId" form:"containerId"`
	//容器stats
	ContainerStats string `gorm:"not null" json:"containerStats" form:"containerStats"`
	//创建时间
	CreateTime time.Time `gorm:"autoCreateTime;not null;index:idx_query" json:"createTime" form:"createTime"`
}

func (b *ContainerStats) TableName() string {
	return "ContainerStats"
}
