package app

import (
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var (
	route     *T
	routeOnce sync.Once
	Mux       = mux.NewRouter().StrictSlash(true)
)

// 加载路由
func (s String) R(key string) (T, error) {
	routeOnce.Do(func() {
		result := make(map[string]map[string]map[string]interface{})
		v := viper.New()
		v.SetConfigType("yml")
		for _, dir := range String("route").Scan("", true) {
			result[dir] = make(map[string]map[string]interface{})
			v.AddConfigPath("route/" + dir)
			for _, file := range String("route/" + dir).Scan(".yml", false) {
				rname := strings.Replace(file, ".yml", "", -1)
				v.SetConfigName(rname)
				if err := v.ReadInConfig(); err != nil {
					DIE("路由文件加载失败....")
				}
				item := make(map[string]interface{})
				for _, k := range v.AllKeys() {
					item[k] = v.Get(k)
				}
				if (item["name"] != nil) && (item["list"] != nil) {
					result[dir][rname] = item
				}
			}
		}
		tmp := TT(TT(result).ToMS(), true)
		route = &tmp
	})
	path := s.ToString()
	if path == "" {
		return TT(nil), errors.New("参数错误")
	}
	keyTmp := TT(key)
	key = (&keyTmp).SwitchValue(key != "", path+"."+key, path).(T).ToString()
	result := (*route).GValue(key, false)
	return result, nil
}

// 路由设置
func (s String) RH() {
	for _, module := range String("route").Scan("", true) {
		if s.ToString() != module {
			continue
		}
		for _, file := range String("route/" + module).Scan(".yml", false) {
			routes, _ := String(module).R(strings.Replace(file, ".yml", "", -1))
			if routes.IsValid() {
				for _, item := range TValue(routes).(map[string]interface{})["list"].([]interface{}) {
					route := item.(map[string]interface{})
					params := route["params"].([]interface{})
					needAuth := route["need_auth"].(bool)
					withPlatform := route["with_platform"].(bool)
					method := strings.ToUpper(route["method"].(string))
					routeStr := ""
					if withPlatform {
						routeStr = "/{platform}/" + route["route"].(string)
					} else {
						routeStr = "/" + route["route"].(string)
					}
					if module == "backend" {
						routeStr = "/backend" + routeStr
					}
					handler := route["handler"].(string)
					Mux.HandleFunc(routeStr, func(rep http.ResponseWriter, req *http.Request) {
						rep.Header().Set("Access-Control-Allow-Origin", "*")
						rep.Header().Set("Access-Control-Max-Age", "10000000")
						rep.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
						rep.Header().Set("Access-Control-Allow-Headers", "Content-Type,Debug,Token,Random_str")
						if req.Method == http.MethodOptions {
							return
						}
						defer PHandler()
						// 初始访问
						h := NT(rep, req)
						// 平台校验
						if withPlatform {
							platform := mux.Vars(req)["platform"]
							if platform == "" {
								h.Output(402, "请求失败", "平台标识获取失败")
							}
							_, err := String("database").C("platform_" + platform)
							if err != nil {
								h.Output(402, "请求失败", "平台标识"+err.Error())
							}
						}
						//
						h.Module = module
						// 上传参数
						if strings.HasSuffix(routeStr, "common/upload/{path}") {
							params = append(params, map[string]interface{}{"field": "file", "rules": []interface{}{"must", "file"}, "label": "上传文件"})
						}
						// 参数校验
						h.Verify(params, module, needAuth)
						// 业务实现
						h.RHH(String(req.URL.Path).Split("/")[1], handler)
					}).Methods(method, "OPTIONS")
				}
			}
		}
	}
}
