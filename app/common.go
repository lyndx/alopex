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
	Float float64
	String string
	Int int64
	Bool bool
	T reflect.Value
)

// 运行时相关 //////////////////////////////////////////////////////////////////
// 格式化命令行显示
func Dump(args ...string) {
	length := len(args)
	if length == 0 {
		fmt.Println()
	}
	if length == 1 {
		fmt.Println(args[0])
	}
	if length == 2 {
		if runtime.GOOS != "windows" {
			code := 30
			switch args[0] {
			case "red":
				code = 31
			case "green":
				code = 32
			case "yellow":
				code = 33
			case "blue":
				code = 34
			case "white":
				code = 37
			}
			args[1] = fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", code, args[1])
		}
		fmt.Println(args[1])
	}
}

// 打印空行
func ELine(num int) {
	if num < 1 {
		num = 1
	}
	for i := 0; i < num; i++ {
		Dump()
	}
}

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
				if (k%2 == 0) || (k == 1) || (v == "") || strings.HasSuffix(v, "/app.PHandler()") {
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
		Dump()
		Dump("red", title)
		if len(stack) > 0 {
			Dump("----------------------------------------------")
			for _, v := range stack {
				Dump("green", "》> "+v)
			}
			Dump()
		}
	}
}

// 终结运行
func DIE(message string, args ...bool) {
	if !((len(args) > 0) && args[0]) {
		Dump("red", "[ERROR]"+message)
	}
	Dump(message)
	os.Exit(1)
}

// 字符串相关 //////////////////////////////////////////////////////////////////
// 首字母大写
func (s String) UFrist() string {
	str := string(s)
	if str != "" {
		str = strings.ToUpper(string([]rune(str)[0])) + str[1:]
	}
	return str
}

// 字符串转整数
func (s String) ToInt() Int {
	str := string(s)
	if match, _ := regexp.MatchString(`^[0-9]+$`, str); match {
		i, _ := strconv.Atoi(str)
		return Int(i)
	}
	return Int(-1)
}

// 字符串转浮点数
func (s String) ToFloat() Float {
	str := string(s)
	if match, _ := regexp.MatchString(`^[0-9]+(\.[0-9]+)?$`, str); match {
		f, _ := strconv.ParseFloat(str, 64)
		return Float(f)
	}
	return Float(0)
}

// 字符串转布尔值
func (s String) ToBool() Bool {
	str := strings.ToLower(string(s))
	b, _ := strconv.ParseBool(str)
	return Bool(b)
}

// 分割字符串
func (s String) Split(str string) []string {
	list := make([]string, 0)
	for _, v := range strings.Split(s.ToString(), str) {
		list = append(list, v)
	}
	return list
}

