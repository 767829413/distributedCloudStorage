package conn

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var (
	db *sql.DB
)

const (
	QueryGet  = 1
	QueryList = 2
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

func GetDb() *sql.DB {
	return db
}

func Exec(txn *sql.Tx, query string, args ...interface{}) bool {
	var (
		err    error
		stmt   *sql.Stmt
		result sql.Result
	)
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
		log.Println(err.Error())
		return false
	} else {
		if num <= 0 {
			log.Println(query)
			return false
		}
		return true
	}
}

func Get(queryType int, query string, args ...interface{}) (row *sql.Row, rows *sql.Rows, err error) {
	var (
		stmt *sql.Stmt
	)
	if stmt, err = db.Prepare(query); err != nil {
		log.Println(err.Error())
		return
	}
	defer stmt.Close()
	switch queryType {
	case QueryGet:
		row = stmt.QueryRow(args...)
		return
	case QueryList:
		if rows, err = stmt.Query(args...); err != nil {
			log.Println(err)
			return
		}
		return
	default:
		err = errors.New("choose the right query type")
		return
	}
}
