package gun

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type SQLCommon interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type Session struct {
	db         *sql.DB
	sql        strings.Builder
	find       []interface{}
}

type Conn struct {
	Username   string //连接名
	Password   string //密码
	Host       string //地址
	DBName     string //数据库名
	Port       int    //端口
}

type DbEngine struct {
	db   SQLCommon
}

func (c *Conn)DbConn(d *Conn) *DbEngine {
	if c.Username == "" || c.Host == "" || c.DBName == "" || c.Port == 0{
		log.Printf("连接名/地址/数据库名/端口 不能为空")
	}
	str := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", c.Username, c.Password, c.Host, c.Port, c.DBName)
	fmt.Println(str)
	dbSQL, err := sql.Open("mysql", str)
	//fmt.Println(db)
	if err != nil{
		log.Printf("数据库连接错误")
	}
	db := &DbEngine{
		db:  dbSQL,
	}
	return db
}

func (c *DbEngine)Select()  {
	
}
