package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-ffmt/ffmt"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Model struct {
	db *gorm.DB
}

var dbs map[string]map[string]*gorm.DB

func init() {
	t, err := String("database").C()
	if err != nil {
		DIE("数据库初始化失败")
	}
	dbs = map[string]map[string]*gorm.DB{"mysql": {}, "sqlite": {}}
	MySQL, SQLite := map[string]map[string]string{}, map[string]string{}
	for k, v := range TValue(t).(map[string]interface{}) {
		kk := String(k).Split(".")
		if len(kk) > 1 {
			switch kk[0] {
			case "mysql":
				key, field := strings.Join(kk[1:len(kk)-1], "."), kk[len(kk)-1]
				if _, ok := MySQL[key]; !ok {
					MySQL[key] = make(map[string]string)
				}
				vv := TT(v)
				if vv.IsValid() {
					MySQL[key][field] = vv.ToString()
				}
			case "sqlite":
				key, vv := kk[1], TT(v)
				if vv.IsValid() {
					SQLite[key] = vv.ToString()
				}
			}
		}
	}
	// 是否为开发模式
	IsDeveloper := true
	tmp, err := String("app").C("is_developer")
	if (err == nil) && tmp.IsBool() && !TValue(tmp, true).(bool) {
		IsDeveloper = false
	}
	// Mysql初始化
	for key, cfg := range MySQL {
		if TT(cfg["host"]).IsEmpty() || TT(cfg["port"]).IsEmpty() || TT(cfg["username"]).IsEmpty() || TT(cfg["password"]).IsEmpty() || TT(cfg["database"]).IsEmpty() || TT(cfg["charset"]).IsEmpty() || TT(cfg["parsetime"]).IsEmpty() || TT(cfg["loc"]).IsEmpty() {
			DIE("数据库Mysql[" + key + "]配置项不完整")
		}
		cs := cfg["username"] + ":" + cfg["password"] + "@tcp(" + cfg["host"] + ":" + cfg["port"] + ")/" + cfg["database"] + "?charset=" + cfg["charset"] + "&parseTime=" + cfg["parsetime"] + "&loc=" + cfg["loc"]
		if IsDeveloper {
			fmt.Println("连接数据库Mysql[" + key + "]：" + cs)
		}
		db, err := gorm.Open("mysql", cs)
		if err != nil {
			DIE("数据库Mysql[" + key + "]连接失败，" + err.Error())
		}
		if err = db.DB().Ping(); err != nil {
			DIE("数据库Mysql[" + key + "]Ping失败" + err.Error())
		}
		db.LogMode(IsDeveloper)
		db.DB().SetMaxIdleConns(10)
		db.DB().SetMaxOpenConns(100)
		dbs["mysql"][key] = db
	}
	// SQLite初始化
	for key, file := range SQLite {
		db, err := gorm.Open("sqlite3", file)
		if err != nil {
			DIE("数据库Sqlite[" + key + "]连接失败，" + err.Error())
		}
		fmt.Println("连接数据库SQLite[" + key + "]：" + file)
		db.LogMode(IsDeveloper)
		db.DB().SetMaxIdleConns(10)
		db.DB().SetMaxOpenConns(100)
		dbs["sqlite"][key] = db
	}
}

func MD(key string) *Model {
	m := Model{}
	if TT(key).IsEmpty() {
		return &m
	}
	kk := String(key).Split(".")
	if len(kk) < 2 {
		return &m
	}
	dt := strings.ToLower(kk[0])
	if _, ok := dbs[dt]; !ok {
		return &m
	}
	dd := strings.Join(kk[1:], ".")
	if _, ok := dbs[dt][dd]; !ok {
		return &m
	}
	m = Model{dbs[dt][dd]}
	return &m
}

type X struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default interface{}
	Extra   string
}

func (m *Model) Select(table string, args ...interface{}) (interface{}, error) {
	if m.db == nil {
		return nil, errors.New("数据库连接失败")
	}
	if TT(table).IsEmpty() {
		return nil, errors.New("数据表名不能为空")
	}
	db := m.db
	if !db.HasTable(table) {
		return nil, errors.New("数据表不存在")
	}
	tables := make([]string, 0)
	db.Raw("SHOW TABLES").Pluck("", &tables)

	for _, tb := range tables {

		ddl := make([]string, 0)
		db.Raw("select * from information_schema.columns where table_name='"+tb+"'").Pluck("column_name", &ddl)
		ffmt.Puts(ddl)
	}

	return tables, nil
}
