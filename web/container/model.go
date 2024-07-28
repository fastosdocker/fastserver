package container

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

type CreateModel struct {
	Name             string                   `json:"Name"`
	Config           container.Config         `json:"Config"`
	HostConfig       container.HostConfig     `json:"HostConfig"`
	NetworkingConfig network.NetworkingConfig `json:"NetworkingConfig"`
}
