package app

import (
	"fast/config"
	"fast/pkg/http"
	"fast/utils"
	"fast/web/base"
	"fmt"
	"os"
	"os/exec"
)

// GetVersion 测试服务是否正常
func GetVersion(c *base.Ctx) {
	rsp := Response{}

	err := http.Get(config.ServerBaseUrl+"getversion", &rsp)
	if err != nil {
		c.Error(err.Error())
		return
	}

	if rsp.Code == 0 {
		c.Error(rsp.Msg)
		return
	}

	rsp.Data.CurVersion = config.Version
	if config.Version < rsp.Data.NewVersion {
		rsp.Data.HasUpdate = true
	}

	c.Success(rsp.Data)
}

// Update 下载更新
func Update(c *base.Ctx) {
	rsp := Response{}

	err := http.Get(config.ServerBaseUrl+"getversion", &rsp)
	if err != nil {
		c.Error(err.Error())
		return
	}

	if rsp.Code == 0 {
		c.Error(rsp.Msg)
		return
	}

	rsp.Data.CurVersion = config.Version
	if config.Version < rsp.Data.NewVersion {
		data := ""
		err = http.Get(rsp.Data.NewUrl, &data)
		if err != nil {
			c.Error(err.Error())
			return
		}

		err = utils.SaveFile("./new_dockercurl", data)
		if err != nil {
			c.Error(err.Error())
			return
		}

		c.Success()
		return
	}

	c.Error("没有可用更新")
}

// Reload 重启服务
func Reload(c *base.Ctx) {
	cmd := exec.Command("kill", "-1", fmt.Sprintf("%d", os.Getpid()))
	err := cmd.Run()
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success()
}

// StartTls 启用或停用tls
func StartTls(c *base.Ctx) {
	var (
		reqJson TlsValidator
		typ     string
		stat    int
	)

	if err := c.ShouldBindJSON(&reqJson); err != nil {
		c.Error(err.Error())
		return
	}

	if err := reqJson.Valid(); err != nil {
		c.Error(err.Error())
		return
	}

	if err := reqJson.Save(); err != nil {
		c.Error(err.Error())
		return
	}

	status := config.Conf.Https.Flag

	if status {
		stat = 1
		if utils.HasPrefix(config.Conf.TLs.Key, "config/server-sign/") {
			typ = "server-sign"
		} else {
			typ = "user-sign"
		}
	} else {
		stat = 0
		typ = "none-sign"
	}

	c.Success(map[string]interface{}{"stat": stat, "typ": typ})
}

// GetTlsStatus 获取tls状态
func GetTlsStatus(c *base.Ctx) {
	var typ string
	var stat int

	status := config.Conf.Https.Flag

	if status {
		stat = 1
		if utils.HasPrefix(config.Conf.TLs.Key, "config/server-sign/") {
			typ = "server-sign"
		} else {
			typ = "user-sign"
		}
	} else {
		stat = 0
		typ = "none-sign"
	}

	c.Success(map[string]interface{}{"stat": stat, "typ": typ})
}

// UploadTls 上传用户签发的tls文件
func UploadTls(c *base.Ctx) {
	keyFile, _ := c.FormFile("keyFile")
	crtFile, _ := c.FormFile("crtFile")
	if keyFile == nil || crtFile == nil {
		c.Error("key文件或crt文件不存在,请重新上传")
		return
	}

	if !utils.HasSuffix(keyFile.Filename, ".key") && !utils.HasSuffix(crtFile.Filename, ".crt") {
		c.Error("key文件或crt文件格式有误,请重新上传")
		return
	}

	keyDst := fmt.Sprintf("config/user-sign/%s", keyFile.Filename)
	crtDst := fmt.Sprintf("config/user-sign/%s", crtFile.Filename)

	// 上传文件至指定的完整文件路径
	if err := c.SaveUploadedFile(keyFile, keyDst); err != nil {
		c.Error(fmt.Sprintf("保存key文件失败,失败原因 %s", err))
		return
	}

	if err := c.SaveUploadedFile(crtFile, crtDst); err != nil {
		c.Error(fmt.Sprintf("保存crt文件失败,失败原因 %s", err))
		return
	}

	c.Success(map[string]interface{}{"crtFile": crtDst, "keyFile": keyDst})
}
