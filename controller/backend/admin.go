package backend

import (
	"strconv"
	"strings"

	"alopex/app"
	"alopex/service"
)

type AdminController struct{}

func init() {
	app.CJoin("admin", AdminController{})
}

func (ctrl AdminController) List(h *app.Http) {
	page, _ := strconv.Atoi(h.P("page").(string))
	size, _ := strconv.Atoi(h.P("size").(string))
	filter, _ := app.String(h.P("filter").(string)).J2O()
	orderby, _ := app.String(h.P("orderby").(string)).J2A()
	where := make([]string, 0)
	for field, v := range filter {
		value := v.(map[string]interface{})
		switch value["type"] {
		case "equal":
			where = append(where, field+"="+value["value"].(string))
		case "like":
			where = append(where, field+" LIKE '%"+value["value"].(string)+"%'")
		case "between":
			a := value["value"].([]interface{})
			where = append(where, field+" BETWEEN '"+a[0].(string)+"' AND '"+a[1].(string)+"'")
		}
	}
	conditions := map[string]string{
		"source": "admins",
		"fields": "id,username,password,email,date_format(created_at,'%Y/%m/%d %H:%m:%s') created_at,status",
		"where":  strings.Join(where, " AND "),
	}
	if len(orderby) > 0 {
		conditions["orderby"] = orderby[0].(string) + " " + orderby[1].(string)
	}
	as := service.AdminService{}
	//
	result, err := as.GetAdminList(conditions, page, size)
	if err != nil {
		h.Output(402, "管理员列表获取失败")
	}
	h.Output(200, result, "加载成功")
}
