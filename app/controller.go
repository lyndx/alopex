package app

import (
	"reflect"
)

type Controller struct{}

var controllerMapper map[string]interface{}

func init() {
	controllerMapper = make(map[string]interface{})
}

func CJoin(name string, controller interface{}) {
	cv := RV(controller)
	if cv.IsValid() {
		path := String(cv.Type().PkgPath()).Split("/controller/")[1] + "." + name
		controllerMapper[path] = controller
	}
}

func (h *Http) CTodo(module string, controller string, action string) {
	key := module + "." + controller
	if _, ok := controllerMapper[key]; !ok {
		h.Output(402, "请求失败", "请求执行业务方法不存在")
	}
	obj := controllerMapper[key]
	handler := RV(obj).MethodByName(String(action).UFrist())
	if !handler.IsValid() {
		h.Output(402, "请求失败", "请求执行业务方法不存在")
	}
	handler.Call([]reflect.Value{RV(h)})
}
