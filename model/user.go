package model

import (
	"distributedCloudStorage/common"
	"distributedCloudStorage/db"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type User struct {
	UserName string `json:"user_name"`
	UserPwd  string `json:"user_pwd"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Token    string `json:"token"`
	CreateAt int64  `json:"create_at"`
	SignupAt string `json:"signup_at"`
	jwt.StandardClaims
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
	user.SignupAt = userDb.SignupAt
	return
}

//Generate JWT Token
func (user *User) GenerateJwtToken(createAt int64) (err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name": user.UserName,
		"create_at": createAt,
	})
	if user.Token, err = token.SignedString([]byte(common.SecretKey)); err != nil {
		log.Println(err)
		return
	}
	return
}

//Save user token
func (user *User) SaveToken(createAt int64) bool {
	userToken := db.NewUserToken(user.UserName, user.Token, createAt)
	return userToken.Save()
}

//Check token
func (user *User) CheckToken(name string, tokenString string) bool {
	var (
		err error
	)
	if _, err = jwt.ParseWithClaims(tokenString, user, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(common.SecretKey), nil
	}); err != nil {
		log.Println("Parse token err: ", err.Error())
		return false
	}
	if user.UserName != name {
		return false
	}
	userToken := db.NewUserToken(user.UserName, tokenString, user.CreateAt)
	if err != userToken.Get() {
		return false
	}
	timeDiff := time.Now().Unix() - userToken.CreateAt
	if !(userToken.UserToken != user.Token || userToken.CreateAt != user.CreateAt || timeDiff > common.UserExpireTime) {
		return false
	}
	return true
}

func (user *User) GetUserInfo() (err error) {
	userDb := db.NewUser()
	if err = userDb.GetInfo(user.UserName); err != nil {
		return
	}
	user.Phone = userDb.Phone
	user.Email = userDb.Email
	user.SignupAt = userDb.SignupAt
	return
}
