package repo

import (
	"context"
	"encoding/json"
	"fast/pkg/docker"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"io"
)

func FetchingDockerContainerRepo() {
	var (
		cli = docker.Get()
		ctx = context.Background()
	)
	// 获取所有正在运行的容器
	containers, err := cli.ContainerList(ctx, container.ListOptions{
		Size:    true,
		All:     false,
		Latest:  true,
		Filters: filters.Args{},
	})
	if err != nil {
		panic(err)
	}

	// 控制并发
	concurrencyLimit := make(chan struct{}, len(containers))

	// 处理容器统计信息的函数
	handleContainerStats := func(ctx context.Context, cont types.Container) {
		defer func() { concurrencyLimit <- struct{}{} }()
		var (
			cID    = cont.ID[:12]
			reader container.StatsResponseReader
			buf    = newBuffer(cID)
		)

		fmt.Printf("Fetching stats for container %s...\n", cID)

		// 开始接收容器的统计信息
		reader, err = cli.ContainerStats(ctx, cont.ID, true)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer reader.Body.Close()

		// 解码JSON流并处理数据
		dec := json.NewDecoder(reader.Body)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Operation cancelled for container", cID)
				return
			default:
				var statsResponse container.StatsResponse
				err = dec.Decode(&statsResponse)
				if err != nil {
					if err == io.EOF {
						break
					}
					continue
				}
				// 将新行添加到缓冲区，如果缓冲区已满，移除最旧的行
				newRecord := []string{
					statsResponse.ID[:12],
					statsResponse.Name,
					statsResponse.Read.Format("2006-01-02 15:04:05"),
					fmt.Sprintf("%v", statsResponse.MemoryStats.Usage),
					fmt.Sprintf("%v", statsResponse.CPUStats.CPUUsage.TotalUsage),
					fmt.Sprintf("%v", statsResponse.Networks["eth0"].RxBytes),
					fmt.Sprintf("%v", statsResponse.Networks["eth0"].TxBytes),
				}

				buf.append(newRecord)
			}
		}
	}
	// 启动goroutines来处理容器统计信息
	for _, c := range containers {
		concurrencyLimit <- struct{}{}
		go handleContainerStats(ctx, c)
	}
}
