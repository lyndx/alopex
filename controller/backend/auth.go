package backend

import "alopex/app"

type AuthController struct {
}

func init() {
	app.CJoin("auth", AuthController{})
}

func (ctrl AuthController) Login(h *app.Http) {

	h.Output(444, "xddd")
}
