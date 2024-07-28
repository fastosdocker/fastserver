package store

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

func extractAppInfo(mdContent string) map[string]string {
	// 定义正则表达式模式
	re := regexp.MustCompile(`应用名称: (.+)\n版本号: (.+)\n图标地址: (.+)\n应用描述: (.+)\n分组名称: (.+)\n`)
	strings.Replace(mdContent, "： ", ": ", 0)
	// 匹配正则表达式
	matches := re.FindStringSubmatch(mdContent)

	// 构造应用信息的映射
	appInfo := make(map[string]string)
	if len(matches) == 6 {
		appInfo["应用名称"] = matches[1]
		appInfo["版本号"] = matches[2]
		appInfo["图标地址"] = matches[3]
		appInfo["应用描述"] = matches[4]
		appInfo["分组名称"] = matches[5]
	}

	return appInfo
}

// 使用正则表达式提取YAML代码块
func extractYAMLCode(mdContent string) string {
	// 通过正则表达式匹配```yml和```之间的内容
	re := regexp.MustCompile("(?s)```yaml(.*?)```")
	matches := re.FindStringSubmatch(mdContent)

	if len(matches) >= 2 {
		return matches[1]
	}

	return ""
}
func readYAML(path string) DockerCompose {
	var compose DockerCompose
	// 读取原始的 YAML 文件
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(path)
		log.Println(err)
		return compose
	}

	// 解析 YAML

	err = yaml.Unmarshal(data, &compose)
	if err != nil {
		fmt.Println(path)
		log.Println(err)
		return compose
	}
	return compose
	// 将修改后的数据结构重新编码为 YAML
	// newData, err := yaml.Marshal(&compose)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }

	// // 写入到新的 YAML 文件
	// err = ioutil.WriteFile("new-docker-compose.yaml", newData, 0644)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
}
func one(fileName string, id int, yamlPath string) (JSONData, error) {
	var data = JSONData{}
	mdContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error reading markdown file:", err)
		return data, err
	}
	// 使用正则表达式提取应用信息
	appInfo := extractAppInfo(string(mdContent))
	if len(appInfo["应用名称"]) > 0 {

		data.Name = appInfo["应用名称"][:len(appInfo["应用名称"])-1]
	} else {
		return data, errors.New("yaml解析错误")
	}
	if len(appInfo["应用描述"]) > 0 {
		data.Description = appInfo["应用描述"][:len(appInfo["应用描述"])-1]
	}
	if len(appInfo["版本号"]) > 0 {
		data.V = appInfo["版本号"][:len(appInfo["版本号"])-1]
	}
	if len(appInfo["图标地址"]) > 0 {
		data.Image = appInfo["图标地址"][:len(appInfo["图标地址"])-1]
	}
	if len(appInfo["分组名称"]) > 0 {
		data.Class = appInfo["分组名称"][:len(appInfo["分组名称"])-1]
	}
	data.Compose = readYAML(yamlPath)
	if data.Compose.Version == "" {
		fmt.Println(fileName)
	}

	data.ID = id

	return data, err
}
func yamlAddNote(path string) error {
	byYaml, _ := ioutil.ReadFile(path + "/docker-compose.yaml")
	byMd, _ := ioutil.ReadFile(path + "/README.md")
	note := extractYAMLCode(string(byMd))
	write := string(byYaml) + "\n" + note
	err := ioutil.WriteFile(path+"/.docker-compose.bak", []byte(write), 777)
	return err
}
