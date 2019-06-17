package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
)

type Http struct {
	rp     http.ResponseWriter
	rq     *http.Request
	params *map[string]interface{}
	errors []string
}

// 获取参数清单
func (h *Http) param() {
	// 参数清单
	params := make(map[string]interface{})
	// Post表单
	h.rq.PostFormValue("")
	for k, v := range h.rq.PostForm {
		params[k] = v
		if len(v) == 1 {
			params[k] = v[0]
		}
	}
	// 常规数据
	h.rq.ParseForm()
	for k, v := range h.rq.Form {
		params[k] = v
		if len(v) == 1 {
			params[k] = v[0]
		}
	}
	// JSON数据
	bs, _ := ioutil.ReadAll(h.rq.Body)
	if len(bs) > 0 {
		tmp := make(map[string]interface{})
		if json.Unmarshal(bs, &tmp) == nil {
			for k, v := range tmp {
				params[k] = v
			}
		}
	}
	// 文件数据
	if h.rq.MultipartForm == nil {
		h.rq.ParseMultipartForm(32 << 20)
	}
	if (h.rq.MultipartForm != nil) && (h.rq.MultipartForm.File != nil) {
		for k, v := range h.rq.MultipartForm.File {
			if len(v) == 1 {
				params[k] = v[0]
			} else if len(v) > 1 {
				params[k] = v
			}
		}
	}
	// 头部数据
	for k, v := range h.rq.Header {
		params[k] = v
		if len(v) == 1 {
			params[k] = v[0]
		}
	}
	h.params = &params
}

// 参数校验
func (h *Http) verify() error {
	result := make(map[string]interface{})
	isTrue, messages := true, []string{}

	for _, cfg := range *h.params {
		// 字段名称
		field := cfg["field"].(string)
		// 校验规则
		rules := make([]string, 0)
		if _, ok := cfg["rules"]; ok && (reflect.TypeOf(cfg["rules"]).String() == "string") && (cfg["rules"] != "") {
			rules = strings.Split(cfg["rules"].(string), " ")
		}
		// 参数是否必填判断
		if _, ok := fields[field]; (!ok) || (fields[field] == nil) || (fields[field] == "") {
			fields[field] = ""
			if _, ok = cfg["default"]; ok {
				fields[field] = cfg["default"]
			}
		}
		// 参数正则判断
		value, isOk, msg := CheckField(fields[field], rules)
		if !isOk {
			isTrue = false
			if _, ok := cfg["label"]; ok {
				messages = append(messages, cfg["label"].(string)+msg)
			}
		} else {
			result[field] = value
		}
	}
	c.Values().Set("__", result)
	if !isTrue {
		Output(&c, strings.Join(messages, "；"), "请求失败", ERROR, nil)
	} else {
		c.Next()
	}
	return nil
}

func init() {

	var xx = func(writer http.ResponseWriter, request *http.Request) {

		writer.Write([]byte("sfdsfsdfsdfs"))
	}
	http.HandleFunc("/", xx)
	err := http.ListenAndServe(":81", nil)
	if err != nil {
		fmt.Println("[ERROR] 服务启动异常，" + err.Error())
		os.Exit(1)
	}
}
