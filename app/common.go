package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

type (
	Float  float64
	String string
	Int    int64
	Bool   bool
	T      reflect.Value
)

// 异常消息
func PHandler() {
	if err := recover(); err != nil {
		_, filename, _, _ := runtime.Caller(0)
		RPath, PID, TIME := path.Dir(path.Dir(filename)), strconv.Itoa(os.Getpid()), time.Now().Format("2006/01/02 15:04:05")
		EMsg := fmt.Sprintf("%v", err)
		if EMsg == "EOF" {
			return
		}
		title, stack := "["+TIME+"][PID:"+PID+"]\n  “"+EMsg+"”", make([]string, 0)
		tmp := strings.Split(string(debug.Stack()), "\n")
		if len(tmp) > 0 {
			for k, v := range tmp {
				if (k%2 == 0) || (k == 1) || (v == "") || (v == "alopex/app.PHandler()") {
					continue
				}
				rec, _ := regexp.Compile(`\(0x[a-z0-9]+`)
				v = rec.ReplaceAllString(v, "(?")
				rec, _ = regexp.Compile(` ?0x[a-z0-9]+\)`)
				v = rec.ReplaceAllString(v, "?)")
				rec, _ = regexp.Compile(`, 0x[a-z0-9]+,`)
				v = rec.ReplaceAllString(v, ",?,")
				rec, _ = regexp.Compile(`, 0x[a-z0-9]+,`)
				v = rec.ReplaceAllString(v, ",?,")
				path := strings.TrimLeft(strings.Split(tmp[k+1], " +0x")[0], "\t")
				if strings.HasPrefix(path, RPath) {
					stack = append(stack, v+"  "+path)
				}
			}
		}
		fmt.Println("\n" + title)
		if len(stack) > 0 {
			fmt.Println("----------------------------------------------")
			for _, v := range stack {
				fmt.Println("》> " + v)
			}
			fmt.Println()
		}
	}
}

// 实例化T对象
func TT(t interface{}) T {
	tmp := reflect.ValueOf(t)
	if !tmp.IsValid() {
		return T(reflect.ValueOf(nil))
	}
	tp := tmp.Type().String()
	if tp == "reflect.Value" {
		return TT(t.(reflect.Value).Interface())
	}
	if tp == "app.T" {
		return T(reflect.ValueOf(t.(T).Value()))
	}
	if tp == "interface {}" {
		tmp = tmp.Elem()
		return TT(tmp)
	}
	return T(reflect.ValueOf(tmp))
}

// 首字母大写
func (s String) UFrist() string {
	str := string(s)
	if str == "" {
		return str
	}
	str = strings.ToLower(str)
	str = strings.ToUpper(string([]rune(str)[0])) + str[1:]
	return str
}

// 字符串转整数
func (s String) ToInt() Int {
	str := string(s)
	if match, _ := regexp.MatchString(`^[0-9]+$`, str); match {
		i, _ := strconv.Atoi(str)
		return Int(i)
	}
	return Int(0)
}

// 字符串转浮点数
func (s String) ToFloat() Float {
	str := string(s)
	if match, _ := regexp.MatchString(`^[0-9]+(\.[0-9]+)$`, str); match {
		f, _ := strconv.ParseFloat(str, 64)
		return Float(f)
	}
	return Float(0)
}

// 字符串转布尔值
func (s String) ToBool() Bool {
	str := string(s)
	if match, _ := regexp.MatchString(`^(true|false)$`, strings.ToLower(str)); match {
		b, _ := strconv.ParseBool(str)
		return Bool(b)
	}
	return Bool(false)
}

// 目录/文件扫描
func (s String) Scan(suffix string, isDir bool) []string {
	files := make([]string, 0)
	err := filepath.Walk(s.ToString(), func(path string, f os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		filename := f.Name()
		// 文件清单 / 目录清单
		if (filepath.Base(s.ToString()) != filename) && (!strings.HasPrefix(filename, ".")) && ((isDir && f.IsDir()) || ((!isDir) && (!f.IsDir()))) {
			if (suffix != "") && (!strings.HasSuffix(filename, suffix)) {
				return errors.New("文件后缀错误")
			}
			files = append(files, filename)
		}
		return nil
	})
	if err != nil {
		return make([]string, 0)
	}
	return files
}

