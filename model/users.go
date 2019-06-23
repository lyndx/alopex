package model

import "alopex/app"

func init() {
	app.Tables["main.user_details"] = app.RT(UserDetails{})
	app.Tables["main.users"] = app.RT(Users{})
}

type UserDetails struct {
	Id        int64   `N:"id" X:"width=11,unsigned,not_null,primary,auto_increment" M:"编号"`
	Balance   float64 `N:"balance" X:"width=11,prec=2,unsigned,default=0.00" M:"余额"`
	Brief     string  `N:"brief" X:"length=-2,default=''" M:"简介"`
	UserId    int64   `N:"user_id" X:"width=11,unsigned,not_null" M:"用户编号"`
	CreatedAt string  `N:"created_at" X:"timestamp,default=current_timestamp" M:"添加时间"`
	UpdatedAt string  `N:"updated_at" X:"timestamp,default=current_timestamp,on_update_current_timestamp" M:"更新时间"`
	DeletedAt string  `N:"deleted_at" X:"timestamp" M:"删除时间"`
}

type Users struct {
	Id        int64  `N:"id" X:"width=11,unsigned,not_null,primary,auto_increment" M:"编号"`
	Email     string `N:"email" X:"length=50,default=''" M:"邮箱"`
	Mobile    string `N:"mobile" X:"length=11,default=''" M:"手机号"`
	Nickname  string `N:"nickname" X:"length=100,default=''" M:"昵称"`
	Password  string `N:"password" X:"length=100,not_null" M:"密码"`
	Status    int64  `N:"status" X:"width=1,unsigned,default=0" M:"状态(0冻结,1正常)"`
	Username  string `N:"username" X:"length=100,not_null,unique" M:"用户名"`
	CreatedAt string `N:"created_at" X:"timestamp,default=current_timestamp" M:"添加时间"`
	UpdatedAt string `N:"updated_at" X:"timestamp,default=current_timestamp,on_update_current_timestamp" M:"更新时间"`
	DeletedAt string `N:"deleted_at" X:"timestamp" M:"删除时间"`
}
