package container

import (
	"context"
	"fast/pkg/db"
	"fast/pkg/docker"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types"
)

// Stats 记录容器数据
type Stats struct{}

func (s Stats) Name() string {
	return "记录容器数据"
}

func (s Stats) GetDuration() time.Duration {
	return time.Second * 6
}

func (s Stats) Run() error {
	// 获取容器列表
	containers, err := docker.Get().ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
	}

	for _, container := range containers {
		var (
			stats = types.ContainerStats{}
			data  []byte
		)
		// 获取容器状态
		stats, err = docker.Get().ContainerStats(context.Background(), container.ID, false)
		if err != nil {
			continue
		}

		data, err = io.ReadAll(stats.Body)
		if err != nil {
			continue
		}

		defer stats.Body.Close()

		//把数据保存到数据库
		err = db.Create(&db.ContainerStats{
			ContainerID:    container.ID,
			ContainerStats: string(data),
			CreateTime:     time.Now(),
		})
		if err != nil {
			log.Printf("保存容器统计信息失败: %s\n", err)
			continue
		}
	}

	return nil
}

// Clean 清理容器过期数据
type Clean struct {
}

func (c Clean) Name() string {
	return "清理容器过期数据"
}

func (c Clean) GetDuration() time.Duration {
	return time.Hour * 1
}

func (c Clean) Run() error {
	// 获取容器列表
	containers, err := docker.Get().ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return err
	}

	createTime := time.Now().Add(-30 * time.Minute).Format("2006-01-02 15:04:05.000")
	for _, v := range containers {
		//删除旧的历史记录
		err = db.Delete(&db.ContainerStats{}, fmt.Sprintf("`container_id`='%s' and `create_time`<'%s'", v.ID, createTime))
		if err != nil {
			return err
		}
	}

	return nil
}
