package backend

import (
	"fmt"

	"alopex/app"
)

type AuthController struct{}

func init() {
	app.CJoin("auth", AuthController{})
}

func (ctrl AuthController) Login(h *app.Http) {
	a, _ := app.MD("main").Select("users u left join user_details d on u.id=d.user_id", true, "u.*,d.balance,d.brief", "u.id <23")
	b, c, e := app.MD("main").Change("add", "users(username,password,nickname)", "'sdfsd','sdfsd','sdfsdf'")
	fmt.Println(b, c, e)
	h.Output(200, a)
}
