package task

import (
	"alopex/app"
	"fmt"
)

type TestTask struct{}

func init() {
	app.Tasks["test"] = app.RV(TestTask{})
}

func (t TestTask) Spec() string {
	return "*/5 * * * * ?"
}

func (t TestTask) Handler() {
	fmt.Println("------------")
}
