package route

import (
	. "encoding/json"
	"errors"
	. "fmt"
	. "os"
	. "path/filepath"
	"reflect"
	"strings"
	. "time"

	. "qp/config"
	. "qp/model"
	. "qp/model/biz"
	. "qp/tool"

	"github.com/iris-contrib/middleware/cors"
	. "github.com/kataras/iris"
)

// 路由对象
type Object struct {
	handler Handler
}

// 获取跨域执行方法
func getCorsHandler() Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // 允许通过的主机名称
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Origin", "Accept", "X-Requested-With", "PLATFORM", "Content-Type"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		Debug:            false,
	})
}

// 获取平台校验执行方法
func getPlatformHandler() Handler {
	return func(c Context) {
		// 平台标识
		platform := c.Params().Get("platform")
		// 请求时间开始计时
		c.Values().Set("requestCurrentTime", Now().UnixNano()/1e3)
		// 平台标识验证
		if PlatformDbConfigs[platform] != nil {
			CurrentPlatformPool = PlatformPools[platform]
			c.Next()
		} else {
			Output(&c, "没有这个平台号", "无法识别的平台号", ERROR, nil)
		}
	}
}

// 获取用户认证执行方法
func getAuthenticateHandler(needAuth bool) Handler {
	return func(c Context) {
		if needAuth {
			token, userId, ok := JwtToken(c)
			if ok && CheckToken(token, userId) {
				c.Next()
			} else {
				return
			}
		}
		c.Next()
	}
}

// 获取参数校验执行方法
func getParamsHandler(params []map[string]interface{}) Handler {
	return func(c Context) {
		result := make(map[string]interface{})
		fields := BuildParams(c)
		isTrue, messages := true, []string{}
		for _, cfg := range params {
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
	}
}

// 获取业务执行方法
func getOperateHandler(controllerName string, actionName string) Handler {
	return func(c Context) {
		module := "mobile"
		if IsBackendService {
			module = "backend"
		}
		if _, ok := ControllerMapping[module][controllerName]; ok {
			ctrl := *ControllerMapping[module][controllerName]
			_, err := Todo(ctrl, actionName, []interface{}{&c})
			if err != nil {
				Output(&c, err.Error(), "请求失败", ERROR, nil)
			}
		} else {
			Output(&c, "请求执行方法不存在", "请求失败", ERROR, nil)
		}
	}
}

// 路由映射
func mapping(app *Application, routeName string, method string, needAuth bool, params []map[string]interface{}, handler []string) error {
	// 请求方法, 控制器名，操作名
	method, controllerName, actionName := Ufirst(method), handler[0], Ufirst(handler[1])
	// 跨域访问中间件
	corsHandler := Object{getCorsHandler()}
	// 平台校验中间件
	platformHandler := Object{getPlatformHandler()}
	// 用户认证中间件
	authenticateHandler := Object{getAuthenticateHandler(needAuth)}
	// 参数校验中间件
	paramsHandler := Object{getParamsHandler(params)}
	// 业务执行方法
	operateHandler := Object{getOperateHandler(controllerName, actionName)}
	// 加载完成
	_, err := Todo(app, method, []interface{}{
		routeName,
		corsHandler.handler,
		platformHandler.handler,
		authenticateHandler.handler,
		paramsHandler.handler,
		operateHandler.handler,
	},
	)
	return err
}

// 路由加载
func handle(app *Application, routingName string, routingConfig interface{}) error {
	if (routingConfig == nil) || (reflect.TypeOf(routingConfig).String() != "map[interface {}]interface {}") {
		return errors.New("路由[" + routingName + "]的配置数据错误")
	}
	ncfg := routingConfig.(map[interface{}]interface{})
	// 校验请求方式
	method := strings.ToUpper(CheckStringField(ncfg, "method", ""))
	if (method == "") || (strings.Index("POST,GET,OPTIONS,DELETE,PUT,HEAD", method) == -1) {
		return errors.New("路由[" + routingName + "]的请求方式错误")
	}
	// 是否带上平台标识
	withPlatform := CheckBoolField(ncfg, "withPlatform", false)
	if withPlatform {
		routingName = "/{platform:string}/" + routingName
	} else {
		routingName = "/" + routingName
	}
	// 是否需要认证
	needAuth := CheckBoolField(ncfg, "needAuth", false)
	// 参数校验
	params := make([]map[string]interface{}, 0)
	for field, cfg := range MapII2MapSI(CheckMapField(ncfg, "params", nil)) {
		cc := cfg.(map[string]interface{})
		cc["field"] = field
		params = append(params, cc)
	}
	// 操作校验
	handler := strings.Split(strings.ToLower(CheckStringField(ncfg, "handler", "")), ".")
	if len(handler) != 2 {
		return errors.New("路由[" + routingName + "]的执行方法未定义")
	}
	// 打印路由
	bs, _ := Marshal(params)
	Echo("  [" + routingName + "] [" + method + "] -> [" + handler[0] + "." + handler[1] + "] -> " + string(bs))
	// 路由映射
	return mapping(app, routingName, method, needAuth, params, handler)
}

// 遍历配置
func scanRoutingPath(app *Application, routingPath string) {
	err := Walk(routingPath, func(path string, f FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 文件名称
		filename := f.Name()
		// 文件校验
		if (f == nil) || f.IsDir() {
			return nil
		}
		if !strings.HasSuffix(filename, ".yml") {
			return errors.New("路由文件后缀错误")
		}
		// 解析配置
		routings, err := Yaml(routingPath + filename)
		if err != nil {
			return err
		}
		filename = strings.TrimRight(filename, ".yml")
		Echo("模块[" + strings.ToUpper(filename) + "] ->")
		for rName, rConfig := range routings {
			err := handle(app, rName, rConfig)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		Echo(true, "路由配置文件解析错误，"+err.Error()+"....")
	}
}

// 路由配置加载
func Loading(app *Application) {
	// 服务平台
	platform := "移动端"
	// 路由目录
	routingPath := "./route/mobile/"
	if IsBackendService {
		platform = "管理端"
		routingPath = "./route/backend/"
	}
	Echo("\n>>>> =======================================================================\n" + platform + "路由加载....")
	// 执行加载路由配置
	scanRoutingPath(app, routingPath)
	//
	Println()
}
