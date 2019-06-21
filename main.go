package main

import (
	"flag"
	. "fmt"
	. "net/http"
	. "strings"

	. "alopex/app"
	_ "alopex/controller/backend"
)

func main() {
	migrate := *flag.String("migrate", "mysql:games_db.qp", "迁移数据库表数据为模型构造，请填写[数据库类型（mysql/sqlite）:数据库名字]，如：mysql:ckgame")
	flag.Parse()
	if !TT(migrate).IsEmpty() {
		Println("\n开始执行 数据库表->模型构造 迁移任务.....")
		_, err := MD(Replace(migrate, ":", ".", -1)).TM()
		if err != nil {
			DIE(err.Error(), true)
		}
		DIE("执行完成.....", true)
	}
	//
	defer PHandler()
	//
	ELine(1)
	// 加载后端服务
	IsBackendService, _ := String("app").C("is_backend_service")
	if IsBackendService.IsValid() && IsBackendService.IsBool() && TValue(IsBackendService).(bool) {
		String("backend").RH()
	}
	// 监听服务端口
	if err := ListenAndServe(":81", nil); err != nil {
		DIE("服务启动异常，" + err.Error())
	}

}
