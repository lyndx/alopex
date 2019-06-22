package app

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/grsmv/inflect"
	_ "github.com/mattn/go-sqlite3"
)

type Model struct {
	db *sql.DB
}

var dbs map[string]map[string]*sql.DB

var Tables map[string]reflect.Value

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
			Dump("yellow", "连接数据库Mysql["+key+"]："+cs)
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
			Dump("yellow", "连接数据库SQLite["+key+"]："+file)
		}
		dbs["sqlite"][key] = db
	}
	// 表初始化
	Tables = make(map[string]reflect.Value)
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

// 生成模型构造
func (m *Model) TS(mn string, fs map[string]map[string]string) string {
	str := "type " + String(mn).UFrist() + " struct {\n"
	items, lfs := make([]string, 0), ""
	ID, CT, UT, DT := "", "", "", ""
	for f, fd := range fs {
		item := "%s@%s `N:\"%s\" X:\"%s\"`"
		a, b, c, d := String(CamelizeDownFirst(f)).UFrist(), "string", f, make([]string, 0)
		if len(lfs) < len(a) {
			lfs = a
		}
		ft := String(fd["type"]).Split(" ")
		if strings.HasSuffix(ft[0], "INT") {
			b = "int"
			d = append(d, "width="+strings.Replace(String(ft[0]).Split("(")[1], ")", "", -1))
		} else if ft[0] == "DECIMAL" {
			b = "float"
			x := String(strings.Replace(String(ft[0]).Split("(")[1], ")", "", -1)).Split(",")
			d = append(d, "width="+x[0], "prec="+x[1])
		} else if strings.HasPrefix(ft[0], "VARCHAR") {
			d = append(d, "length="+strings.Replace(String(ft[0]).Split("(")[1], ")", "", -1))
		} else if ft[0] == "TIMESTAMP" {
			d = append(d, "timestamp")
		}
		if len(ft) == 2 {
			d = append(d, "unsigned")
		}
		if fd["is_null"] == "NO" {
			d = append(d, "not_null")
		}
		if fd["key"] == "PRI" {
			d = append(d, "primary")
		}
		if fd["key"] == "UNI" {
			d = append(d, "unique")
		}
		if fd["default"] != "" {
			def := fd["default"]
			if def == "CURRENT_TIMESTAMP" {
				def = "current_timestamp"
			}
			if (b == "string") && (def != "current_timestamp") {
				def = "'" + def + "'"
			}
			d = append(d, "default="+def)
		} else if fd["is_null"] == "YES" {
			def := fd["default"]
			if def == "CURRENT_TIMESTAMP" {
				def = "current_timestamp"
			}
			if (b == "string") && (def != "current_timestamp") {
				def = "'" + def + "'"
			}
			d = append(d, "default="+def)
		}
		if fd["extra"] == "AUTO_INCREMENT" {
			d = append(d, "auto_increment")
		} else if fd["extra"] == "ON UPDATE CURRENT_TIMESTAMP" {
			d = append(d, "on_update_current_timestamp")
		}
		sort.Reverse(sort.StringSlice(d))
		item = fmt.Sprintf(item, a, b, c, strings.Join(d, ",")) + "\n"
		if a == "Id" {
			ID = item
		} else if a == "CreatedAt" {
			CT = item
		} else if a == "UpdatedAt" {
			UT = item
		} else if a == "DeletedAt" {
			DT = item
		} else {
			items = append(items, item)
		}
	}
	CUD := make([]string, 0)
	if CT != "" {
		CUD = append(CUD, CT)
	}
	if UT != "" {
		CUD = append(CUD, UT)
	}
	if DT != "" {
		CUD = append(CUD, DT)
	}
	sort.Sort(sort.StringSlice(items))
	if ID != "" {
		tmp := String(ID).Split("@")
		f, length, max := tmp[0], len(tmp[0]), len(lfs)
		for i := 0; i < (max-length)+1; i++ {
			f += " "
		}
		str += "    " + f + tmp[1]
	}
	for _, v := range items {
		tmp := String(v).Split("@")
		f, length, max := tmp[0], len(tmp[0]), len(lfs)
		for i := 0; i < (max-length)+1; i++ {
			f += " "
		}
		str += "    " + f + tmp[1]
	}
	for _, v := range CUD {
		tmp := String(v).Split("@")
		f, length, max := tmp[0], len(tmp[0]), len(lfs)
		for i := 0; i < (max-length)+1; i++ {
			f += " "
		}
		str += "    " + f + tmp[1]
	}
	return str + "}"
}

// 获取表信息
func (m *Model) getTablesDDL() (map[string]map[string]map[string]string, error) {
	db := m.db
	if db == nil {
		return nil, errors.New("数据库连接失败.....")
	}
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
		rs, err := db.Query("DESC " + tb)
		if err != nil {
			return nil, err
		}
		if err := rs.Err(); err != nil {
			return nil, err
		}
		fs := make(map[string]map[string]string)
		for rs.Next() {
			f, t, n, k, d, e := "", "", "", "", sql.NullString{}, ""
			if err := rs.Scan(&f, &t, &n, &k, &d, &e); err != nil {
				return nil, err
			}
			if f != "" {
				fs[f] = map[string]string{
					"type":    strings.ToUpper(t),
					"is_null": strings.ToUpper(n),
					"key":     strings.ToUpper(k),
					"default": d.String,
					"extra":   strings.ToUpper(e),
				}
			}
		}
		tables[tb] = fs
	}
	rows.Close()
	return tables, nil
}

// 反转DDL为模型构造文件
func (m *Model) TM() (bool, error) {
	db := m.db
	if db == nil {
		return false, errors.New("数据库连接失败.....")
	}
	tables, err := m.getTablesDDL()
	if err != nil {
		return false, err
	}
	if len(tables) < 1 {
		return false, errors.New("表数据解析失败.....")
	}
	fds, mapper := make(map[string][]string), make(map[string][]string)
	for tb, fs := range tables {
		f := strings.ToLower(Pluralize(tb))
		ws := String(tb).Split("_")
		if len(ws) > 1 {
			f = strings.ToLower(Pluralize(ws[0]))
		}
		if _, ok := fds[f]; !ok {
			fds[f] = make([]string, 0)
			mapper[f] = make([]string, 0)
		}
		mn := CamelizeDownFirst(tb)
		mapper[f] = append(mapper[f], mn)
		fds[f] = append(fds[f], m.TS(mn, fs))
	}
	for file, ctx := range fds {
		function := "func init() {\n"
		for _, fm := range mapper[file] {
			function += "	app.Tables[\"" + Underscore(fm) + "\"] = app.RV(" + String(fm).UFrist() + "{})\n"
		}
		function += "}\n\n"
		file = file + ".go"
		if String("model").IsExist(file) {
			if err := os.Remove("model/" + file); err != nil {
				return false, err
			}
		}
		_, err := String("model/" + file).Write("package model\n\nimport \"alopex/app\"\n\n" + function + strings.Join(ctx, "\n\n"))
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// 查询
func (m *Model) Select(table string, args ...interface{}) (interface{}, error) {
	_, err := m.HasTable(table)
	if err != nil {
		return nil, err
	}

	// do more ...
	return nil, nil
}
