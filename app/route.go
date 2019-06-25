package app

import (
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/spf13/viper"
	//"github.com/valyala/fasthttp"
)

var (
	route     *T
	routeOnce sync.Once
)

// 加载路由
func (s String) R(key string) (T, error) {
	routeOnce.Do(func() {
		result := make(map[string]map[string][]interface{})
		v := viper.New()
		v.SetConfigType("yml")
		for _, dir := range String("route").Scan("", true) {
			result[dir] = make(map[string][]interface{})
			v.AddConfigPath("route/" + dir)
			for _, file := range String("route/"+dir).Scan(".yml", false) {
				rname := strings.Replace(file, ".yml", "", -1)
				v.SetConfigName(rname)
				if err := v.ReadInConfig(); err != nil {
					DIE("路由文件加载失败....")
				}
				items := make(map[string]interface{})
				for _, k := range v.AllKeys() {
					items[k] = v.Get(k)
				}
				if items["list"] != nil {
					result[dir][rname] = items["list"].([]interface{})
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
		for _, file := range String("route/"+module).Scan(".yml", false) {
			routes, _ := String(module).R(strings.Replace(file, ".yml", "", -1))
			if routes.IsValid() {
				for _, item := range TValue(routes).([]interface{}) {
					route := item.(map[string]interface{})
					params := route["params"].([]interface{})
					needAuth := route["need_auth"].(bool)
					withPlatform := route["with_platform"].(bool)
					if withPlatform{
						//route = "/"
					}
					handler := route["handler"].(string)
					//

					//m := func(ctx *fasthttp.RequestCtx) {
					//	switch string(ctx.Path()) {
					//	case "/foo":
					//		fooHandlerFunc(ctx)
					//	case "/bar":
					//		barHandlerFunc(ctx)
					//	case "/baz":
					//		bazHandler.HandlerFunc(ctx)
					//	default:
					//		ctx.Error("not found", fasthttp.StatusNotFound)
					//	}
					//}
					//fasthttp.ListenAndServe(":80", m)

					http.HandleFunc("/"+module+"/"+route["route"].(string), func(rep http.ResponseWriter, req *http.Request) {
						defer PHandler()
						// 初始访问
						h := NT(rep, req)
						h.Module = module
						// 参数校验
						h.Verify(params, module, needAuth, withPlatform)
						// 业务实现
						h.RHH(String(req.URL.Path).Split("/")[1], handler)
					})
				}
			}
		}
	}
}
