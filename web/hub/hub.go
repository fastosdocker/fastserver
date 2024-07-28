package hub

import (
	"fast/config"
	"fast/web/base"
)

// Add 添加docker hub账号
func Add(c *base.Ctx) {
	name := c.PostForm("name")
	username := c.PostForm("username")
	password := c.PostForm("password")
	url := c.PostForm("url")
	authentication := c.PostForm("authentication")
	if name == "" || username == "" || password == "" || url == "" || authentication == "" {
		c.Error("参数错误")
		return
	}

	for _, v := range config.Conf.DockerHubUser {
		if v["name"] == name {
			c.Success()
			return
		}
	}

	config.Conf.DockerHubUser = append(config.Conf.DockerHubUser, map[string]string{"name": name, "username": username, "password": password, "url": url, "authentication": authentication})

	if err := config.Save(); err != nil {
		c.Error(err.Error())
		return
	}
	c.Success()
}

// Del 删除docker hub账号
func Del(c *base.Ctx) {
	name := c.PostForm("name")
	if name == "" {
		c.Error("参数错误")
		return
	}

	for k, v := range config.Conf.DockerHubUser {
		if v["name"] == name {
			config.Conf.DockerHubUser = append(config.Conf.DockerHubUser[:k], config.Conf.DockerHubUser[k+1:]...)
			break
		}
	}

	if err := config.Save(); err != nil {
		c.Error(err.Error())
		return
	}

	c.Success()
}

// Save 修改docker hub账号
func Save(c *base.Ctx) {
	name := c.PostForm("name")
	username := c.PostForm("username")
	password := c.PostForm("password")
	url := c.PostForm("url")
	authentication := c.PostForm("authentication")
	if name == "" || username == "" || password == "" || url == "" || authentication == "" {
		c.Error("参数错误")
		return
	}

	for _, v := range config.Conf.DockerHubUser {
		if v["name"] == name {
			v["username"] = username
			v["password"] = password
			v["url"] = url
			v["authentication"] = authentication
			break
		}
	}

	if err := config.Save(); err != nil {
		c.Error(err.Error())
		return
	}

	c.Success()
}

// Find 查询docker hub账号列表
func Find(c *base.Ctx) {
	if len(config.Conf.DockerHubUser) == 0 {
		c.Success([]string{})
		return
	}

	c.Success(config.Conf.DockerHubUser)
}