// 整数转字符串
func (i Int) ToString() string {
	return strconv.Itoa(int(i))
}

// 浮点数转字符串
func (f Float) ToString(prec int) string {
	return strconv.FormatFloat(float64(f), 'f', prec, 64)
}

// 布尔值转字符串
func (b Bool) ToString() string {
	return strconv.FormatBool(bool(b))
}

// 字符串转换
func (s String) ToString() string {
	return string(s)
}

// 校验是否为有效参数
func (t T) IsValid() bool {
	return reflect.Value(t).IsValid()
}

// 任何类型转字符串
func (t T) ToString() string {
	if !t.IsValid() {
		return ""
	}
	vv := reflect.ValueOf(TT(t).MapParse()).Interface()
	switch vv.(type) {
	case string:
		return vv.(string)
	case int:
		return Int(vv.(int)).ToString()
	case byte:
		return Int(int64(vv.(byte))).ToString()
	case int8:
		return Int(int64(vv.(int8))).ToString()
	case int32:
		return Int(int64(vv.(int32))).ToString()
	case int64:
		return Int(vv.(int64)).ToString()
	case float32:
		return Float(float64(vv.(float32))).ToString(2)
	case float64:
		return Float(vv.(float64)).ToString(2)
	case bool:
		return Bool(vv.(bool)).ToString()
	case []interface{}, map[string]interface{}, map[interface{}]interface{}:
		bs, _ := json.Marshal(vv)
		return string(bs)
	}
	return ""
}

// 对象/数组再解析，便于转成JSON字符串
func (t T) MapParse() interface{} {
	if !t.IsValid() {
		return nil
	}
	v := reflect.Value(t)
	tt := v.Type().String()
	if tt == "interface {}" {
		v = v.Elem()
		tt = v.Type().String()
	}
	if strings.HasPrefix(tt, "[]") {
		result := make([]interface{}, 0)
		for i := 0; i < v.Len(); i++ {
			vv := T(v.Index(i)).MapParse()
			result = append(result, vv)
		}
		return result
	} else if strings.HasPrefix(tt, "map[") {
		result := make(map[string]interface{})
		for _, k := range v.MapKeys() {
			key := T(k).ToString()
			if key != "" {
				kv := v.MapIndex(k)
				result[key] = T(kv).MapParse()
			}
		}
		return result
	}
	return v.Interface()
}

// 根据KEY获取数据，针对数组和对象
func (t T) GetValue(key string, isStrict bool) T {
	v := reflect.ValueOf(t.MapParse())
	match, _ := regexp.MatchString(`^[0-9a-zA-Z_\/]+(\.[0-9a-zA-Z_\/]+)*$`, key)
	if (!v.IsValid()) || (!match) {
		return TT(nil)
	}
	tt := v.Type().String()
	if tt == "interface {}" {
		v = v.Elem()
		return TT(v).GetValue(key, isStrict)
	}
	if tt == "app.T" {
		return v.Interface().(T).GetValue(key, isStrict)
	}
	if tt == "reflect.Value" {
		v = v.Interface().(reflect.Value)
		tt = v.Type().String()
	}
	isMap, isArr := strings.HasPrefix(tt, "map["), strings.HasPrefix(tt, "[]")
	if !(isMap || isArr) {
		return TT(nil)
	}
	keys := strings.Split(key, ".")
	nk, kv := -1, reflect.ValueOf(keys[0])
	if match, _ := regexp.MatchString(`^(0|[1-9][0-9]*)$`, keys[0]); match {
		nk, _ = strconv.Atoi(keys[0])
	}
	if isMap {
		vv := v.MapIndex(kv)
		if (!vv.IsValid()) && (!isStrict) {
			tmp := make(map[string]interface{})
			for _, mk := range v.MapKeys() {
				if mk.Type().String() == "string" {
					mks := mk.String()
					if strings.HasPrefix(mks, keys[0]+".") {
						mks = mks[len(keys[0]+"."):]
						tmp[mks] = v.MapIndex(mk).Interface()
					}
				}
			}
			if len(tmp) > 0 {
				vv = reflect.ValueOf(tmp)
			}
		}
		if (!vv.IsValid()) && (nk > -1) {
			vv = v.MapIndex(reflect.ValueOf(nk))
		}
		if vv.IsValid() {
			if len(keys) > 1 {
				key = strings.Join(keys[1:], ".")
				return T(vv).GetValue(key, isStrict)
			}
			return T(vv)
		}
	} else if nk > -1 {
		vv := v.Index(nk)
		if vv.IsValid() {
			if len(keys) > 1 {
				key = strings.Join(keys[1:], ".")
				return T(vv).GetValue(key, isStrict)
			}
			return T(vv)
		}
	}
	return TT(nil)
}

