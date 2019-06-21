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
	a, e := app.MD("mysql.games_db.qp").Select("users")
	fmt.Println(e)
	h.Output(200, a)
}
