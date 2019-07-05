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
)

type Model struct {
	dbname string
	db     *sql.DB
}

var dbs map[string]*sql.DB

var Tables map[string]reflect.Type
var Services map[string]interface{}

func init() {
	Dump()
	t, err := String("database").C()
	if err != nil {
		DIE("数据库初始化失败")
	}
	dbs = make(map[string]*sql.DB, 0)
	configs := make(map[string]map[string]string)
	for k, v := range TValue(t).(map[string]interface{}) {
		kk := String(k).Split(".")
		if len(kk) != 2 {
			DIE("数据库配置错误")
		}
		vv := TT(v)
		if !vv.IsValid() {
			DIE("数据库配置错误")
		}
		key, item := strings.ToLower(kk[0]), strings.ToLower(kk[1])
		if _, ok := configs[key]; !ok {
			configs[key] = make(map[string]string)
		}
		configs[key][item] = vv.ToString()
	}
	// 初始化
	for key, cfg := range configs {
		if TT(cfg["host"]).IsEmpty() || TT(cfg["port"]).IsEmpty() || TT(cfg["username"]).IsEmpty() || TT(cfg["password"]).IsEmpty() || TT(cfg["database"]).IsEmpty() || TT(cfg["charset"]).IsEmpty() || TT(cfg["parsetime"]).IsEmpty() || TT(cfg["loc"]).IsEmpty() {
			DIE("数据库Mysql[" + key + "]配置项不完整")
		}
		cs := cfg["username"] + ":" + cfg["password"] + "@tcp(" + cfg["host"] + ":" + cfg["port"] + ")/" + cfg["database"] + "?charset=" + cfg["charset"] + "&parseTime=" + cfg["parsetime"] + "&loc=" + cfg["loc"]
		if IsDeveloper {
			Dump("yellow", "连接数据库["+key+"]："+cs)
		}
		db, err := sql.Open("mysql", cs)
		if err != nil {
			DIE("数据库[" + key + "]连接失败，" + err.Error())
		}
		if err = db.Ping(); err != nil {
			DIE("数据库[" + key + "]Ping失败" + err.Error())
		}
		db.SetMaxIdleConns(100)
		db.SetMaxOpenConns(1000)
		db.SetConnMaxLifetime(100 * time.Second)
		dbs[key] = db
	}
	//
	Tables = make(map[string]reflect.Type)
	Services = make(map[string]interface{})
}

