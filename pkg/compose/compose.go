package compose

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/cmd/formatter"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
)

var (
	yamlPath string
	Svc      api.Service
)

func Init(opt Option) error {
	yamlPath = opt.YamlPath

	err := install(opt.SH)
	if err != nil {
		return err
	}

	return initCli()
}

func install(sh string) error {
	log.Println("检查docker-compose是否安装")
	cmd := exec.Command("docker-compose", "--version")
	// 执行命令并捕获输出
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("检查docker-compose未安装,开始执行安装命令")
		cmd = exec.Command("/bin/sh", sh)
		cmd.Dir = "." // 设置执行命令的工作目录为当前目录

		var output []byte
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("compose安装失败: %s", err)
		}

		log.Println("命令输出：", string(output))
		return nil
	}

	log.Printf("%s, 正在启动...\n", string(out))

	return nil
}

func initCli() error {
	dockerCli, err := command.NewDockerCli()
	if err != nil {
		return fmt.Errorf("创建docker cli失败: %s", err)
	}

	opt := flags.NewClientOptions()
	err = dockerCli.Initialize(opt)
	if err != nil {
		return fmt.Errorf("初始化docker cli失败: %s", err)
	}

	Svc = compose.NewComposeService(dockerCli)

	return nil
}

func Start(name string, jsonStr string, orphans string) error {

	yaml, err := jsonToYaml([]byte(jsonStr))
	if err != nil {
		return err
	}
	fmt.Println(string(yaml))
	err = createYaml(name, yaml)
	if err != nil {
		return fmt.Errorf("创建yaml文件失败: %s", err)
	}

	var cmd *exec.Cmd
	if orphans == "1" {
		cmd = exec.Command("docker-compose", "-f", getYamlPath(name), "up", "-d", "--remove-orphans")
	} else {
		cmd = exec.Command("docker-compose", "-f", getYamlPath(name), "up", "-d")
	}

	// 执行命令并捕获输出
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to create Docker Compose container:", err, string(out))
		return fmt.Errorf("启动项目失败: %s", err)
	}

	fmt.Println(out)

	return err
}

// Get 查询项目信息
func Get(name string) ([]api.ContainerSummary, error) {
	ps, err := Svc.Ps(context.Background(), name, api.PsOptions{})
	if err != nil {
		return nil, fmt.Errorf("查询项目信息: %s", err)
	}

	return ps, nil
}

// Log 获取项目日志
func Log(name string, rows string, services []string) (string, error) {
	// 查询日志命令
	writer := &stringWriter{}
	consumer := formatter.NewLogConsumer(context.Background(), writer, false, false)

	err := Svc.Logs(context.Background(), name, consumer, api.LogOptions{Tail: rows, Services: services})
	if err != nil {
		return "", fmt.Errorf("获取项目日志: %s", err)
	}

	return writer.String(), nil
}

// Stop 停止项目命令
func Stop(name string) error {
	cmd := exec.Command("docker-compose", "-f", getYamlPath(name), "stop")
	// 执行命令并捕获输出
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("停止项目失败: %s", err)
	}

	return nil
}

// Remove 删除项目
func Remove(name string) error {
	err := Stop(name)
	if err != nil {
		return err
	}
	cmd := exec.Command("docker-compose", "-f", getYamlPath(name), "rm", "--force")
	// 执行命令并捕获输出
	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("删除项目失败: %s", err)
	}

	err = os.RemoveAll(filepath.Dir(getYamlPath(name)))
	if err != nil {
		return fmt.Errorf("删除yaml文件失败: %s", err)
	}

	return nil
}

// List 获取项目列表
func List() ([]api.Stack, error) {
	return Svc.List(context.Background(), api.ListOptions{})
}
