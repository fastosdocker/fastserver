package home

import (
	"context"
	"fast/pkg/docker"
	"fast/web/base"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	volumetypes "github.com/docker/docker/api/types/volume"
)

type HomerRsp struct {
	Info    types.Info
	Network []types.NetworkResource
	Volume  volumetypes.VolumeListOKBody
}

// Overview 服务器总览
func Overview(c *base.Ctx) {

	info, err := docker.Get().Info(context.Background())
	if err != nil {
		c.Error("获取info信息失败: " + err.Error())
		return
	}

	network, err := docker.Get().NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		c.Error("获取network信息失败: " + err.Error())
		return
	}

	volume, err := docker.Get().VolumeList(context.Background(), filters.Args{})
	if err != nil {
		c.Error("获取volume信息失败: " + err.Error())
		return
	}

	c.Success(HomerRsp{Info: info, Network: network, Volume: volume})
}
