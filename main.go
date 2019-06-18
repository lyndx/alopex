package main

import (
	"fmt"
	"net/http"
	"os"

	"alopex/app"
)

func main() {
	defer app.PHandler()

	// 加载后端服务
	IsBackendService, _ := app.String("app").C("is_backend_service")
	if IsBackendService.IsValid() && IsBackendService.IsBool() && IsBackendService.Value().(bool) {
		app.String("backend").RH()
	}

	// 监听服务端口
	err := http.ListenAndServe(":81", nil)
	if err != nil {
		fmt.Println("[ERROR] 服务启动异常，" + err.Error())
		os.Exit(1)
	}

}
