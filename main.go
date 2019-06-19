package main

import (
	"net/http"

	"alopex/app"
)

func main() {
	//
	defer app.PHandler()
	//
	app.ELine(1)
	// 加载后端服务
	IsBackendService, _ := app.String("app").C("is_backend_service")
	if IsBackendService.IsValid() && IsBackendService.IsBool() && app.TValue(IsBackendService).(bool) {
		app.String("backend").RH()
	}
	// 监听服务端口
	if err := http.ListenAndServe(":81", nil); err != nil {
		app.DIE("服务启动异常，" + err.Error())
	}

}
