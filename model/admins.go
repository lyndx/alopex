package model

import "alopex/app"

func init() {
	app.Tables["alopex.admin_roles"] = app.RT(AdminRoles{})
	app.Tables["alopex.admin_rules"] = app.RT(AdminRules{})
	app.Tables["alopex.admins"] = app.RT(Admins{})
}

type AdminRoles struct {
	Id        int64  `N:"id" X:"width=11,unsigned,not_null,primary,auto_increment" M:"编号"`
	Brief     string `N:"brief" X:"length=100,default=''" M:"角色简介"`
	Name      string `N:"name" X:"length=100,not_null,unique" M:"角色名称"`
	RuleIds   string `N:"rule_ids" X:"length=1000,not_null" M:"权限编号集合，英文逗号分割"`
	CreatedAt string `N:"created_at" X:"timestamp,default=current_timestamp" M:"添加时间"`
	UpdatedAt string `N:"updated_at" X:"timestamp,default=current_timestamp,on_update_current_timestamp" M:"更新时间"`
	DeletedAt string `N:"deleted_at" X:"timestamp" M:"删除时间"`
}

type AdminRules struct {
	Id        int64  `N:"id" X:"width=11,unsigned,not_null,primary,auto_increment" M:"编号"`
	Name      string `N:"name" X:"length=50,not_null" M:"权限名称"`
	Route     string `N:"route" X:"length=100,not_null,unique" M:"权限路由"`
	Sort      int64  `N:"sort" X:"width=11,unsigned,default=0" M:"排序序号"`
	Type      int64  `N:"type" X:"width=1,unsigned,not_null" M:"权限分类（1左侧导航，2Tab导航，3动作事件）"`
	CreatedAt string `N:"created_at" X:"timestamp,default=current_timestamp" M:"添加时间"`
	UpdatedAt string `N:"updated_at" X:"timestamp,default=current_timestamp,on_update_current_timestamp" M:"更新时间"`
	DeletedAt string `N:"deleted_at" X:"timestamp" M:"删除时间"`
}

type Admins struct {
	Id        int64  `N:"id" X:"width=11,unsigned,not_null,primary,auto_increment" M:"编号"`
	Email     string `N:"email" X:"length=100,not_null,unique" M:"邮件"`
	Password  string `N:"password" X:"length=100,not_null" M:"密码"`
	RoleIds   string `N:"role_ids" X:"length=100,not_null" M:"角色编号集合，英文逗号分隔"`
	Status    int64  `N:"status" X:"width=1,unsigned,default=1" M:"状态（0冻结，1正常）"`
	Token     string `N:"token" X:"length=1000,default=''" M:"最后登录的TOKEN"`
	Username  string `N:"username" X:"length=100,not_null,unique" M:"帐号"`
	CreatedAt string `N:"created_at" X:"timestamp,default=current_timestamp" M:"添加时间"`
	UpdatedAt string `N:"updated_at" X:"timestamp,default=current_timestamp,on_update_current_timestamp" M:"更新时间"`
	DeletedAt string `N:"deleted_at" X:"timestamp" M:"删除时间"`
}
