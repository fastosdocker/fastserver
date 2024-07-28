package user

type RegisterReq struct {
	UserName  string `json:"username" form:"username"`
	PassWord  string `json:"password" form:"password"`
	PassWord2 string `json:"password2" form:"password2"`
}

type PassWordReq struct {
	OldPassword string `json:"old_password" form:"old_password"`
	NewPassword string `json:"new_password" form:"new_password"`
}
