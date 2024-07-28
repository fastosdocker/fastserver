package app

import (
	"errors"
	"fast/config"
	"fast/utils"
	"fast/web/base"
	"fmt"
	"os"
)

type Response struct {
	base.Response
	Data struct {
		NewVersion string
		NewUrl     string
		CurVersion string
		HasUpdate  bool
	}
}

type TlsValidator struct {
	Typ string `json:"typ"`
	Crt string `json:"crt"`
	Key string `json:"key"`
}

func (t *TlsValidator) Valid() error {
	if t.Typ == "none-sign" {
		if !config.Conf.Https.Flag {
			return errors.New("当前服务状态为none-sign, 请不要重复操作")
		}
		return nil
	} else if t.Typ == "server-sign" {
		if config.Conf.Https.Flag && utils.HasPrefix(config.Conf.TLs.Key, "config/server-sign/") {
			return errors.New("当前服务状态为server-sign, 请不要重复操作")
		}
		return nil
	} else if t.Typ == "user-sign" {
		if config.Conf.Https.Flag && utils.HasPrefix(config.Conf.TLs.Key, "config/user-sign/") {
			return errors.New("当前服务状态为user-sign, 请不要重复操作")
		}
		// 验证crt和key
		if t.Key == "" || t.Crt == "" {
			return errors.New("缺少crt或key")
		}
		if !utils.HasSuffix(t.Crt, ".crt") {
			return errors.New("crt文件名错误")
		}
		if !utils.HasSuffix(t.Key, ".key") {
			return errors.New("key文件名称错误")
		}
		_, err := os.Stat(t.Key)
		if err != nil {
			return errors.New("key文件不存在，请重新上传")
		}
		_, err = os.Stat(t.Crt)
		if err != nil {
			return errors.New("crt文件不存在，请重新上传")
		}
		return nil
	} else {
		return errors.New("未识别的typ，请重新输入")
	}
}

func (t *TlsValidator) Save() error {
	if t.Typ == "none-sign" {
		// 不签发，http
		// 修改配置信息
		config.Conf.Https.Flag = false
		config.Conf.TLs.Key = ""
		config.Conf.TLs.Pem = ""
	} else if t.Typ == "server-sign" {
		// 系统自签，https
		config.Conf.Https.Flag = true
		// 生成证书
		err, keyPath, crtPath := utils.CreateTlsFile("config/server-sign")
		if err != nil {
			return errors.New(fmt.Sprintf("生成证书失败，%s", err.Error()))
		}
		config.Conf.TLs.Key = keyPath.(string)
		config.Conf.TLs.Pem = crtPath.(string)
	} else if t.Typ == "user-sign" {
		// 用户自签， https
		config.Conf.Https.Flag = true
		config.Conf.TLs.Key = t.Key
		config.Conf.TLs.Pem = t.Crt
	}
	// 保存配置信息
	if err := config.Save(); err != nil {
		return errors.New("写入配置文件失败：" + err.Error())
	}
	return nil
}
