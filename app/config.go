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
	config     *T
	configOnce sync.Once
)

// 加载配置
func (s String) C(key string) (T, error) {
	configOnce.Do(func() {
		result := make(map[string]map[string]interface{})
		v := viper.New()
		v.SetConfigType("yml")
		v.AddConfigPath("config")
		for _, file := range String("config").Scan(".yml", false) {
			cname := strings.Replace(file, ".yml", "", -1)
			v.SetConfigName(cname)
			if err := v.ReadInConfig(); err != nil {
				fmt.Println("[ERROR] 配置文件加载失败....")
				os.Exit(1)
			}
			items := make(map[string]interface{})
			for _, k := range v.AllKeys() {
				items[k] = v.Get(k)
			}
			result[cname] = items
		}
		tmp := TT(TT(result).MapParse())
		config = &tmp
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
	result := TT((*config).GetValue(key, false))
	if !result.IsValid() {
		return TT(nil), errors.New("配置项数据获取失败")
	}
	return result, nil
}