// 返回真实数值
func (t T) Value() interface{} {
	if !t.IsValid() {
		return nil
	}
	v := reflect.Value(t).Interface()
	vv := reflect.ValueOf(v)
	tp := vv.Type().String()
	if tp == "reflect.Value" {
		vv = reflect.ValueOf(v.(reflect.Value).Interface())
		return T(vv).Value()
	}
	if tp == "interface {}" {
		vv = reflect.ValueOf(vv.Elem())
		return T(vv).Value()
	}
	return v
}

// 根据条件取对应值
func (t *T) SwitchValue(conditions bool, trueValue interface{}, falseValue interface{}) interface{} {
	if conditions {
		*t = TT(trueValue)
	} else {
		*t = TT(falseValue)
	}
	return *t
}

// 判断是否为文件
func (t T) IsFile(checkIsFileArray bool) bool {
	if !t.IsValid() {
		return false
	}
	t = TT(t)
	if checkIsFileArray {
		return reflect.Value(t).Type().String() == "[]*multipart.FileHeader"
	}
	return reflect.Value(t).Type().String() == "*multipart.FileHeader"
}

// 判断是否为字符串
func (t T) IsString() bool {
	if !t.IsValid() {
		return false
	}
	t = TT(t)
	return reflect.Value(t).Type().String() == "string"
}

// 判断是否为整数
func (t T) IsInt() bool {
	if !t.IsValid() {
		return false
	}
	t = TT(t)
	tp := reflect.Value(t).Type().String()
	if strings.HasPrefix(tp, "int") && (tp != "interface {}") {
		return true
	} else if t.IsString() {
		match, _ := regexp.MatchString(`^(0|[1-9][0-9]*)$`, t.ToString())
		return match
	}
	return false
}

// 判断是否为浮点数
func (t T) IsFloat() bool {
	if !t.IsValid() {
		return false
	}
	t = TT(t)
	tp := reflect.Value(t).Type().String()
	if strings.HasPrefix(tp, "float") {
		return true
	} else if t.IsString() {
		match, _ := regexp.MatchString(`^(0|[1-9][0-9]*)(\.[0-9]+)?$`, t.ToString())
		return match
	}
	return false
}

// 判断是否为布尔值
func (t T) IsBool() bool {
	if !t.IsValid() {
		return false
	}
	t = TT(t)
	tp := reflect.Value(t).Type().String()
	if strings.HasPrefix(tp, "bool") {
		return true
	} else if t.IsString() {
		match, _ := regexp.MatchString(`^(TRUE|FALSE|YES|NO)$`, strings.ToUpper(t.ToString()))
		return match
	}
	return false
}

// 判断是否为数组
func (t T) IsArray() bool {
	if !t.IsValid() {
		return false
	}
	t = TT(t)
	tp := reflect.Value(t).Type().String()
	if strings.HasPrefix(tp, "[]") {
		return true
	}
	return false
}

// 判断是否为空
func (t T) IsEmpty() bool {
	if !t.IsValid() {
		return true
	}
	t = TT(t)
	if t.IsString() && (t.ToString() == "") {
		return true
	}
	return false
}
