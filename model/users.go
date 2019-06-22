package model

import "alopex/app"

func init() {
	app.Tables["user_invites"] = app.RV(UserInvites{})
	app.Tables["user_bank_cards"] = app.RV(UserBankCards{})
	app.Tables["user_login_logs"] = app.RV(UserLoginLogs{})
	app.Tables["users"] = app.RV(Users{})
	app.Tables["user_feedbacks"] = app.RV(UserFeedbacks{})
	app.Tables["user_banks"] = app.RV(UserBanks{})
	app.Tables["user_groups"] = app.RV(UserGroups{})
}

type UserInvites struct {
    Id       string `N:"id" X:"not_null,primary,auto_increment"`
    Created  string `N:"created" X:"not_null"`
    FriendId string `N:"friend_id" X:"not_null"`
    UserId   string `N:"user_id" X:"not_null"`
}

type UserBankCards struct {
    Id         string `N:"id" X:"unsigned,not_null,primary,auto_increment"`
    Address    string `N:"address" X:"length=100,not_null"`
    BankName   string `N:"bank_name" X:"length=50,not_null"`
    CardNumber string `N:"card_number" X:"length=50,not_null"`
    Created    string `N:"created" X:"not_null,default='0'"`
    Name       string `N:"name" X:"length=50,not_null"`
    Remark     string `N:"remark" X:"length=100,not_null"`
    Status     string `N:"status" X:"not_null,default='0'"`
    Updated    string `N:"updated" X:"not_null,default='0'"`
    UserId     string `N:"user_id" X:"not_null,default='0'"`
}

type UserLoginLogs struct {
    Id         string `N:"id" X:"not_null,primary,auto_increment"`
    Addr       string `N:"addr" X:"length=150,not_null"`
    Ip         string `N:"ip" X:"not_null"`
    LoginFrom  string `N:"login_from" X:"length=20,not_null"`
    LoginTime  string `N:"login_time" X:"not_null,default='0'"`
    LogoutTime string `N:"logout_time" X:"not_null,default='0'"`
    UserId     string `N:"user_id" X:"not_null"`
}

type Users struct {
    Id              string `N:"id" X:"not_null,primary,auto_increment"`
    Birthday        string `N:"birthday" X:"length=10,not_null,default='1970-01-01'"`
    Created         string `N:"created" X:"not_null"`
    DownIds         string `N:"downIds" X:"length=500,default=''"`
    Email           string `N:"email" X:"length=50,not_null"`
    IsModifiedUname string `N:"is_modified_uname" X:"not_null,default='0'"`
    LastLoginTime   string `N:"last_login_time" X:"not_null,default='0'"`
    LastPlatformId  string `N:"last_platform_id" X:"default='0'"`
    MobileType      string `N:"mobile_type" X:"not_null,default='1'"`
    Name            string `N:"name" X:"not_null"`
    ParentId        string `N:"parent_id" X:"not_null,default='0'"`
    Password        string `N:"password" X:"not_null"`
    Path            string `N:"path" X:"default=''"`
    Phone           string `N:"phone" X:"not_null"`
    ProxyStatus     string `N:"proxy_status" X:"default='1'"`
    Qq              string `N:"qq" X:"length=15,not_null"`
    RegIp           string `N:"reg_ip" X:"length=15,not_null"`
    SafePassword    string `N:"safe_password" X:"length=40,not_null"`
    Sex             string `N:"sex" X:"not_null,default='1'"`
    Status          string `N:"status" X:"not_null,default='1'"`
    Token           string `N:"token" X:"length=200,not_null"`
    TokenCreated    string `N:"token_created" X:"not_null"`
    UniqueCode      string `N:"unique_code" X:"length=50,not_null"`
    Updated         string `N:"updated" X:"timestamp,default=current_timestamp,on_update_current_timestamp"`
    UserGroupId     string `N:"user_group_id" X:"not_null"`
    UserName        string `N:"user_name" X:"length=20,not_null,unique"`
    UserType        string `N:"user_type" X:"not_null,default='0'"`
    VipLevel        string `N:"vip_level" X:"not_null,default='1'"`
    Wechat          string `N:"wechat" X:"length=30,not_null"`
    WxOpenId        string `N:"wx_open_id" X:"length=32,not_null"`
}

type UserFeedbacks struct {
    Id       string `N:"id" X:"not_null,primary,auto_increment"`
    Created  string `N:"created" X:"default=''"`
    ImageUrl string `N:"image_url" X:"length=500,default=''"`
    Suggest  string `N:"suggest" X:"length=400,not_null"`
    UserId   string `N:"user_id" X:"not_null"`
    Version  string `N:"version" X:"length=100,default=''"`
}

type UserBanks struct {
    Id   string `N:"id" X:"not_null,primary,auto_increment"`
    Logo string `N:"logo" X:"length=255,not_null"`
    Name string `N:"name" X:"length=30,not_null"`
}

type UserGroups struct {
    Id        string `N:"id" X:"not_null,primary,auto_increment"`
    GroupName string `N:"group_name" X:"length=30,not_null"`
    IsDefault string `N:"is_default" X:"not_null,default='0'"`
    Remark    string `N:"remark" X:"length=255,not_null"`
}