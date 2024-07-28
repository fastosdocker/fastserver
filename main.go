package main

import (
	"fast/cron"
	"fast/pkg/compose"
	"fast/pkg/db"
	"fast/utils"
	"fast/web"
	dockerRepo "fast/web/repo"
	"log"
	"time"
)

func main() {
	//初始化docker-compose
	if err := compose.Init(compose.Option{
		SH:       "install_compose.sh",
		YamlPath: "data/yaml/",
	}); err != nil {
		log.Fatalln(err)
	}
	//初始化数据库
	dsn := "data/test.db"
	if err := utils.MkDir("data"); err != nil {
		log.Fatalf("创建data目录失败: %s", err)
	}

	if err := db.Init(dsn); err != nil {
		log.Fatalln(err)
	}
	s := time.Now()
	//启动定时任务服务
	cron.Start()
	dockerRepo.FetchingDockerContainerRepo()
	//启动web服务
	web.Start(s)
}
