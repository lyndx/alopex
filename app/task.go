package app

import (
	"reflect"

	"github.com/robfig/cron"
)

var Tasks map[string]reflect.Value

func init() {
	Tasks = make(map[string]reflect.Value)
}

func InitTask() {
	tmp, err := String("task").C("list")
	if (err != nil) || (!tmp.IsArray()) {
		DIE("定时任务配置异常")
	}
	list := make(map[string]map[string]interface{})
	rtmp := reflect.Value(tmp)
	for i := 0; i < rtmp.Len(); i++ {
		name := rtmp.Index(i).Interface().(string)
		if _, ok := Tasks[name]; !ok {
			DIE("定时任务配置异常，任务不存在")
		}
		funcSpec := Tasks[name].MethodByName("Spec")
		if !funcSpec.IsValid() {
			DIE("定时任务“" + name + "”的方法Spec不存在")
		}
		spec := funcSpec.Call([]reflect.Value{})[0].Interface().(string)
		funcHandler := Tasks[name].MethodByName("Handler")
		if !funcHandler.IsValid() {
			DIE("定时任务“" + name + "”的方法Spec不存在")
		}
		list[name] = map[string]interface{}{
			"spec":    spec,
			"handler": funcHandler.Interface().(func()),
		}
	}
	if len(list) > 0 {
		c := cron.New()
		for name, value := range list {
			c.AddFunc(value["spec"].(string), value["handler"].(func()))
			Dump("green", "定时任务“"+name+"”启动成功")
		}
		c.Start()
	}
}
