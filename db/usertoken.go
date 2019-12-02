package db

import "distributedCloudStorage/db/conn"

type UserToke struct {
	UserName  string `sql:"user_name"`
	UserToken string `sql:"user_token"`
	CreateAt  int64  `sql:"create_at"`
}

const userTokeTable = "tbl_user_token"

func NewUserToken(name string, token string, expira int64) *UserToke {
	return &UserToke{
		UserName:  name,
		UserToken: token,
		CreateAt:  expira,
	}
}

//user name and password to user registration
func (userToken *UserToke) Save() bool {
	return conn.Exec("replace into `"+userTokeTable+"` (`user_name`,`user_token`,`create_at`) values (?,?,?)", userToken.UserName, userToken.UserToken, userToken.CreateAt)
}
