package compose

import (
	"errors"
	"fast/pkg/compose"
	"fast/pkg/db"
	"fast/web/base"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func Run(c *base.Ctx) {
	id := c.Query("id")

	if id == "" {
		c.Error("id不能为空")

		return
	}
	orphans := c.Query("orphans")

	app := db.Application{}
	err := db.FindOne(&app, db.Query{Where: fmt.Sprintf("`id`=%s", id)})
	if err != nil {
		c.Error(fmt.Sprintf("查询失败: %s", err))
		return
	}

	if app.ID == 0 {
		c.Error("应用不存在")
		return
	}

	err = compose.Start(app.Name, app.Compose, orphans)
	if err != nil {
		c.Error(err.Error())
		return
	}

	c.Success()
}

// Create 创建项目
func Create(c *base.Ctx) {
	app := db.Application{BaseModel: &db.BaseModel{}}
	app.Name = c.PostForm("name")
	app.Compose = c.PostForm("compose")
	app.Description = c.PostForm("description")
	app.V = c.PostForm("v")
	app.Image = c.PostForm("image")
	app.Note = c.PostForm("note")
	orphans := c.PostForm("orphans")

	if app.Name == "" {
		c.Error("项目名称不能为空")
		return
	}
	if app.Compose == "" {
		c.Error("yaml数据不能为空")
		return
	}
	if app.Image == "" {
		c.Error("图标不能为空")
		return
	}

	app.Name = strings.ToLower(app.Name)

	appLog := db.ApplicationLog{}
	_ = c.ShouldBind(&appLog)

	err := db.Get().Transaction(func(tx *gorm.DB) error {
		myApp := db.Application{BaseModel: &db.BaseModel{}}
		err := db.FindOne(&myApp, db.Query{Where: fmt.Sprintf("`name`='%s'", app.Name)})
		if err != nil {
			return err
		}

		if myApp.ID == 0 {
			err = db.Create(&app)
			if err != nil {
				return err
			}
		} else {
			app.ID = myApp.ID
			err = db.Save(app)
			if err != nil {
				return err
			}
		}

		// 插入日志
		appLog.BaseModel = &db.BaseModel{}
		err = db.Create(&appLog)
		if err != nil {
			return err
		}

		if app.BaseModel == nil || app.ID == -1 {
			return errors.New("请添加新的版本号")
		}

		if err != nil {
			return err
		}

		_ = compose.Start(app.Name, app.Compose, orphans)

		return nil
	})
	if err != nil {
		c.Error(err.Error())
		return
	}

	c.Success()
}

// Log 日志
func Log(c *base.Ctx) {
	name := c.Query("name")
	if name == "" {
		c.Error("参数错误,应用名称不能为空")
		return
	}

	imageName := c.Query("imageName")
	rows := c.DefaultQuery("rows", "1000")
	var services []string
	if imageName != "" {
		services = append(services, imageName)
	}
	project, err := compose.Log(name, rows, services)
	if err != nil {
		c.Error(err.Error())
		return
	}

	c.Success(project)
}

// Stop 停止项目
func Stop(c *base.Ctx) {
	name := c.Query("name")
	if name == "" {
		c.Error("参数错误,用用名称不能为空")
		return
	}

	err := compose.Stop(name)
	if err != nil {
		c.Error(err.Error())
		return
	}

	c.Success()
}

// Del 删除容器
func Del(c *base.Ctx) {
	id := c.Query("id")
	if id == "" {
		c.Error("参数错误,ID不能为空")
		return
	}

	app := db.Application{BaseModel: &db.BaseModel{}}
	err := db.FindOne(&app, db.Query{Where: fmt.Sprintf("`id`=%s", id)})
	if err != nil {
		c.Error(fmt.Sprintf("查询失败: %s", err))
		return
	}

	if app.ID == 0 {
		c.Error("应用不存在")
		return
	}

	err = db.Get().Transaction(func(tx *gorm.DB) error {
		err = db.Delete(&app, fmt.Sprintf("`id`=%d", app.ID))
		if err != nil {
			c.Error(err.Error())
			return err
		}

		err = compose.Remove(app.Name)
		if err != nil {
			c.Error(err.Error())
		}

		return nil
	})
	if err != nil {
		c.Error("删除应用失败: " + err.Error())
		return
	}

	c.Success()
}

// Update 更新容器
func Update(c *base.Ctx) {
	app := db.Application{}
	err := c.BindJSON(&app)
	if err != nil {
		// 容器更新失败
		c.Error("数据格式不正确" + err.Error())
		return
	}
	if app.ID == 0 {
		c.Error("缺少id")
		return
	}
	if app.Compose == "" {
		c.Error("缺少Compose数据")
		return
	}
	if app.Name == "" {
		c.Error("缺少Name数据")
		return
	}
	if app.V == "" {
		c.Error("版本号不能为空")
		return
	}

	myApp := db.Application{BaseModel: &db.BaseModel{}}
	err = db.FindOne(&myApp, db.Query{Where: fmt.Sprintf("`id`=%d", app.ID)})
	if err != nil {
		c.Error(fmt.Sprintf("查询失败: %s", err))
		return
	}

	if myApp.ID == 0 {
		c.Error("应用不存在")
		return
	}

	if myApp.V == app.V {
		c.Error("该版本号已经存在，请输入新的版本号")
		return
	}

	err = db.Get().Transaction(func(tx *gorm.DB) error {
		err = db.Save(app)
		if err != nil {
			c.Error(err.Error())
			return err
		}

		// 更新容器
		err = compose.Start(app.Name, app.Compose, "")
		if err != nil {
			c.Error(err.Error())
			return err
		}

		return nil
	})
	if err != nil {
		c.Error("更新应用失败: " + err.Error())
		return
	}

	c.Success()
}

// Get 获取项目信息
func Get(c *base.Ctx) {
	name := c.Query("name")
	if name == "" {
		c.Error("参数错误,应用名称不能为空")
		return
	}

	get, err := compose.Get(name)
	if err != nil {
		c.Error(err.Error())
		return
	}

	c.Success(get)
}

// List 获取项目列表
func List(c *base.Ctx) {
	// apps := make([]*db.Application, 0)
	// // 获取数据，将数据转为map
	// err := db.Find(&apps)
	// if err != nil {
	// 	c.Error(fmt.Sprintf("查询失败: %s", err))
	// 	return
	// }

	// appMap := make(map[string]*db.Application)
	// for i := range apps {
	// 	appMap[apps[i].Name] = apps[i]
	// }
	// if len(apps) == 0 {
	// 	c.Success(apps)
	// 	return
	// }

	// // 数据库数据跟容器数据同步
	// list, err := compose.List()
	// if err != nil {
	// 	c.Error(fmt.Sprintf("获取项目列表失败: %s", err))
	// 	return
	// }
	// //  只负责程序生成的容器
	// for i := 0; i < len(list); i++ {
	// 	if appMap[list[i].Name] != (&db.Application{}) {
	// 		m := appMap[list[i].Name]
	// 		if m == nil {
	// 			continue
	// 		}
	// 		m.Status = list[i].Status
	// 	}
	// }

	// c.Success(apps)

	pageTem := c.DefaultQuery("page", "1")
	sizeTem := c.DefaultQuery("size", "20")
	page, _ := strconv.Atoi(pageTem)
	size, _ := strconv.Atoi(sizeTem)
	apps := make([]*db.Application, 0)
	res := struct {
		App     []*db.Application `json:"app"`
		Page    int               `json:"page"`
		Size    int               `json:"size"`
		PageAll int               `json:"pageAll"`
	}{}
	// 获取数据，将数据转为map
	err := db.Find(&apps)
	if err != nil {
		c.Error(fmt.Sprintf("查询失败: %s", err))
		return
	}
	pageAll := len(apps)
	sum := 0
	if pageAll == 0 {
		c.Success(res)
		return
	}

	if pageAll%size != 0 {
		sum = 1
	}
	pageAll = pageAll/size + sum
	err = db.Find(&apps, db.Query{Limit: page, Offset: size * page})
	if err != nil {
		c.Error(fmt.Sprintf("查询失败: %s", err))
		return
	}

	appMap := make(map[string]*db.Application)
	for i := range apps {
		appMap[apps[i].Name] = apps[i]
	}

	// 数据库数据跟容器数据同步
	list, err := compose.List()
	if err != nil {
		c.Error(fmt.Sprintf("获取项目列表失败: %s", err))
		return
	}
	//  只负责程序生成的容器
	for i := 0; i < len(list); i++ {
		if appMap[list[i].Name] != (&db.Application{}) {
			m := appMap[list[i].Name]
			if m == nil {
				continue
			}
			m.Status = list[i].Status
		}
	}
	res.App = apps
	res.Page = page
	res.Size = size
	res.PageAll = pageAll
	c.Success(res)
}
