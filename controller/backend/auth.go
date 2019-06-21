package backend

import (
	"alopex/app"
)

type AuthController struct{}

func init() {
	app.CJoin("auth", AuthController{})
}

func (ctrl AuthController) Login(h *app.Http) {
	a, _ := app.MD("mysql.games_db.qp").Select("users")
	h.Output(200, a)
}
