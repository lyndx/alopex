package app

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/grsmv/inflect"
	_ "github.com/mattn/go-sqlite3"
)

type Model struct {
	db *sql.DB
}

var dbs map[string]map[string]*sql.DB

func init() {
	t, err := String("database").C()
	if err != nil {
		DIE("数据库初始化失败")
	}
	dbs = map[string]map[string]*sql.DB{"mysql": {}, "sqlite": {}}
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
		db, err := sql.Open("mysql", cs)
		if err != nil {
			DIE("数据库Mysql[" + key + "]连接失败，" + err.Error())
		}
		if err = db.Ping(); err != nil {
			DIE("数据库Mysql[" + key + "]Ping失败" + err.Error())
		}
		db.SetMaxIdleConns(100)
		db.SetMaxOpenConns(1000)
		db.SetConnMaxLifetime(100 * time.Second)
		dbs["mysql"][key] = db
	}
	// SQLite初始化
	for key, file := range SQLite {
		db, err := sql.Open("sqlite3", file)
		if err != nil {
			DIE("数据库Sqlite[" + key + "]连接失败，" + err.Error())
		}
		if err = db.Ping(); err != nil {
			DIE("数据库Sqlite[" + key + "]Ping失败" + err.Error())
		}
		if IsDeveloper {
			fmt.Println("连接数据库SQLite[" + key + "]：" + file)
		}
		dbs["sqlite"][key] = db
	}
}

// 获取数据库连接
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

// 检查表是否存在
func (m *Model) HasTable(table string) (bool, error) {
	if table == "" {
		return false, errors.New("表名不能为空")
	}
	if m.db == nil {
		return false, errors.New("数据库连接失效")
	}
	db := m.db
	row := db.QueryRow("SHOW TABLES LIKE '" + table + "'")
	var tb string
	if err := row.Scan(&tb); err != nil {
		return false, err
	}
	if tb == table {
		return true, nil
	}
	return false, errors.New("表不存在")
}

// 反转DDL为模型构造文件
func (m *Model) TM(ddl map[string]map[string]map[string]map[string]string) {
	//if file == "" {
	//	DIE("文件错误.....")
	//}
	//file = inflect.Pluralize(file)
	//if String("model").IsExist(file+".go") {
	//
	//}
	//s := Pluralize("ip")
	//fmt.Println(s)
}

// 查询
func (m *Model) Select(table string, args ...interface{}) (interface{}, error) {
	_, err := m.HasTable(table)
	if err != nil {
		return nil, err
	}
	db := m.db
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	tables := make(map[string]map[string]map[string]string)
	for rows.Next() {
		tb := ""
		if err := rows.Scan(&tb); err != nil {
			return nil, err
		}
		if tb != "" {
			tables[tb] = make(map[string]map[string]string)
			rs, err := db.Query("DESC " + tb)
			if err != nil {
				return nil, err
			}
			if err := rs.Err(); err != nil {
				return nil, err
			}
			for rs.Next() {
				f, t, n, k, d, e := "", "", "", "", sql.NullString{}, ""
				if err := rs.Scan(&f, &t, &n, &k, &d, &e); err != nil {
					return nil, err
				}
				if f != "" {
					fd := map[string]string{
						"type":    strings.ToUpper(t),
						"is_null": strings.ToUpper(n),
						"key":     strings.ToUpper(k),
						"default": d.String,
						"extra":   strings.ToUpper(e),
					}
					tables[tb][f] = fd
				}
			}
		}
	}
	models := map[string]map[string]map[string]map[string]string{}
	for tb, ddl := range tables {
		ws := String(tb).Split("_")
		if len(ws) == 1 {
			key := strings.ToLower(inflect.Pluralize(tb))
			tb = String(key).UFrist()
			models[key] = map[string]map[string]map[string]string{tb: ddl}
		} else {
			tb = String(inflect.CamelizeDownFirst(tb)).UFrist()
			key := strings.ToLower(inflect.Pluralize(ws[0]))
			if _, ok := models[key]; ok {
				models[key][tb] = ddl
			} else {
				models[key] = map[string]map[string]map[string]string{tb: ddl}
			}
		}
	}
	rows.Close()
	m.TM(models)
	return models, nil
}
