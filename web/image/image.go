package image

import (
	"context"
	"fast/pkg/docker"
	"fast/utils"
	"fast/web/base"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/filters"
)

// Save 导出镜像
func Save(c *base.Ctx) {
	image := c.PostForm("image")

	imageInspect, _, err := docker.Get().ImageInspectWithRaw(c, image)
	if err != nil {
		c.Error(err.Error())
		return
	}
	r, err := docker.Get().ImageSave(c, []string{image})
	if err != nil {
		c.Error(err.Error())
		return
	}

	file, err := io.ReadAll(r)
	if err != nil {
		c.Error(err.Error())
		return
	}

	filename := imageInspect.RepoTags[len(imageInspect.RepoTags)-1] + ".tar"
	path := "data/image/" + filename
	dirPath := filepath.Dir(path)
	utils.MkDir(dirPath)
	err = os.WriteFile(path, file, os.ModePerm)
	if err != nil {
		c.Error(err.Error())
		return
	}
	defer os.Remove(path)

	c.File(path)
}

// Load 导入镜像
func Load(c *base.Ctx) {
	image, _, err := c.Request.FormFile("image")
	if err != nil {
		c.Error(err.Error())
		return
	}

	_, err = docker.Get().ImageLoad(context.Background(), image, false)
	if err != nil {
		c.Error(err.Error())
		return
	}

	c.Success()
}

// Prune 删除未使用的镜像
func Prune(c *base.Ctx) {
	dangling := c.PostForm("dangling")
	if dangling != "1" && dangling != "0" {
		c.Error("参数错误")
		return
	}

	args := filters.NewArgs(filters.KeyValuePair{Key: "dangling", Value: "0"})

	_, err := docker.Get().ImagesPrune(context.Background(), args)
	if err != nil {
		c.Error(err.Error())
		return
	}

	c.Success()
}
