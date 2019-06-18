package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type Http struct {
	STime  time.Time
	Rep    http.ResponseWriter
	Req    *http.Request
	Params *map[string]interface{}
}

// 创建Http实例
func NT(rep http.ResponseWriter, req *http.Request) *Http {
	h := new(Http)
	h.Rep = rep
	h.Req = req
	h.STime = time.Now()
	h.Params = h.params()
	return h
}

// 获取参数清单
func (h *Http) params() *map[string]interface{} {
	// 参数清单
	params := make(map[string]interface{})
	// Post表单
	h.Req.PostFormValue("")
	for k, v := range h.Req.PostForm {
		params[k] = v
		if len(v) == 1 {
			params[k] = v[0]
		}
	}
	// 常规数据
	h.Req.ParseForm()
	for k, v := range h.Req.Form {
		params[k] = v
		if len(v) == 1 {
			params[k] = v[0]
		}
	}
	// JSON数据
	bs, _ := ioutil.ReadAll(h.Req.Body)
	if len(bs) > 0 {
		tmp := make(map[string]interface{})
		if json.Unmarshal(bs, &tmp) == nil {
			for k, v := range tmp {
				params[k] = v
			}
		}
	}
	// 文件数据
	if h.Req.MultipartForm == nil {
		h.Req.ParseMultipartForm(32 << 20)
	}
	if (h.Req.MultipartForm != nil) && (h.Req.MultipartForm.File != nil) {
		for k, v := range h.Req.MultipartForm.File {
			if len(v) == 1 {
				params[k] = v[0]
			} else if len(v) > 1 {
				params[k] = v
			}
		}
	}
	// 头部数据
	for k, v := range h.Req.Header {
		k = strings.ToLower(k)
		params["_"+k] = v
		if len(v) == 1 {
			params["_"+k] = v[0]
		}
	}
	return &params
}

// 字段校验
func (h *Http) checkField(field interface{}, rules []string) (interface{}, bool, string) {
	Field, IsTrue, MSG := TT(field), TT(true), TT("")
	IsValid := Field.IsValid()
	for _, rule := range rules {
		switch rule {
		case "file":
			IsTrue.SwitchValue(IsValid, Field.IsFile(false), true)
			MSG.SwitchValue(reflect.Value(IsTrue).Bool(), "", "必须为文件")
		case "files":
			IsTrue.SwitchValue(IsValid, Field.IsFile(true), true)
			MSG.SwitchValue(reflect.Value(IsTrue).Bool(), "", "必须为文件数组")
		case "must":
			IsTrue.SwitchValue(IsValid, Field.IsEmpty(), true)
			MSG.SwitchValue(reflect.Value(IsTrue).Bool(), "", "为必填字段")
		case "string":
			IsTrue.SwitchValue(IsValid, Field.IsString(), true)
			MSG.SwitchValue(reflect.Value(IsTrue).Bool(), "", "必须为字符串格式")
		case "int":
			Field = TT(Field.ToString())
			IsTrue.SwitchValue(IsValid, Field.IsInt(), true)
			MSG.SwitchValue(reflect.Value(IsTrue).Bool(), "", "必须为整数")
		case "float":
			Field = TT(Field.ToString())
			IsTrue.SwitchValue(IsValid, Field.IsFloat(), true)
			MSG.SwitchValue(reflect.Value(IsTrue).Bool(), "", "必须为浮点数")
		case "bool":
			Field = TT(Field.ToString())
			IsTrue.SwitchValue(IsValid, Field.IsBool(), true)
			MSG.SwitchValue(reflect.Value(IsTrue).Bool(), "", "必须为布尔值")
		case "array":
			IsTrue.SwitchValue(IsValid, Field.IsArray(), true)
			MSG.SwitchValue(reflect.Value(IsTrue).Bool(), "", "必须为数组")
		default:
			if IsValid {
				match := false
				if Field.IsString() {
					match, _ = regexp.MatchString(rule, Field.ToString())
				}
				IsTrue = TT(match)
			}
			MSG.SwitchValue(reflect.Value(IsTrue).Bool(), "", "正则不匹配")
		}
		if IsTrue.ToString() == "false" {
			return nil, false, MSG.ToString()
		}
	}
	return Field.Value(), true, ""
}

// 参数校验
func (h *Http) Verify(configs []map[string]interface{}) {
	params := *(*h).Params
	result, isTrue, messages := make(map[string]interface{}), true, make([]string, 0)
	for _, config := range configs {
		// 字段名称
		field := config["field"].(string)
		// 校验规则
		rules := make([]string, 0)
		if _, ok := config["rules"]; ok && (reflect.TypeOf(config["rules"]).String() == "[]interface {}") {
			for _, v := range config["rules"].([]interface{}) {
				rules = append(rules, v.(string))
			}
		}
		// 参数是否必填判断
		if _, ok := params[field]; (!ok) || (params[field] == nil) || (params[field] == "") {
			params[field] = ""
			if _, ok = config["default"]; ok {
				params[field] = config["default"]
			}
		}
		// 参数正则判断
		value, isOk, msg := h.checkField(params[field], rules)
		if !isOk {
			isTrue = false
			if _, ok := config["label"]; ok {
				messages = append(messages, config["label"].(string)+msg)
			}
		} else {
			result[field] = value
		}
	}
	(*h.Params)["__"] = result
	if !isTrue {
		h.Output(402, "请求失败", strings.Join(messages, "；"))
	}
}

// 请求返回
func (h *Http) Output(code int, args ...interface{}) {
	if h.Rep == nil {
		panic("EOF")
	}
	result := make(map[String]interface{})
	result["code"] = code
	if code == 200 {
		if len(args) > 0 {
			result["data"] = args[0]
			if len(args) > 1 {
				t := TT(args[1])
				if t.IsString() && (!t.IsEmpty()) {
					result["message"] = args[1]
				}
			}
		}
	} else if code == 503 {
		if len(args) > 0 {
			result["message"] = args[0].(string) + "\n" + strings.Join(args[1].([]string), "\n")
		}
	} else {
		if len(args) > 0 {
			t := TT(args[0])
			if t.IsString() && (!t.IsEmpty()) {
				result["message"] = args[0]
			}
		}
	}
	result["duration"] = fmt.Sprintf("%s", time.Since(h.STime))
	if t, _ := String("").C("app.debug"); t.IsValid() && t.IsBool() && (t.ToString() == "true") {
		result["params"] = h.Params
	}
	bs, _ := json.Marshal(result)
	h.Rep.Write(bs)
	h.Rep = nil
}
