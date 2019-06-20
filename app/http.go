package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
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
	// 跨域处理
	if origin := req.Header.Get("Origin"); origin != "" {
		req.Header.Set("Access-Control-Allow-Origin", "*")
		req.Header.Set("Access-Control-Max-Age", "172800")
		req.Header.Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
		req.Header.Set("Access-Control-Allow-Headers", "Authorization,Accept-Language,Cache-Control,Content-Type")
		req.Header.Set("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
	}
	if req.Method == "OPTIONS" {
		rep.Header().Set("content-type", "text/plain")
		rep.Write([]byte("Options Request!"))
		panic("EOF")
	} else {
		rep.Header().Set("content-type", "application/json")
	}
	//
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
		params[k] = v
		if len(v) == 1 {
			params[k] = v[0]
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
			MSG.SwitchValue(TValue(IsTrue).(bool), "", "必须为文件")
		case "files":
			IsTrue.SwitchValue(IsValid, Field.IsFile(true), true)
			MSG.SwitchValue(TValue(IsTrue).(bool), "", "必须为文件数组")
		case "must":
			IsTrue.SwitchValue(IsValid, !Field.IsEmpty(), true)
			MSG.SwitchValue(TValue(IsTrue).(bool), "", "为必填字段")
		case "string":
			IsTrue.SwitchValue(IsValid, Field.IsString(), true)
			MSG.SwitchValue(TValue(IsTrue).(bool), "", "必须为字符串格式")
		case "int":
			Field = TT(Field.ToString())
			IsTrue.SwitchValue(IsValid, Field.IsInt(), true)
			MSG.SwitchValue(TValue(IsTrue).(bool), "", "必须为整数")
		case "float":
			Field = TT(Field.ToString())
			IsTrue.SwitchValue(IsValid, Field.IsFloat(), true)
			MSG.SwitchValue(TValue(IsTrue).(bool), "", "必须为浮点数")
		case "bool":
			Field = TT(Field.ToString())
			IsTrue.SwitchValue(IsValid, Field.IsBool(), true)
			MSG.SwitchValue(TValue(IsTrue).(bool), "", "必须为布尔值")
		case "array":
			IsTrue.SwitchValue(IsValid, Field.IsArray(), true)
			MSG.SwitchValue(TValue(IsTrue).(bool), "", "必须为数组")
		default:
			if IsValid {
				match := false
				if Field.IsString() {
					match, _ = regexp.MatchString(rule, Field.ToString())
				}
				IsTrue = TT(match)
			}
			MSG.SwitchValue(TValue(IsTrue).(bool), "", "正则不匹配")
		}
		if IsTrue.ToString() == "false" {
			return TValue(Field), false, MSG.ToString()
		}
	}
	return TValue(Field), true, ""
}

// 参数校验
func (h *Http) Verify(configs []interface{}, needAuth bool, withPlatform bool) {
	params := *(*h).Params
	result, isTrue, messages := make(map[string]interface{}), true, make([]string, 0)
	for _, item := range configs {
		config := item.(map[string]interface{})
		// 字段名称
		field := config["field"].(string)
		// 校验规则
		rules := make([]string, 0)
		if _, ok := config["rules"]; ok && TT(config["rules"]).IsArray() {
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
		}
		result[field] = value
	}
	(*h.Params)["__"] = result
	if !isTrue {
		h.Output(402, "请求失败", strings.Join(messages, "；"))
	}
}

// 请求返回
func (h *Http) Output(code int, args ...interface{}) {
	defer func() {
		panic("EOF")
	}()
	if h.Rep == nil {
		panic("EOF")
	}
	result := make(map[String]interface{})
	result["code"] = code
	if len(args) > 0 {
		if code == 200 {
			result["data"] = args[0]
			if len(args) > 1 {
				t := TT(args[1])
				if t.IsString() && (!t.IsEmpty()) {
					result["message"] = args[1]
				}
			}
		} else {
			t := TT(args[0])
			if t.IsString() && (!t.IsEmpty()) {
				result["message"] = args[0]
			}
			if len(args) > 1 {
				t := TT(args[1])
				if t.IsString() && (!t.IsEmpty()) {
					result["message_detail"] = args[1]
				}
			}
		}
	}
	result["duration"] = time.Since(h.STime).String()
	if t, _ := String("app").C("is_developer"); t.IsValid() && t.IsBool() && TValue(t).(bool) {
		PA, PB := make(map[string]interface{}), (*h.Params)["__"]
		for k, v := range *h.Params {
			if k != "__" {
				vv := TT(v, true)
				if vv.IsFile(false) {
					f := v.(*multipart.FileHeader)
					PA[k] = map[string]interface{}{"name": f.Filename, "size": Float(float64(f.Size)/float64(1024)).ToString(2) + "KB", "type": f.Header.Get("Content-Type")}
				} else if vv.IsFile(true) {
					fs := v.([]*multipart.FileHeader)
					fitems := make([]map[string]interface{}, 0)
					for _, f := range fs {
						fitems = append(fitems, map[string]interface{}{"name": f.Filename, "size": Float(float64(f.Size)/float64(1024)).ToString(2) + "KB", "type": f.Header.Get("Content-Type")})
					}
					PA[k] = fitems
				} else if k == "content-type" {
					PA[k] = String(v.(string)).Split(";")[0]
				} else {
					PA[k] = v
				}
			}
		}
		result["request"] = map[string]interface{}{
			"params": PA,
			"needed": PB,
		}
	}
	bs, _ := json.Marshal(result)
	h.Rep.WriteHeader(200)
	h.Rep.Header().Set("token", "")
	if token, ok := (*h.params())["token"]; ok && (token != "") {
		h.Rep.Header().Set("token", token.(string))
	}
	// 可以考虑GZIP传输
	h.Rep.Write(bs)
	h.Rep = nil
}

// 控制器
func (h *Http) RHH(handler string) {
	hh := String(handler).Split(".")
	if len(hh) != 2 {
		h.Output(402, "请求失败", "路由配置中Handler项格式错误")
	}
	contrl, action := hh[0], hh[1]

	fmt.Println(contrl, action)

}
