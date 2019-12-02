package model

import (
	"distributedCloudStorage/db"
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	UserName string `json:"user_name"`
	UserPwd  string `json:"user_pwd"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Token    string `json:"token"`
}

func NewUser(name string, pwd string) *User {
	return &User{
		UserName: name,
		UserPwd:  pwd,
	}
}

func (user *User) Save() bool {
	userDb := db.NewUser()
	userDb.UserName = user.UserName
	userDb.UserPwd = user.UserPwd
	return userDb.Save()
}

func (user *User) Get() (err error) {
	userDb := db.NewUser()
	if err = userDb.Get(user.UserName, user.UserPwd); err != nil {
		return
	}
	user.Phone = userDb.Phone
	user.Email = userDb.Email
	return
}

//Generate JWT Token
func (user *User) GenerateJwtToken(createAt int64) (err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":      user.UserName,
		"timestamp": createAt,
	})

	if user.Token, err = token.SigningString(); err != nil {
		return
	}
	return
}

//Save user token
func (user *User) SaveToken(createAt int64) bool {
	userToken := db.NewUserToken(user.UserName, user.Token, createAt)
	return userToken.Save()
}
