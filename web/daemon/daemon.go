package daemon

import (
	"fast/config"
	"fast/utils"
	"fast/web/base"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Save 修改docker中的daemon.json配置文件
func Save(c *base.Ctx) {
	var (
		err error
	)

	content := strings.TrimSpace(c.PostForm("content"))
	if content == "" {
		c.Error("内容不能为空")
		return
	}

	dir := filepath.Dir(config.Conf.DaemonPath)

	if ok, _ := utils.FileExists(dir); !ok {
		err = utils.MkDir(dir)
		if err != nil {
			c.Error(err.Error())
			return
		}
	}

	err = utils.SaveFile(config.Conf.DaemonPath, content)
	if err != nil {
		c.Error(err.Error())
		return
	}

	c.Success()
}

// Load 获取docker中的daemon.json配置文件
func Load(c *base.Ctx) {
	fileObj, err := os.OpenFile(config.Conf.DaemonPath, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		c.Error(err.Error())
		return
	}

	defer fileObj.Close()

	data, err := ioutil.ReadAll(fileObj)
	if err != nil {
		c.Error(err.Error())
		return
	}

	c.Success(string(data))
}
