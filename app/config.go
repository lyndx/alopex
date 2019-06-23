package app

import (
	"errors"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	config      *T
	IsDeveloper bool
	configOnce  sync.Once
)

// 加载配置
func (s String) C(args ...string) (T, error) {
	configOnce.Do(func() {
		result := make(map[string]map[string]interface{})
		v := viper.New()
		v.SetConfigType("yml")
		v.AddConfigPath("config")
		for _, file := range String("config").Scan(".yml", false) {
			cname := strings.Replace(file, ".yml", "", -1)
			v.SetConfigName(cname)
			if err := v.ReadInConfig(); err != nil {
				DIE("配置文件加载失败....")
			}
			items := make(map[string]interface{})
			for _, k := range v.AllKeys() {
				items[k] = v.Get(k)
			}
			result[cname] = items
		}
		tmp := TT(TT(result).ToMS(), true)
		config = &tmp
		//
		IsDeveloper = false
		x := TT(config.GValue("app.is_developer", false))
		if x.IsValid() && x.IsBool() && TValue(x, true).(bool) {
			IsDeveloper = true
		}
	})
	path := s.ToString()
	if path == "" {
		return TT(nil), errors.New("参数错误")
	}
	key := path
	if len(args) == 1 {
		key += "." + args[0]
	}
	result := TT((*config).GValue(key, false))
	if !result.IsValid() {
		return TT(nil), errors.New("配置项数据获取失败")
	}
	return result, nil
}
