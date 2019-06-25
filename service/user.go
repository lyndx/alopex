package service

import (
	"math"
	"strconv"

	"alopex/app"
)

type UserService struct{}

func init() {
	app.Services["user"] = UserService{}
}

// 通过用户编号获取用户信息
func (u UserService) GetUserById(userId string) (interface{}, error) {
	user, err := app.MD("main").Select("users", true, "*", "id="+userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 通过用户名获取用户信息
func (u UserService) GetUserByUsername(username string) (interface{}, error) {
	user, err := app.MD("main").Select("users", true, "*", "username='"+username+"'")
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 获取用户清单
func (u UserService) GetUserList(conditions map[string]string, page int, size int) (map[string]interface{}, error) {
	source, fields, where, groupby, orderby, limit := conditions["source"], conditions["fields"], conditions["where"], conditions["groupby"], conditions["orderby"], ""
	if page > 0 {
		limit = strconv.Itoa(page)
		if size > 0 {
			limit += "," + strconv.Itoa(size)
		} else {
			limit += ",30"
		}
	}
	list, err := app.MD("main").Select(source, false, fields, where, groupby, orderby, limit)
	if err != nil {
		return nil, err
	}
	totalPage, totalRow := 0, 0
	if page > 0 {
		justOne := false
		if groupby == "" {
			justOne = true
		}
		rows, err := app.MD("main").Select(source, justOne, "COUNT(*) count", where, groupby)
		if err != nil {
			return nil, err
		}
		if groupby != "" {
			totalRow = len(rows.([]map[string]string))
		} else {
			totalRow, _ = strconv.Atoi(rows.(map[string]string)["count"])
		}
		totalPage = int(math.Ceil(float64(totalRow) / float64(size)))
	}
	pager := map[string]int{"page": page, "size": size, "total_page": totalPage, "total_row": totalPage}
	result := map[string]interface{}{"list": list, "pager": pager}
	return result, nil
}
