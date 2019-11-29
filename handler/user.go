package handler

import (
	"distributedCloudStorage/common"
	"distributedCloudStorage/model"
	"distributedCloudStorage/util"
	//"distributedCloudStorage/common"
	//"distributedCloudStorage/db"
	//"distributedCloudStorage/util"
	"io/ioutil"
	"log"
	"net/http"
)

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
		name = r.FormValue("name")
		pwd = r.FormValue("pwd")
		enPwd := util.Sha1([]byte(pwd + common.UserPwdSalt))
		user := model.NewUser()
		user.UserName = name
		user.UserPwd = enPwd
		if flag := user.Save(); !flag {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
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
}
