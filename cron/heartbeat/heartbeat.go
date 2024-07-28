package heartbeat

import (
	"fast/pkg/http"
	"fast/utils"
	"fmt"
	"time"
)

// Heartbeat 给服务器发送心跳
type Heartbeat struct{}

func (h Heartbeat) Name() string {
	return "给服务器发送心跳"
}

func (h Heartbeat) GetDuration() time.Duration {
	return time.Hour
}

// Run 发送心跳到服务器
func (h Heartbeat) Run() error {
	var postMap = make(map[string]interface{})
	postMap["key"] = fmt.Sprintf("FAST-%s", fmt.Sprintf("FAST%s-%s", time.Now().Format("20060102")[2:], utils.GetRandomString(12)))

	return http.Post("http://dockernb.com:8083/api/ping", postMap)
}
