package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"alopex/app"
)

func main() {
	http.HandleFunc("/sdfsd", func(rep http.ResponseWriter, req *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println("[BUG] ", e)
			}
		}()

		h := app.Http{time.Now().UnixNano(), rep, req, nil, nil}
		h.GetParams().Verify(nil)

		h.Output(app.HttpError{"402", "请求失败", "", *(h.Params)})
		h.Output(app.HttpError{"402", "sdfsd", "", *(h.Params)})
		h.Output(app.HttpError{"402", "位34通过", "", *(h.Params)})

		fmt.Println(time.Now().Unix())
	})

	err := http.ListenAndServe(":81", nil)
	if err != nil {
		fmt.Println("[ERROR] 服务启动异常，" + err.Error())
		os.Exit(1)
	}

}
