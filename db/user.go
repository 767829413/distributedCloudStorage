package db

import (
	"distributedCloudStorage/db/conn"
)

type User struct {
	UserName string `sql:"user_name"`
	UserPwd  string `sql:"user_pwd"`
	Email    int64  `sql:"email"`
	Phone    string `sql:"phone"`
	SignupAt int    `sql:"signup_at"`
	Status   string `sql:"status"`
}

func NewUser() *User {
	return &User{}
}

//user name and password to user registration
func (user *User) Save() bool {
	return conn.Exec("insert ignore into `tbl_user` (`user_name`,`user_pwd`) values (?,?)", user.UserName, user.UserPwd)
}
