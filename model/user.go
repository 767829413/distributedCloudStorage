package model

import "distributedCloudStorage/db"

type User struct {
	UserName string `json:"user_name"`
	UserPwd  string `json:"user_pwd"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

func NewUser() *User {
	return &User{}
}

func (user *User) Save() bool {
	userDb := db.NewUser()
	userDb.UserName = user.UserName
	userDb.UserPwd = user.UserPwd
	return userDb.Save()
}
