package config

import (
	"fast/utils"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

var (
	Version       = "v1.0.0"
	Conf          conf
	path          = "config/config.yaml"
	ServerBaseUrl = os.Getenv("FAST_STORE")
)

func init() {
	if ok, _ := utils.FileExists(path); !ok {
		// 不存在配置文件
		Conf.Port = "8081"
		Conf.DaemonPath = "/etc/docker/daemon.json"
		Conf.Https.Port = "8082"
		Conf.TLs.Key = "config/server-sign/cert.key"
		Conf.TLs.Pem = "config/server-sign/cert.crt"

		// 创建server-sign证书路径
		err := utils.MkDir("config/server-sign")
		if err != nil {
			log.Fatalf("创建server-sign文件夹失败: %s", err)
		}
		// 创建user-sign 证书路径
		err = utils.MkDir("config/user-sign")
		if err != nil {
			log.Fatalf("创建user-sign文件夹失败: %s", err)
		}

		err = Save()
		if err != nil {
			log.Fatalf("初始化config.yaml配置文件错误: %s", err)
		}

		return
	}

	// 存在配置文件
	byt, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("config.yaml读取失败: %s", err)
	}

	err = yaml.Unmarshal(byt, &Conf)
	if err != nil {
		log.Fatalf("反序列化配置文件失败: %s", err)
	}
}

func Save() error {
	data, err := yaml.Marshal(&Conf)
	if err != nil {
		return err
	}

	return utils.SaveFile(path, string(data))
}
