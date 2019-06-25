package service

import (
	"strconv"
	"time"

	"alopex/app"

	"github.com/fwhezfwhez/jwt"
	"github.com/pkg/errors"
)

type AuthService struct{}

func init() {
	app.Services["auth"] = AuthService{}
}

// 生成认证Token & 随机字符串
func (a AuthService) GetToken(module string, userId string) (string, string, error) {
	tmp, err := app.String("auth").C("salt")
	if (err != nil) || tmp.IsEmpty() {
		return "", "", errors.New("认证配置信息获取失败")
	}
	// 加密盐值
	tokenSecret := tmp.ToString()
	tmp, err = app.String("auth").C(module + "_token_expires")
	if (err != nil) || (!tmp.IsInt()) {
		return "", "", errors.New("认证配置信息获取失败")
	}
	te := app.TValue(tmp).(int)
	if te < 1 {
		return "", "", errors.New("认证Token有效期配置错误")
	}
	// Token有效期
	tokenExpires := strconv.Itoa(int(time.Now().Add(time.Duration(te) * time.Minute).Unix()))
	// 随机字符串
	randomStr := app.String("").GUID()
	if randomStr == "" {
		return "", "", errors.New("随机字符串生成失败")
	}
	// 生成Token
	token := jwt.GetToken()
	token.AddHeader("typ", "JWT").AddHeader("alg", "HS256")
	token.AddPayLoad("user_id", userId).AddPayLoad("random_str", randomStr).AddPayLoad("exp", tokenExpires)
	tokenStr, _, err := token.JwtGenerator(tokenSecret)
	if (err != nil) || (tokenStr == "") {
		return "", "", errors.New("认证Token生成失败")
	}
	return tokenStr, randomStr, nil
}
