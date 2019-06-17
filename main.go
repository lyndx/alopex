package main

import (
	"reflect"

	"alopex/app"

	"github.com/go-ffmt/ffmt"
)

func main() {
	a, _ := app.String("backend").R("auth")
	ffmt.Puts(reflect.Value(a).Interface())
}
