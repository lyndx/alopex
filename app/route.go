package app

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
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
			for _, file := range String("route/" + dir).Scan(".yml", false) {
				rname := strings.Replace(file, ".yml", "", -1)
				v.SetConfigName(rname)
				if err := v.ReadInConfig(); err != nil {
					fmt.Println("[ERROR] 路由文件加载失败....")
					os.Exit(1)
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
		tmp := TT(TT(result).MapParse())
		route = &tmp
	})
	path := s.ToString()
	if path == "" {
		return TT(nil), errors.New("参数错误")
	}
	if key != "" {
		key = path + "." + key
	} else {
		key = path
	}
	result := (*route).GetValue(key, false)
	return result, nil
}
