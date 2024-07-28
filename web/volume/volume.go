package volume

import (
	"context"
	"fast/pkg/docker"
	"fast/web/base"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
)

// Create 创建数据卷
func Create(c *base.Ctx) {
	params := volume.VolumeCreateBody{}
	err := c.BindJSON(&params)
	if err != nil {
		c.Error("参数错误: " + err.Error())
		return
	}

	_, err = docker.Get().VolumeCreate(context.Background(), params)
	if err != nil {
		c.Error("创建数据卷失败: " + err.Error())
		return
	}

	c.Success()
}

// List 获取数据卷列表
func List(c *base.Ctx) {
	params := make([]filters.KeyValuePair, 0)
	err := c.BindJSON(&params)
	if err != nil {
		c.Error("参数错误: " + err.Error())
		return
	}

	rsp, err := docker.Get().VolumeList(context.Background(), filters.NewArgs(params...))
	if err != nil {
		c.Error("获取数据卷失败: " + err.Error())
		return
	}

	c.Success(rsp)
}

// Remove 删除数据卷
func Remove(c *base.Ctx) {
	id := c.PostForm("id")
	forceStr := c.PostForm("force")

	if id == "" {
		c.Error("数据卷id不能为空")
		return
	}

	force := false
	if forceStr == "true" || forceStr == "1" {
		force = true
	}

	err := docker.Get().VolumeRemove(context.Background(), id, force)
	if err != nil {
		c.Error("删除数据卷失败: " + err.Error())
		return
	}

	c.Success()
}
