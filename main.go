package main

import (
	"alopex/app"
	_ "alopex/controller/backend"
	"flag"
	"fmt"
	. "net/http"
)

func main() {
	migrate := *flag.String("migrate", "mysql:games_db.qp", "迁移数据库表数据为模型构造，请填写[数据库类型（mysql/sqlite）:数据库名字]，如：mysql:ckgame")
	flag.Parse()
	if !app.TT(migrate).IsEmpty() {
		fmt.Println("\n开始执行 数据库表->模型构造 迁移任务.....")

		app.MD("mysql.maindb").Select("platform")

		//app.DIE("执行完成.....", true)
	}
	fmt.Println("migrate:", migrate)

	defer app.PHandler()
	//
	app.ELine(1)
	// 加载后端服务
	IsBackendService, _ := app.String("app").C("is_backend_service")
	if IsBackendService.IsValid() && IsBackendService.IsBool() && app.TValue(IsBackendService).(bool) {
		app.String("backend").RH()
	}
	// 监听服务端口
	if err := ListenAndServe(":81", nil); err != nil {
		app.DIE("服务启动异常，" + err.Error())
	}

}