// 判断文件是否存在
func (s String) IsExist(file string) bool {
	path, err := filepath.Abs(s.ToString())
	if err != nil {
		return false
	}
	_, err = os.Stat(path + string(os.PathSeparator) + file)
	return err == nil
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

// 写文件
func (s String) Write(content string) (bool, error) {
	f, err := os.OpenFile(s.ToString(), os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return false, err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return false, err
	}
	return true, nil
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

// 泛数据相关 //////////////////////////////////////////////////////////////////
// 获取数据的类型，为nil时返回空字符串
func TType(t interface{}) string {
	if RV(t).IsValid() {
		return RT(t).String()
	}
	return ""
}

// 获取数据的最终数值，基础数据类型的数值
func TValue(obj interface{}, args ...bool) interface{} {
	noPointer := false
	if len(args) > 0 {
		noPointer = args[0]
	}
	switch TType(obj) {
	case "app.T":
		v := reflect.Value(obj.(T))
		if v.IsValid() {
			obj = v.Interface()
			return TValue(obj, noPointer)
		}
		return nil
	case "reflect.Value":
		ov := obj.(reflect.Value)
		if ov.IsValid() {
			obj = ov.Interface()
		} else {
			obj = nil
		}
		return TValue(obj, noPointer)
	case "interface {}":
		ov := RV(obj)
		if ov.IsValid() {
			obj = ov.Elem()
			return TValue(obj, noPointer)
		}
		return nil
	default:
		ov := RV(obj)
		if ov.IsValid() {
			if noPointer {
				ot := TType(obj)
				if strings.HasPrefix(ot, "*") {
					obj = ov.Elem()
					return TValue(obj, noPointer)
				}
			}
			return obj
		}
		return nil
	}
}

// 获取泛对象
func RV(t interface{}) reflect.Value {
	return reflect.ValueOf(t)
}

// 获取泛类型
func RT(t interface{}) reflect.Type {
	return reflect.TypeOf(t)
}

// 实例化T对象
func TT(t interface{}, args ...bool) T {
	t = TValue(t, args...)
	tp := TType(t)
	if tp == "" {
		return T(RV(nil))
	}
	return T(RV(t))
}

// 数组去掉空值,字符串数组去掉空字符串，其他数组去掉nil值
func (t *T) Filter(args ...func(value reflect.Value) bool) []interface{} {
	obj := TValue(t)
	if !TT(obj).IsArray() {
		return nil
	}
	list := make([]interface{}, 0)
	ov := RV(obj)
	length := ov.Len()
	for i := 0; i < length; i++ {
		v := ov.Index(i)
		ot := v.Type().String()
		isOk := (v.IsValid() && (ot != "string")) || ((ot == "string") && (v.Interface() != ""))
		if !isOk {
			continue
		}
		if (len(args) > 0) && (!args[0](v)) {
			continue
		}
		list = append(list, v.Interface())
	}
	vl := TT(list)
	t = &vl
	return list
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
	vv := TValue(t, true)
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
func (t T) ToMS() interface{} {
	obj := TValue(t, true)
	if obj == nil {
		return nil
	}
	v, tp := RV(obj), TType(obj)
	if strings.HasPrefix(tp, "[]") {
		result := make([]interface{}, 0)
		for i := 0; i < v.Len(); i++ {
			vv := TT(v.Index(i), true).ToMS()
			result = append(result, vv)
		}
		return result
	} else if strings.HasPrefix(tp, "map[") {
		result := make(map[string]interface{})
		for _, k := range v.MapKeys() {
			kk := TValue(k, true)
			if kk != nil {
				key := TT(kk).ToString()
				result[key] = TT(v.MapIndex(k), true).ToMS()
			}
		}
		return result
	}
	return v.Interface()
}

// 根据KEY获取数据，针对数组和对象
func (t T) GValue(key string, isStrict bool) T {
	match, _ := regexp.MatchString(`^[0-9a-zA-Z_\/]+(\.[0-9a-zA-Z_\/]+)*$`, key)
	if !match {
		return TT(nil)
	}
	obj := t.ToMS()
	if obj == nil {
		return TT(nil)
	}
	tp := TType(obj)
	IsMap, IsArr := strings.HasPrefix(tp, "map["), strings.HasPrefix(tp, "[]")
	if !(IsMap || IsArr) {
		return TT(nil)
	}
	keys := String(key).Split(".")
	nk, ov := int(String(keys[0]).ToInt()), RV(obj)
	if IsMap {
		kv := RV(keys[0])
		vv := RV(TValue(ov.MapIndex(kv), true))
		if (!vv.IsValid()) && (!isStrict) {
			tmp := make(map[string]interface{})
			for _, mk := range ov.MapKeys() {
				mktmp := TValue(mk, true)
				if mktmp == nil {
					continue
				}
				mktype := TType(mktmp)
				if mktype != "string" {
					continue
				}
				mkey := mk.String()
				if strings.HasPrefix(mkey, keys[0]+".") {
					mkey = mkey[len(keys[0]+"."):]
					tmp[mkey] = TValue(ov.MapIndex(mk).Interface())
				}
			}
			if len(tmp) > 0 {
				vv = RV(tmp)
			}
		}
		if (!vv.IsValid()) && (nk > -1) {
			vv = ov.MapIndex(RV(nk))
		}
		if vv.IsValid() {
			if len(keys) > 1 {
				key = strings.Join(keys[1:], ".")
				return TT(vv, true).GValue(key, isStrict)
			}
			return TT(vv, true)
		}
	} else if nk > -1 {
		vv := ov.Index(nk)
		if vv.IsValid() {
			if len(keys) > 1 {
				key = strings.Join(keys[1:], ".")
				return TT(vv, true).GValue(key, isStrict)
			}
			return TT(vv, true)
		}
	}
	return TT(nil)
}

// 根据条件取对应值
func (t *T) SwitchValue(conditions bool, trueValue interface{}, falseValue interface{}) interface{} {
	if conditions {
		*t = TT(trueValue, true)
	} else {
		*t = TT(falseValue, true)
	}
	return *t
}

// 判断是否为文件
func (t T) IsFile(checkIsFileArray bool) bool {
	obj := TValue(t, true)
	if obj == nil {
		return false
	}
	tp := TType(obj)
	if checkIsFileArray {
		return tp == "[]*multipart.FileHeader"
	}
	return tp == "multipart.FileHeader"
}

// 判断是否为字符串
func (t T) IsString() bool {
	obj := TValue(t, true)
	if obj == nil {
		return false
	}
	return TType(obj) == "string"
}

// 判断是否为整数
func (t T) IsInt() bool {
	obj := TValue(t, true)
	if obj == nil {
		return false
	}
	tp := TType(obj)
	t = TT(obj)
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
	obj := TValue(t, true)
	if obj == nil {
		return false
	}
	tp := TType(obj)
	if strings.HasPrefix(tp, "float") {
		return true
	} else if t.IsString() {
		match, _ := regexp.MatchString(`^(0|[1-9][0-9]*)(\.[0-9]+)?$`, TT(obj).ToString())
		return match
	}
	return false
}

// 判断是否为布尔值
func (t T) IsBool() bool {
	obj := TValue(t, true)
	if obj == nil {
		return false
	}
	tp := TType(obj)
	if strings.HasPrefix(tp, "bool") {
		return true
	} else if t.IsString() {
		match, _ := regexp.MatchString(`^(TRUE|FALSE)$`, strings.ToUpper(TT(obj).ToString()))
		return match
	}
	return false
}

// 判断是否为数组
func (t T) IsArray() bool {
	obj := TValue(t, true)
	if obj == nil {
		return false
	}
	tp := TType(obj)
	if strings.HasPrefix(tp, "[]") {
		return true
	}
	return false
}

// 判断是否为空
func (t T) IsEmpty() bool {
	obj := TValue(t, true)
	if obj == nil {
		return false
	}
	t = TT(obj)
	if t.IsString() && (t.ToString() == "") {
		return true
	}
	return false
}
