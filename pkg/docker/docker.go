package docker

import (
	"context"
	"github.com/docker/docker/client"
	"log"
)

var (
	cli *client.Client
)

func init() {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	cli.NegotiateAPIVersion(context.Background())
	if err != nil {
		log.Fatalf("创建docker客户端失败: %s", err)
	}
}

// Get 获取docker client
func Get() *client.Client {
	return cli
}
