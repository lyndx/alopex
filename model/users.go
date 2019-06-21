package model

type Users struct {
	CreatedAt      string `N:"created_at" X:"timestamp,default=current_timestamp"`
	Id             string `N:"id" X:"unsigned,not_null,primary,auto_increment"`
	LastLoginIp    string `N:"last_login_ip" X:"default=''"`
	LastLoginToken string `N:"last_login_token" X:"default=''"`
	Password       string `N:"password" X:"not_null"`
	RegisterIp     string `N:"register_ip" X:"not_null,default='0'"`
	Status         string `N:"status" X:"unsigned,default='1.000'"`
	UpdatedAt      string `N:"updated_at" X:"timestamp,default=current_timestamp,on_update_current_timestamp"`
	Username       string `N:"username" X:"not_null,unique"`
}
