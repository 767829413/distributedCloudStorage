package conn

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var db *sql.DB

func init() {
	var (
		err error
	)
	if db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3339)/file_server?charset=utf8mb4"); err != nil {
		log.Println("connect mysql fail: ", err.Error())
	}
	db.SetMaxOpenConns(1000)
	if err = db.Ping(); err != nil {
		log.Println("ping mysql fail: ", err.Error())
		os.Exit(1)
	}
}

func MysqlConn() *sql.DB {
	return db
}
