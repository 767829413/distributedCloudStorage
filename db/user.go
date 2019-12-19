package db

import (
	"database/sql"
	"distributedCloudStorage/db/conn"
	"log"
)

type User struct {
	UserName string `sql:"user_name"`
	UserPwd  string `sql:"user_pwd"`
	Email    string `sql:"email"`
	Phone    string `sql:"phone"`
	SignupAt string `sql:"signup_at"`
	Status   string `sql:"status"`
}

const userTable = "tbl_user"

func NewUser() *User {
	return &User{}
}

//user name and password to user registration
func (user *User) Save(txn *sql.Tx) bool {
	return conn.Exec(txn, "insert ignore into `"+userTable+"` (`user_name`,`user_pwd`) values (?,?)", user.UserName, user.UserPwd)
}

//查询
func (user *User) Get(name string, pwd string) (err error) {
	var (
		row *sql.Row
	)
	if row, _, err = conn.Get(conn.QueryGet, "select `user_name`,`user_pwd`,`email`,`phone`,`signup_at`,`status` from `"+userTable+"` where `user_name` = ? and `user_pwd` = ? limit 1", name, pwd); err != nil {
		log.Println(err.Error())
		return
	}
	if err = row.Scan(&user.UserName, &user.UserPwd, &user.Email, &user.Phone, &user.SignupAt, &user.Status); err != nil {
		log.Println(err.Error())
		return
	}
	return
}

//查询信息
func (user *User) GetInfo(name string) (err error) {
	var (
		row *sql.Row
	)
	if row, _, err = conn.Get(conn.QueryGet, "select `user_name`,`user_pwd`,`email`,`phone`,`signup_at`,`status` from `"+userTable+"` where `user_name` = ? limit 1", name); err != nil {
		log.Println(userTable, " ", err.Error())
		return
	}
	if err = row.Scan(&user.UserName, &user.UserPwd, &user.Email, &user.Phone, &user.SignupAt, &user.Status); err != nil {
		log.Println(userTable, " ", err.Error())
		return
	}
	return
}
