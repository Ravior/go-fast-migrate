package util

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DbHelper = &dbHelper{
	dbs: make(map[string]*sql.DB, 5),
}

type dbHelper struct {
	dbs map[string]*sql.DB
}

func (d *dbHelper) GetConn(conn string) *sql.DB {
	host := ConfigHelper.GetString(fmt.Sprintf("db.%s.host", conn), "127.0.0.1")
	port := ConfigHelper.GetInt(fmt.Sprintf("db.%s.port", conn), 3306)
	username := ConfigHelper.GetString(fmt.Sprintf("db.%s.username", conn), "root")
	password := ConfigHelper.GetString(fmt.Sprintf("db.%s.password", conn), "")
	database := ConfigHelper.GetString(fmt.Sprintf("db.%s.database", conn), "migrate")

	connId := StrHelper.Md5(fmt.Sprintf("%s%d%s%s%s", host, port, username, password, database))
	db, ok := d.dbs[connId]

	if !ok {
		db = d.NewDb(host, port, username, password, database)
		d.dbs[connId] = db
	}

	return db
}

func (d *dbHelper) NewDb(host string, port int, username string, password, dbname string) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", username, password, host, port, dbname))
	if err != nil {
		SysHelper.Exit("DB Connect fail, err: %v", err)
	}

	return db
}

func (d *dbHelper) Exec(sql string, conns ...interface{}) (sql.Result, error) {
	conn := ToolHelper.GetStrParam(ArrHelper.Get(conns, 0), "default")
	LogHelper.Debug("SQL Execute: %s", sql)
	return d.GetConn(conn).Exec(sql)
}

func (d *dbHelper) Query(sql string, conns ...interface{}) (*sql.Rows, error) {
	conn := ToolHelper.GetStrParam(ArrHelper.Get(conns, 0), "default")
	LogHelper.Debug("SQL Execute: %s", sql)
	return d.GetConn(conn).Query(sql)
}

func (d *dbHelper) QueryRow(sql string, conns ...interface{}) *sql.Row {
	conn := ToolHelper.GetStrParam(ArrHelper.Get(conns, 0), "default")
	LogHelper.Debug("SQL Execute: %s", sql)
	return d.GetConn(conn).QueryRow(sql)
}
