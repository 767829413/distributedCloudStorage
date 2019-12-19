package handler

import (
	"distributedCloudStorage/common"
	"distributedCloudStorage/db/conn"
	"distributedCloudStorage/model"
	"distributedCloudStorage/util"
	"github.com/gin-gonic/gin"
	"time"

	//"distributedCloudStorage/common"
	//"distributedCloudStorage/db"
	//"distributedCloudStorage/util"
	"io/ioutil"
	"log"
	"net/http"
)

// User registration
func Signup(c *gin.Context) {
	var (
		data []byte
		name string
		pwd  string
		err  error
	)
	switch c.Request.Method {
	case http.MethodPost:
		_ = c.Request.ParseForm()
		name = c.Request.FormValue("username")
		pwd = c.Request.FormValue("password")
		enPwd := util.Sha1([]byte(pwd + common.UserPwdSalt))
		user := model.NewUser(name, enPwd)
		txn, _ := conn.GetDb().Begin()
		if flag := user.Save(txn); !flag {
			_ = txn.Rollback()
			c.Writer.WriteHeader(http.StatusInternalServerError)
			_, _ = c.Writer.Write([]byte("FAIL"))
			return
		}
		_ = txn.Commit()
		c.Writer.WriteHeader(http.StatusOK)
		_, _ = c.Writer.Write([]byte("SUCCESS"))
	case http.MethodGet:
		if data, err = ioutil.ReadFile(common.StaticFileDir + "/view/signup.html"); err != nil {
			log.Println("reade static file err : ", err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		c.Writer.WriteHeader(http.StatusOK)
		_, _ = c.Writer.Write(data)
	default:
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
	}
	return
}

// User login api
func SignIn(c *gin.Context) {
	var (
		//data []byte
		name string
		pwd  string
		err  error
	)
	_ = c.Request.ParseForm()
	name = c.Request.FormValue("username")
	pwd = c.Request.FormValue("password")
	enPwd := util.Sha1([]byte(pwd + common.UserPwdSalt))
	user := model.NewUser(name, enPwd)
	if err = user.Get(); err != nil {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	createAt := time.Now().Unix()
	if err = user.GenerateJwtToken(createAt); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	txn, _ := conn.GetDb().Begin()
	if flag := user.SaveToken(txn, createAt); !flag {
		_ = txn.Rollback()
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := util.NewRespMsg(0, "SUCCESS", struct {
		Token    string `json:"token"`
		UserName string `json:"user_name"`
		Location string `json:"Location"`
	}{
		Token:    user.Token,
		UserName: user.UserName,
		Location: "http://" + c.Request.Host + "/static/view/home.html",
	})
	_ = txn.Commit()
	_, _ = c.Writer.Write(resp.JSONBytes())
}

//Get user information
func Info(c *gin.Context) {
	var (
		err error
	)
	name := c.Request.Form.Get("username")
	user := model.NewUser(name, "")
	if err = user.GetUserInfo(); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := util.NewRespMsg(0, "SUCCESS", struct {
		UserName string `json:"username"`
		SignupAt string `json:"regtime"`
	}{
		UserName: user.UserName,
		SignupAt: user.SignupAt,
	})
	//w.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	_, _ = c.Writer.Write(resp.JSONBytes())
}
