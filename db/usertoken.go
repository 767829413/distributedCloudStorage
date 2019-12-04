package db

import (
	"database/sql"
	"distributedCloudStorage/db/conn"
	"log"
)

type UserToke struct {
	UserName  string `sql:"user_name"`
	UserToken string `sql:"user_token"`
	CreateAt  int64  `sql:"create_at"`
}

const userTokeTable = "tbl_user_token"

func NewUserToken(name string, token string, createAt int64) *UserToke {
	return &UserToke{
		UserName:  name,
		UserToken: token,
		CreateAt:  createAt,
	}
}

//user name and password to user registration
func (userToken *UserToke) Save() bool {
	return conn.Exec("replace into `"+userTokeTable+"` (`user_name`,`user_token`,`create_at`) values (?,?,?)", userToken.UserName, userToken.UserToken, userToken.CreateAt)
}

//Get user token
func (userToken *UserToke) Get() (err error) {
	var (
		row *sql.Row
	)
	if row, err = conn.Get("select `user_token`,`create_at` from `"+userTokeTable+"` where `user_name` = ? limit 1", userToken.UserName); err != nil {
		log.Println(err.Error())
		return
	}
	if err = row.Scan(&userToken.UserToken, &userToken.CreateAt); err != nil {
		log.Println(err.Error())
		return
	}
	return
}
