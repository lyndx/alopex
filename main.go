package main

import (
	"flag"
	. "net/http"

	. "alopex/app"
	_ "alopex/controller/backend"
	_ "alopex/model"
)

func main() {
	D := flag.String("d", "", "操作的数据库名字")
	T := flag.String("t", "", "迁移方式：m2t->数据模型构造为数据表；t2m->数据表为数据模型构造")
	flag.Parse()
	d, t := *D, *T
	if (d != "") && (t != "") {
		Dump("green", "开始执行 迁移任务.....")
		D := MD(d)
		var err error = nil
		if t == "m2t" {
			_, err = D.MT()
		} else if t == "t2m" {
			_, err = D.TM()
		}
		if err != nil {
			DIE(err.Error())
		}
		DIE("执行完成.....", true)
	}
	//
	defer PHandler()
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
