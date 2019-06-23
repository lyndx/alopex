package backend

import (
	"alopex/app"
)

type AuthController struct{}

func init() {
	app.CJoin("auth", AuthController{})
}

func (ctrl AuthController) Login(h *app.Http) {
	a, _ := app.MD("main").Select("users u left join user_details d on u.id=d.user_id", "u.*,d.balance,d.brief", "u.id <23")
	h.Output(200, a)
}
