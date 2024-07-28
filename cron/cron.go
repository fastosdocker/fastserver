package cron

import (
	"fast/cron/container"
	"fast/cron/heartbeat"
	"log"
	"time"
)

var (
	crons []Cron
)

func init() {
	Register(
		container.Stats{},
		container.Clean{},
		heartbeat.Heartbeat{},
	)
}

type Cron interface {
	Name() string
	Run() error
	GetDuration() time.Duration
}

func Register(cron ...Cron) {
	crons = append(crons, cron...)
}

func Start() {
	for _, v := range crons {
		log.Printf("启动[%s]定时任务", v.Name())
		go func(cron Cron) {
			var (
				err      error
				errCount int
			)

			for {
				err = cron.Run()
				if err != nil {
					log.Println(err)

					if errCount < 300 {
						errCount++
					}
					time.Sleep(time.Second * time.Duration(errCount))

					continue
				}

				time.Sleep(cron.GetDuration())
			}
		}(v)
	}
}
