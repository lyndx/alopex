package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Http struct {
	STime  int64
	Rep    http.ResponseWriter
	Req    *http.Request
	Params *map[string]interface{}
	Errors []string
}

type HttpError struct {
	Code    string                 `json:"code"`
	Msg     string                 `json:"msg"`
	MsgMore string                 `json:"msg_more"`
	Data    map[string]interface{} `json:"data"`
}

// 获取参数清单
func (h *Http) GetParams() *Http {
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
	h.Params = &params
	return h
}

// 字段校验
func (h *Http) checkField(field interface{}, rules []string) (interface{}, bool, string) {
	Field, IsTrue, MSG := T(reflect.ValueOf(field)), T(reflect.ValueOf(true)), T(reflect.ValueOf(""))
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
			Field = T(reflect.ValueOf(Field.ToString()))
			IsTrue.SwitchValue(IsValid, Field.IsInt(), true)
			MSG.SwitchValue(reflect.Value(IsTrue).Bool(), "", "必须为整数")
		case "float":
			Field = T(reflect.ValueOf(Field.ToString()))
			IsTrue.SwitchValue(IsValid, Field.IsFloat(), true)
			MSG.SwitchValue(reflect.Value(IsTrue).Bool(), "", "必须为浮点数")
		case "bool":
			Field = T(reflect.ValueOf(Field.ToString()))
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
				IsTrue = T(reflect.ValueOf(match))
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
		h.Output(HttpError{"402", "请求失败", strings.Join(messages, "；"), nil})
	}
}

// 请求返回
func (h *Http) Output(err HttpError) {
	time.Sleep(1)
	if h.Rep == nil {
		panic(map[string]interface{}{"sdfsd": 23423})
	}
	bs, _ := json.Marshal(err)
	tmp := make(map[string]interface{})
	json.Unmarshal(bs, &tmp)
	tmp["duration"] = strconv.Itoa(int(float64(time.Now().UnixNano())-float64(h.STime))) + "毫秒"
	bs, _ = json.Marshal(tmp)
	h.Rep.Write(bs)
	h.Rep = nil
}
