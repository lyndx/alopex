package backend

import (
	"alopex/app"
	"alopex/service"
)

type AuthController struct{}

func init() {
	app.CJoin("auth", AuthController{})
}

func (ctrl AuthController) Login(h *app.Http) {
	us := service.UserService{}
	user, err := us.GetUserByUsername("x")
	if err != nil {
		h.Output(402, "用户信息获取失败")
	}
	as := service.AuthService{}
	tokenStr, randomStr, err := as.GetToken(h.Module, "1")
	if err != nil {
		h.Output(402, "登录失败")
	}
	h.Rep.Header().Set("token", tokenStr)
	h.Rep.Header().Set("random_str", randomStr)
	h.Output(200, user)
}
