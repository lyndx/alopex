package model

import "alopex/app"

func init() {
	app.Tables["admin_accesses"] = app.RV(AdminAccesses{})
	app.Tables["admin_roles"] = app.RV(AdminRoles{})
	app.Tables["admin_ips"] = app.RV(AdminIps{})
	app.Tables["admin_role_sets"] = app.RV(AdminRoleSets{})
	app.Tables["admins"] = app.RV(Admins{})
	app.Tables["admin_login_logs"] = app.RV(AdminLoginLogs{})
	app.Tables["admin_nodes"] = app.RV(AdminNodes{})
	app.Tables["admin_logs"] = app.RV(AdminLogs{})
}

type AdminAccesses struct {
    Id           string `N:"id" X:"unsigned,not_null,primary,auto_increment"`
    AccessMethod string `N:"access_method" X:"not_null"`
    Level        string `N:"level" X:"not_null"`
    Module       string `N:"module" X:"unsigned,not_null"`
    NodeId       string `N:"node_id" X:"unsigned,not_null"`
    RoleId       string `N:"role_id" X:"unsigned,not_null"`
}

type AdminRoles struct {
    Id       string `N:"id" X:"unsigned,not_null,primary,auto_increment"`
    MenuIds  string `N:"menu_ids" X:"default=''"`
    Name     string `N:"name" X:"length=20,not_null"`
    ParentId string `N:"parent_id" X:"default=''"`
    Remark   string `N:"remark" X:"length=255,default=''"`
    Status   string `N:"status" X:"unsigned,default=''"`
}

type AdminIps struct {
    Id string `N:"id" X:"not_null,primary,auto_increment"`
    Ip string `N:"ip" X:"length=20,not_null"`
}

type AdminRoleSets struct {
    Id      string `N:"id" X:"unsigned,not_null,primary,auto_increment"`
    AdminId string `N:"admin_id" X:"default=''"`
    RoleId  string `N:"role_id" X:"unsigned,default=''"`
}

type Admins struct {
    Id            string `N:"id" X:"unsigned,not_null,primary,auto_increment"`
    ChargeAlert   string `N:"charge_alert" X:"not_null,default='1'"`
    Created       string `N:"created" X:"not_null"`
    Email         string `N:"email" X:"length=255,not_null"`
    ForceOut      string `N:"force_out" X:"not_null,default='0'"`
    IsOtp         string `N:"is_otp" X:"not_null,default='0'"`
    IsOtpFirst    string `N:"is_otp_first" X:"not_null,default='1'"`
    LoginIp       string `N:"login_ip" X:"length=1000,not_null"`
    ManualMax     string `N:"manual_max" X:"not_null,default='0'"`
    Name          string `N:"name" X:"length=255,not_null"`
    Password      string `N:"password" X:"length=255,not_null"`
    Permission    string `N:"permission" X:"not_null,default='1'"`
    RoleId        string `N:"role_id" X:"not_null"`
    Status        string `N:"status" X:"not_null"`
    Token         string `N:"token" X:"default=''"`
    Updated       string `N:"updated" X:"not_null,default='0'"`
    WithdrawAlert string `N:"withdraw_alert" X:"not_null,default='0'"`
}

type AdminLoginLogs struct {
    Id        string `N:"id" X:"unsigned,not_null,primary,auto_increment"`
    AdminId   string `N:"admin_id" X:"not_null"`
    AdminName string `N:"admin_name" X:"length=50,default=''"`
    Ip        string `N:"ip" X:"not_null"`
    LoginTime string `N:"login_time" X:"not_null,default='0'"`
}

type AdminNodes struct {
    Id       string `N:"id" X:"unsigned,not_null,primary,auto_increment"`
    Level    string `N:"level" X:"unsigned,not_null"`
    Method   string `N:"method" X:"not_null"`
    Name     string `N:"name" X:"length=50,not_null"`
    ParentId string `N:"parent_id" X:"unsigned,not_null"`
    Remark   string `N:"remark" X:"length=255,default=''"`
    Route    string `N:"route" X:"length=100,not_null"`
    Seq      string `N:"seq" X:"unsigned,default=''"`
    Status   string `N:"status" X:"default='0'"`
    Title    string `N:"title" X:"length=50,default=''"`
    Type     string `N:"type" X:"default=''"`
}

type AdminLogs struct {
    Id        string `N:"id" X:"not_null,primary,auto_increment"`
    AdminId   string `N:"admin_id" X:"not_null"`
    AdminName string `N:"admin_name" X:"length=50,default=''"`
    Content   string `N:"content" X:"length=255,not_null"`
    Created   string `N:"created" X:"default='0'"`
    Node      string `N:"node" X:"length=50,not_null"`
    Type      string `N:"type" X:"length=255,not_null"`
}