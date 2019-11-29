package conn

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var (
	db *sql.DB
)

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

func Exec(query string, args ...interface{}) bool {
	var (
		err    error
		stmt   *sql.Stmt
		result sql.Result
	)
	txn, _ := db.Begin()
	if stmt, err = txn.Prepare(query); err != nil {
		log.Println(err.Error())
		return false
	}

	defer stmt.Close()
	if result, err = stmt.Exec(args...); err != nil {
		log.Println(err.Error())
		return false
	}
	if num, err := result.RowsAffected(); err != nil {
		return false
	} else {
		if num <= 0 {
			return false
		}
		_ = txn.Commit()
		return true
	}
}

func Get(query string, args ...interface{}) (row *sql.Row, err error) {
	var (
		stmt *sql.Stmt
	)
	if stmt, err = db.Prepare(query); err != nil {
		log.Println(err.Error())
		return
	}
	defer stmt.Close()
	row = stmt.QueryRow(args...)
	return
}
