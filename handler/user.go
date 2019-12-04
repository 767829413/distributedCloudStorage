package handler

import (
	"distributedCloudStorage/common"
	"distributedCloudStorage/db/conn"
	"distributedCloudStorage/model"
	"distributedCloudStorage/util"
	"time"

	//"distributedCloudStorage/common"
	//"distributedCloudStorage/db"
	//"distributedCloudStorage/util"
	"io/ioutil"
	"log"
	"net/http"
)

// User registration
func Signup(w http.ResponseWriter, r *http.Request) {
	var (
		data []byte
		name string
		pwd  string
		err  error
	)
	switch r.Method {
	case http.MethodPost:
		_ = r.ParseForm()
		name = r.FormValue("username")
		pwd = r.FormValue("password")
		enPwd := util.Sha1([]byte(pwd + common.UserPwdSalt))
		user := model.NewUser(name, enPwd)
		txn, _ := conn.GetDb().Begin()
		if flag := user.Save(txn); !flag {
			_ = txn.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("FAIL"))
			return
		}
		_ = txn.Commit()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("SUCCESS"))
	case http.MethodGet:
		if data, err = ioutil.ReadFile(common.StaticFileDir + "/view/signup.html"); err != nil {
			log.Println("reade static file err : ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	return
}

// User login api
func SignIn(w http.ResponseWriter, r *http.Request) {
	var (
		//data []byte
		name string
		pwd  string
		err  error
	)
	_ = r.ParseForm()
	name = r.FormValue("username")
	pwd = r.FormValue("password")
	enPwd := util.Sha1([]byte(pwd + common.UserPwdSalt))
	user := model.NewUser(name, enPwd)
	if err = user.Get(); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	createAt := time.Now().Unix()
	if err = user.GenerateJwtToken(createAt); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	txn, _ := conn.GetDb().Begin()
	if flag := user.SaveToken(txn, createAt); !flag {
		_ = txn.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := util.NewRespMsg(0, "SUCCESS", struct {
		Token    string `json:"token"`
		UserName string `json:"user_name"`
		Location string `json:"Location"`
	}{
		Token:    user.Token,
		UserName: user.UserName,
		Location: "http://" + r.Host + "/home.html",
	})
	_ = txn.Commit()
	_, _ = w.Write(resp.JSONBytes())
}

//Get user information
func Info(w http.ResponseWriter, r *http.Request) {
	var (
		err error
	)
	name := r.Form.Get("username")
	user := model.NewUser(name, "")
	if err = user.GetUserInfo(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp.JSONBytes())
}
