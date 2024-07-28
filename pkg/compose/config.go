package compose

import (
	"bytes"
	"encoding/json"
	"fast/utils"
	"os"
	"path/filepath"

	"github.com/compose-spec/compose-go/types"
	"gopkg.in/yaml.v3"
)

type Option struct {
	SH       string
	YamlPath string
}

type ProjectConfig struct {
	// 版本
	Version string `json:"version,omitempty"`
	// 服务列表
	Services map[string]Service `json:"services,omitempty"`
	// 网络
	Networks types.Networks `json:"networks,omitempty"`
	// 数据卷
	Volumes types.Volumes `json:"volumes,omitempty"`
}

type Service struct {
	// 镜像名称
	Image string `yaml:"image" json:"image,omitempty"`
	// 端口
	Ports []types.ServicePortConfig `yaml:"ports" json:"ports,omitempty"`
	// 重启策略
	Restart string `yaml:"restart" json:"restart,omitempty"`
	// 文件映射
	Volumes []string `yaml:"volumes" json:"volumes,omitempty"`
	// 容器名称
	ContainerName string `yaml:"container_name" json:"container_name,omitempty"`
	// 环境变量
	Environment map[string]interface{} `yaml:"environment" json:"environment,omitempty"`
	// 网络
	Networks []string `yaml:"networks" json:"networks,omitempty"`
	// 标签
	Labels map[string]string `yaml:"labels" json:"labels,omitempty"`
}

type stringWriter struct {
	buffer bytes.Buffer
}

func (sw *stringWriter) Write(p []byte) (n int, err error) {
	return sw.buffer.Write(p)
}

func (sw *stringWriter) String() string {
	return sw.buffer.String()
}

func jsonToYaml(jsonData []byte) ([]byte, error) {
	data := make(map[string]interface{}, 0)
	// jsonData=append([]byte("{"),jsonData...)
	// jsonData=append(jsonData,'}')

	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	return yamlBytes, err
}

func createYaml(name string, yaml []byte) error {
	// yaml 保存位置
	path := yamlPath + name + "/docker-compose.yml"
	// 确保目录存在
	dirPath := filepath.Dir(path)
	err := utils.MkDir(dirPath)
	if err != nil {
		return err
	}
	// 创建并写入txt文件
	return os.WriteFile(path, yaml, 0644)
}

func getYamlPath(name string) string {
	return yamlPath + name + "/docker-compose.yml"
}
