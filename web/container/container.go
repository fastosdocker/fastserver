package container

import (
	"context"
	"fast/pkg/db"
	"fast/pkg/docker"
	"fast/web/base"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
)

// List 获取容器列表
func List(c *base.Ctx) {
	data, err := docker.Get().ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		c.Error("获取容器列表失败: " + err.Error())
		return
	}

	c.Success(data)
}

// Create 创建容器
func Create(c *base.Ctx) {
	params := CreateModel{}
	err := c.BindJSON(&params)
	if err != nil {
		return
	}

	_, err = docker.Get().ContainerCreate(context.Background(), &params.Config, &params.HostConfig, &params.NetworkingConfig, nil, params.Name)
	if err != nil {
		c.Error("创建容器失败: " + err.Error())
		return
	}

	c.Success()
}

func Delete(c *base.Ctx) {
	containerId := c.PostForm("containerId")
	if containerId == "" {
		c.Error("容器ID不能为空")
		return
	}

	err := docker.Get().ContainerRemove(context.Background(), containerId, types.ContainerRemoveOptions{})
	if err != nil {
		c.Error("删除容器失败: " + err.Error())
		return
	}

	//同时删除容器的统计记录
	err = db.Delete(&db.ContainerStats{}, fmt.Sprintf("`container_id`='%s'", containerId))
	if err != nil {
		c.Error("删除容器stats记录失败: " + err.Error())
		return
	}

	c.Success()
}

func HistoryStats(c *base.Ctx) {
	// 根据容器ID获取数据
	data := make([]*db.ContainerStats, 0)
	err := db.Find(&data, db.Query{
		Where: fmt.Sprintf("`container_id`='%s' and `create_time`>'%s'", c.PostForm("containerId"), time.Now().Add(-30*time.Minute)),
		Order: "id",
	})
	if err != nil {
		c.Error(fmt.Sprintf("查询失败: %s", err))
		return
	}

	c.Success(data)
}

func Stats(c *base.Ctx) {
	// 根据容器ID获取数据
	containerId := c.PostForm("containerId")
	if containerId == "" {
		c.Error("容器ID不能为空")
		return
	}

	stats, err := docker.Get().ContainerStats(context.Background(), containerId, false)
	if err != nil {
		c.Error("获取容器stats失败: " + err.Error())
		return
	}

	data, err := io.ReadAll(stats.Body)
	if err != nil {
		c.Error("获取容器stats失败: " + err.Error())
		return
	}

	c.Success(string(data))
}
