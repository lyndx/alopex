package model

type AdminRules struct {
	CreatedAt string `N:"created_at" X:"timestamp,default=current_timestamp"`
	Id        string `N:"id" X:"unsigned,not_null,primary,auto_increment"`
	IsMenu    string `N:"is_menu" X:"unsigned,default='0'"`
	Module    string `N:"module" X:"not_null"`
	Name      string `N:"name" X:"not_null,unique"`
	Route     string `N:"route" X:"not_null"`
	UpdatedAt string `N:"updated_at" X:"timestamp,default=current_timestamp,on_update_current_timestamp"`
}

type Admins struct {
	CreatedAt      string `N:"created_at" X:"timestamp,default=current_timestamp"`
	Email          string `N:"email" X:"default=''"`
	Id             string `N:"id" X:"unsigned,not_null,primary,auto_increment"`
	LastLoginIp    string `N:"last_login_ip" X:"default=''"`
	LastLoginToken string `N:"last_login_token" X:"default=''"`
	Password       string `N:"password" X:"not_null"`
	RegisterIp     string `N:"register_ip" X:"not_null"`
	RoleIds        string `N:"role_ids" X:"not_null"`
	Status         string `N:"status" X:"unsigned,default='1'"`
	UpdatedAt      string `N:"updated_at" X:"timestamp,default=current_timestamp,on_update_current_timestamp"`
	Username       string `N:"username" X:"not_null,unique"`
}

type AdminRoles struct {
	CreatedAt string `N:"created_at" X:"timestamp,default=current_timestamp"`
	Id        string `N:"id" X:"unsigned,not_null,primary,auto_increment"`
	Name      string `N:"name" X:"not_null,unique"`
	RuleIds   string `N:"rule_ids" X:"not_null"`
	Status    string `N:"status" X:"unsigned,default='1'"`
	UpdatedAt string `N:"updated_at" X:"timestamp,default=current_timestamp,on_update_current_timestamp"`
}
