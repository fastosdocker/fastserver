package user

import (
	"fast/pkg/db"
	"fast/utils"
	"fast/web/mw"
	"fmt"

	"fast/web/base"
)

func Login(r *base.Ctx) {
	username, _ := r.GetPostForm("username")
	password, _ := r.GetPostForm("password")
	if username == "" || password == "" {
		r.Error("登陆信息不能为空")
		return
	}

	user := db.User{BaseModel: &db.BaseModel{}}
	err := db.FindOne(&user, db.Query{Where: fmt.Sprintf("`user_name`='%s' and `password`='%s'", username, utils.Md5(password))})
	if err != nil {
		r.Error(fmt.Sprintf("查询失败: %s", err))
		return
	}

	if user.ID == 0 {
		var count int64
		count, err = db.Count(&user)
		if err != nil {
			r.Error("统计用户失败")
			return
		}

		if count != 0 {
			r.Error("用户名或密码错误")
			return
		}

		user.UserName = username
		user.Password = utils.Md5(password)
		err = db.Create(&user)
		if err != nil {
			r.Error("创建用户失败")
			return
		}
	}

	token, err := mw.GenToken(mw.UserInfo{User: user})
	if err != nil {
		r.Error(err.Error())
		return
	}

	r.Success(token)
}

// Register 注册用户
func Register(c *base.Ctx) {
	req := RegisterReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.Error(fmt.Sprintf("解析参数错误: %s", err))
		return
	}

	if req.UserName == "" {
		c.Error("账号不能为空")
	}

	if req.PassWord == "" {
		c.Error("密码不能为空")
	}

	if len(req.PassWord) <= 6 {
		c.Error("密码必须大于等于六位")
		return
	}

	if req.PassWord != req.PassWord2 {
		c.Error("两次密码不一致")
	}

	user := db.User{BaseModel: &db.BaseModel{}}
	err := db.FindOne(&user, db.Query{Where: fmt.Sprintf("`user_name`='%s'", req.UserName)})
	if err != nil {
		c.Error(fmt.Sprintf("查询用户失败: %s", err))
		return
	}

	if user.ID != 0 {
		c.Error(fmt.Sprintf("[%s]账号已注册", req.UserName))
		return
	}

	user.UserName = req.UserName
	user.Password = utils.Md5(req.PassWord)
	err = db.Create(&user)
	if err != nil {
		c.Error(fmt.Sprintf("注册用户失败: %s", err))
		return
	}

	c.Success()
}

// Password 修改密码
func Password(c *base.Ctx) {
	var (
		req PassWordReq
		err error
	)

	v, exists := c.Get("user")
	if !exists {
		c.Error("请先登陆")
		return
	}

	user := v.(db.User)
	if err = c.ShouldBind(&req); err != nil {
		c.Error(err.Error())
		return
	}

	if req.OldPassword == "" {
		c.Error("密码不能为空")
		return
	}

	if len(req.NewPassword) <= 6 {
		c.Error("新密码必须大于等于六位")
		return
	}

	u := db.User{BaseModel: &db.BaseModel{}}
	err = db.FindOne(&u, db.Query{Where: fmt.Sprintf("`user_name`='%s'", user.UserName)})
	if err != nil {
		c.Error(fmt.Sprintf("查询用户失败: %s", err))
		return
	}

	if u.ID == 0 {
		c.Error(fmt.Sprintf("[%s]账号不存在", user.UserName))
		return
	}

	oldP := utils.Md5(req.OldPassword)
	if oldP != u.Password {
		c.Error("旧密码不正确")
		return
	}

	u.Password = utils.Md5(req.NewPassword)

	err = db.Save(u)
	if err != nil {
		c.Error(err.Error())
		return
	}

	c.Success()
}
