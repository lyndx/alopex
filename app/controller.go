package app

import (
	"fmt"

	"github.com/go-ffmt/ffmt"
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

func CTodo(h *Http, path string, action string) {
	fmt.Println(path, controllerMapper[path], "------------------")
	ffmt.Puts(controllerMapper)
	fmt.Println("---------")
}
