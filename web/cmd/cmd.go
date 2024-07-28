package cmd

import (
	"context"
	"fast/web/base"
	"os/exec"
	"strings"
)

// Shell 执行Shell命令
func Shell(c *base.Ctx) {
	//获取前端传来的docker参数值，存入d
	d := c.PostForm("docker")
	//通过空格区分参数d，存入b
	b := strings.Split(d, " ")
	//判断b的长度是否大于4
	if len(b) > 4 {
		//通过命令行执行d命令
		cmd := exec.CommandContext(context.Background(), "/bin/sh", "-c", d)
		//将结果存入out
		out, err := cmd.CombinedOutput()
		if err != nil {
			c.Error(err.Error())
			return
		}
		//给前端返回out
		c.Success(string(out))
		return
	}
	c.Error("参数有误")
}