// 获取数据库连接
func MD(key string) *Model {
	m := Model{}
	if TT(key).IsEmpty() {
		return &m
	}
	key = strings.ToLower(key)
	if _, ok := dbs[key]; !ok {
		return &m
	}
	dbname, err := String("database").C(key + ".database")
	if err != nil {
		return &m
	}
	m = Model{dbname.ToString(), dbs[key]}
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
		item := "%s@%s `N:\"%s\" X:\"%s\" M:\"%s\"`\n"
		a, b, c, d, e := String(CamelizeDownFirst(f)).UFrist(), "string", f, make([]string, 0), fd["comment"]
		if len(lfs) < len(a) {
			lfs = a
		}
		ft := String(fd["type"]).Split(" ")
		ftf := String(ft[0]).Split("(")
		if strings.HasSuffix(ftf[0], "INT") {
			b = "int64"
			d = append(d, "width="+strings.Replace(ftf[1], ")", "", -1))
		} else if strings.HasPrefix(ftf[0], "DECIMAL") {
			b = "float64"
			x := String(strings.Replace(ftf[1], ")", "", -1)).Split(",")
			d = append(d, "width="+x[0], "prec="+x[1])
		} else if strings.HasPrefix(ftf[0], "VARCHAR") {
			d = append(d, "length="+strings.Replace(ftf[1], ")", "", -1))
		} else if ft[0] == "TIMESTAMP" {
			d = append(d, "timestamp")
		} else if ft[0] == "TEXT" {
			d = append(d, "length=-1")
		} else if ft[0] == "LONGTEXT" {
			d = append(d, "length=-2")
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
		if a == "Id" {
			ID = fmt.Sprintf(item, "Id", "int64", "id", "width=11,unsigned,not_null,primary,auto_increment", "编号")
		} else if a == "CreatedAt" {
			CT = fmt.Sprintf(item, "CreatedAt", "string", "created_at", "timestamp,default=current_timestamp", "添加时间")
		} else if a == "UpdatedAt" {
			UT = fmt.Sprintf(item, "UpdatedAt", "string", "updated_at", "timestamp,default=current_timestamp,on_update_current_timestamp", "更新时间")
		} else if a == "DeletedAt" {
			DT = fmt.Sprintf(item, "DeletedAt", "string", "deleted_at", "timestamp", "删除时间")
		} else {
			_ = sort.Reverse(sort.StringSlice(d))
			items = append(items, fmt.Sprintf(item, a, b, c, strings.Join(d, ","), e))
		}
	}
	sort.Sort(sort.StringSlice(items))
	if CT != "" {
		items = append(items, CT)
	}
	if UT != "" {
		items = append(items, UT)
	}
	if DT != "" {
		items = append(items, DT)
	}
	if ID != "" {
		items = append([]string{ID}, items...)
	}
	for _, v := range items {
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
		rs, err := db.Query("SHOW FULL FIELDS FROM " + tb)
		if err != nil {
			return nil, err
		}
		if err := rs.Err(); err != nil {
			return nil, err
		}
		fs := make(map[string]map[string]string)
		for rs.Next() {
			f, t, co, n, k, d, e, pr, cm := "", "", sql.NullString{}, "", "", sql.NullString{}, "", "", ""
			if err := rs.Scan(&f, &t, &co, &n, &k, &d, &e, &pr, &cm); err != nil {
				return nil, err
			}
			if f != "" {
				fs[f] = map[string]string{
					"type":    strings.ToUpper(t),
					"is_null": strings.ToUpper(n),
					"key":     strings.ToUpper(k),
					"default": d.String,
					"extra":   strings.ToUpper(e),
					"comment": cm,
				}
			}
		}
		tables[tb] = fs
	}
	_ = rows.Close()
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
		return false, errors.New("数据库不能为空.....")
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
			function += "	app.Tables[\"" + m.dbname + "." + Underscore(fm) + "\"] = app.RT(" + String(fm).UFrist() + "{})\n"
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

// 反转模型构造为数据表
func (m *Model) MT() (bool, error) {
	db, dbname := m.db, m.dbname
	if db == nil {
		return false, errors.New("数据库连接失败.....")
	}
	if len(Tables) < 1 {
		return false, errors.New("数据表不能为空.....")
	}
	ts := make([]string, 0)
	for t := range Tables {
		ts = append(ts, t)
	}
	// 清空多余的表
	rows, err := db.Query("SELECT table_name,engine FROM information_schema.tables WHERE table_schema='" + dbname + "'")
	if err != nil {
		return false, err
	}
	if err := rows.Err(); err != nil {
		return false, err
	}
	emapper := make(map[string]string)
	sqls := make([]string, 0)
	for rows.Next() {
		tb, engine := "", ""
		if err := rows.Scan(&tb, &engine); err != nil {
			return false, err
		}
		tb = strings.Replace(tb, m.dbname+".", "", -1)
		emapper[tb] = engine
		sqls = append(sqls, "DROP TABLE IF EXISTS "+tb)
	}
	for _, s := range sqls {
		Dump("blue", s)
		_, err := db.Exec(s)
		if err != nil {
			return false, errors.New("表删除失败....")
		}
	}
	for tn, tb := range Tables {
		if !strings.HasPrefix(tn, m.dbname+".") {
			continue
		}
		length, items, keys := tb.NumField(), make([]string, 0), make([]string, 0)
		if length < 1 {
			return false, errors.New("表字段不能为空....")
		}
		tn = strings.Replace(tn, m.dbname+".", "", -1)
		for i := 0; i < length; i++ {
			tbf := tb.Field(i)
			ft, fgn, fgx, fgm := tbf.Type.String(), tbf.Tag.Get("N"), String(tbf.Tag.Get("X")).Split(","), tbf.Tag.Get("M")
			item, null, ouc, aic := []string{"`" + fgn + "` ", "", ""}, []string{"", "DEFAULT NULL "}, "", ""
			if ft == "int64" {
				item[1] = "INT"
				item[2] = "11"
				null[1] = "DEFAULT 0 "
			} else if ft == "string" {
				item[1] = "VARCHAR"
				item[2] = "200"
				null[1] = "DEFAULT '' "
			} else if ft == "float64" {
				item[1] = "DECIMAL"
				item[2] = "11,2"
				null[1] = "DEFAULT 0.00 "
			}
			for _, v := range fgx {
				if strings.HasPrefix(v, "width=") || strings.HasPrefix(v, "length=") {
					vt := String(item[2]).Split(",")
					vt[0] = String(v).Split("=")[1]
					if vt[0] == "1" {
						if ft == "int64" {
							item[1] = "TINYINT "
						}
					}
					item[2] = strings.Join(vt, ",")
					if (ft == "string") && strings.HasPrefix(v, "length=") {
						if vt[0] == "-1" {
							item[1] = "TEXT "
						} else if vt[0] == "-2" {
							item[1] = "LONGTEXT "
						}
					}
				} else if strings.HasPrefix(v, "perc=") {
					vt := String(item[2]).Split(",")
					if len(vt) == 1 {
						vt = append(vt, "")
					}
					vt[1] = String(v).Split("=")[1]
					item[2] = strings.Join(vt, ",")
				} else if v == "unsigned" {
					item = append(item, "UNSIGNED ")
				} else if v == "not_null" {
					item = append(item, "")
					null[0] = "NOT NULL "
				} else if v == "primary" {
					keys = append(keys, "PRIMARY KEY (`"+fgn+"`)")
				} else if v == "auto_increment" {
					aic = "AUTO_INCREMENT "
				} else if v == "timestamp" {
					item[1] = "TIMESTAMP "
					item[2] = ""
				} else if strings.HasPrefix(v, "default=") {
					tddef := String(v).Split("=")[1]
					if tddef == "current_timestamp" {
						item[1] = "TIMESTAMP "
						tddef = "CURRENT_TIMESTAMP"
					}
					null[1] = "DEFAULT " + tddef + " "
				} else if v == "on_update_current_timestamp" {
					item[1] = "TIMESTAMP "
					ouc = "ON UPDATE CURRENT_TIMESTAMP "
				} else if v == "unique" {
					keys = append(keys, "UNIQUE KEY `index_"+fgn+"` (`"+fgn+"`) USING BTREE")
				}
			}
			if (item[2] != "") && (!strings.HasSuffix(item[1], "TEXT ")) {
				item[1] = item[1] + "(" + item[2] + ") "
				item[2] = ""
			}
			if strings.HasSuffix(item[1], "TEXT ") {
				item[2] = ""
				null = []string{"", "DEFAULT NULL "}
			}
			if null[0] == "NOT NULL " {
				null[1] = ""
			}
			isn := strings.Join(null, "")
			if isn != "" {
				item = append(item, isn)
			}
			if ouc != "" {
				item = append(item, ouc)
			}
			if aic != "" {
				item = append(item, aic)
			}
			if fgm != "" {
				item = append(item, "COMMENT '"+fgm+"'")
			}
			if fgn == "deleted_at" {
				item[3] = "DEFAULT NULL "
			}
			items = append(items, strings.Join(item, ""))
		}
		fieldsStr := "\n    " + items[0]
		if len(items) > 1 {
			fieldsStr = "\n    " + strings.Join(items, ",\n    ")
		}
		keysStr := ""
		if len(keys) > 0 {
			keysStr = keys[0]
			if len(keys) > 1 {
				keysStr = strings.Join(keys, ",\n    ")
			}
			keysStr = ",\n    " + keysStr
		}
		engine := "INNODB"
		if _, ok := emapper[tn]; ok {
			engine = strings.ToUpper(emapper[tn])
		}
		ddl := "CREATE TABLE `" + tn + "` (" + fieldsStr + keysStr + "\n) ENGINE=" + engine + " DEFAULT CHARSET=UTF8"
		if IsDeveloper {
			Dump("yellow", ddl+"\n")
		}
		_, err := db.Exec(ddl)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// 查询 返回
func (m *Model) Select(source string, justOne bool, args ...string) (interface{}, error) {
	if source == "" {
		return nil, errors.New("查询源不能为空")
	}
	template, fields, where := "SELECT %v FROM "+source+" WHERE %v", "*", "1=1"
	if (len(args) > 0) && (args[0] != "") {
		fields = args[0]
	}
	if (len(args) > 1) && (args[1] != "") {
		where = args[1]
	}
	sqlStr := fmt.Sprintf(template, fields, where)
	if (len(args) > 2) && (args[2] != "") {
		sqlStr += fmt.Sprintf(" GROUP BY %v", args[2])
	}
	if (len(args) > 3) && (args[3] != "") {
		sqlStr += fmt.Sprintf(" ORDER BY %v", args[3])
	}
	if (len(args) > 4) && (args[4] != "") {
		sqlStr += fmt.Sprintf(" LIMIT %v", args[4])
	}
	if IsDeveloper {
		Dump("yellow", sqlStr)
	}
	rows, err := m.db.Query(sqlStr)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	cols, _ := rows.Columns()
	result, clen := make([]map[string]string, 0), len(cols)
	for rows.Next() {
		item := make([]interface{}, 0)
		for i := 0; i < clen; i++ {
			x := sql.NullString{}
			item = append(item, &x)
		}
		if err := rows.Scan(item...); err != nil {
			return false, err
		}
		row := make(map[string]string)
		for k, v := range item {
			row[cols[k]] = v.(*sql.NullString).String
		}
		result = append(result, row)
	}
	if justOne {
		if len(result) > 0 {
			return result[0], nil
		}
		return nil, nil
	}
	return result, nil
}

// 数据更新（增删改），返回（操作影响的记录数（如果是单条插入，返回对应插入的ID） & 执行成功与否 & 错误消息）
func (m *Model) Change(operate string, source string, args ...string) (int, bool, error) {
	if (operate == "") || (source == "") {
		return 0, false, errors.New("参数错误")
	}
	sqlStr := ""
	switch operate {
	case "add":
		if len(args) < 1 {
			return 0, false, errors.New("插入数据不能为空")
		}
		sqlStr = "INSERT INTO " + source + " VALUES (" + strings.Join(args, "),(") + ")"
	case "edit":
		if len(args) < 1 {
			return 0, false, errors.New("插入数据不能为空")
		}
		set, where := args[0], ""
		if len(args) > 1 {
			where = " WHERE " + args[1]
		}
		sqlStr = "UPDATE " + source + " SET " + set + where
	case "remove":
		where := ""
		if len(args) > 0 {
			where = " WHERE " + args[0]
		}
		sqlStr = "DELETE FROM " + source + where
	}
	if sqlStr == "" {
		return 0, false, errors.New("操作不允许")
	}
	result, err := m.db.Exec(sqlStr)
	if err != nil {
		return 0, false, err
	}
	if (operate == "add") && (len(args) == 1) {
		id, err := result.LastInsertId()
		if (err != nil) || (id < 1) {
			return 0, false, err
		}
		return int(id), true, nil
	}
	num, err := result.RowsAffected()
	if (err != nil) || (num < 1) {
		return 0, false, err
	}
	return int(num), true, nil
}

// 调用业务方法
func (s String) CService(args ...interface{}) []interface{} {

	return []interface{}{}
}
