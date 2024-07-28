package store

import (
	"encoding/json"
	"fast/web/base"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

var rspForever = new(ApplicationRsp)

var classMap = make(map[string][]JSONData)

func init() {
	var data = []JSONData{}
	cmd := exec.Command("git", "clone", "http://192.168.0.237/fastos/appstore.git")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	dirName, err := os.ReadDir("./appstore")
	if err != nil {
		fmt.Println("readDir::" + err.Error())
		return
	}
	for i, j := range dirName {
		if !j.IsDir() {
			continue
		}
		if j.Name() == ".git" {
			continue
		}

		path := "./appstore/" + j.Name() + "/README.md"
		yamlPath := "./appstore/" + j.Name() + "/docker-compose.yaml"
		if TEM, err := one(path, i, yamlPath); err != nil {
			fmt.Println("yaml解析错误::" + j.Name())
			continue
		} else {
			data = append(data, TEM)
		}

		// 将数据编码为JSON并写入文件
		err = yamlAddNote("./appstore/" + j.Name())
		if err != nil {
			fmt.Println("写入yamlNote数据失败:", err)
			return
		}
	}

	for _, j := range data {
		classMap[j.Class] = append(classMap[j.Class], j)
		fmt.Println(j.Class)
	}

	file, err := os.Create("./jsonData.json")
	if err != nil {
		fmt.Println("创建文件失败:", err)
		return
	}
	defer file.Close()

	// 将数据编码为JSON并写入文件
	encoder := json.NewEncoder(file)
	err = encoder.Encode(classMap)
	if err != nil {
		fmt.Println("写入JSON数据失败:", err)
		return
	}
}
func PageId(c *base.Ctx) {
	class := c.Query("class")
	idString := c.Query("id")
	id, _ := strconv.Atoi(idString)
	for _, j := range classMap[class] {
		if j.ID == id {
			c.Success(j)
			return
		}
	}
	c.Error("id不存在")
	return
}

// 新增 我的应用
// 修改 我的应用
// 删除 我的应用
func Page(c *base.Ctx) {
	className := c.Query("class")
	re := classMap[className]
	var getId []JSONDataGetId
	for _, j := range re {
		tem := JSONDataGetId{}
		tem.Class = j.Class
		tem.Description = j.Description
		tem.ID = j.ID
		tem.Image = j.Image
		tem.Name = j.Name
		tem.V = j.V
		getId = append(getId, tem)
	}

	c.Success(getId)
	return
}
func Class(c *base.Ctx) {
	var className []string
	for k, _ := range classMap {
		className = append(className, k)
	}
	c.Success(className)
	return
}
func Update(c *base.Ctx) {
	classString := c.PostForm("class")
	idString := c.PostForm("id")
	compose := c.PostForm("compose")
	id, _ := strconv.Atoi(idString)
	for i, j := range classMap[classString] {
		if j.ID == id {
			comYaml := []byte(compose)
			err := ioutil.WriteFile("./appstore/"+j.Name+"/docker-compose.yaml", comYaml, 777)
			if err != nil {
				fmt.Println(err)
				c.Error("writeFile error")
				return
			}
			dockerCompose := readYAML("./appstore/" + j.Name + "/docker-compose.yaml")
			classMap[classString][i].Compose = dockerCompose
		}
	}
	c.Success("ok")
	return
}
