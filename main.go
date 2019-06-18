package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"alopex/app"
)

func main() {
	defer app.PHandler()
	http.HandleFunc("/sdfsd", func(rep http.ResponseWriter, req *http.Request) {
		defer app.PHandler()
		h := app.NT(rep, req)

		h.Verify(nil)

		h.Output(402, "请求失败")
		h.Output(402)
		h.Output(200, *(h.Params), "草组成")

		fmt.Println(time.Now().Unix())
	})

	err := http.ListenAndServe(":81", nil)
	if err != nil {
		fmt.Println("[ERROR] 服务启动异常，" + err.Error())
		os.Exit(1)
	}

}
