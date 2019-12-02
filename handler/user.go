package handler

import (
	"distributedCloudStorage/common"
	"distributedCloudStorage/model"
	"distributedCloudStorage/util"
	"github.com/dgrijalva/jwt-go"
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
		if flag := user.Save(); !flag {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("FAIL"))
			return
		}
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
	if flag := user.SaveToken(createAt); !flag {
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
	_, _ = w.Write(resp.JSONBytes())
}

//Get user information
func Info(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	token := r.Form.Get("token")
	//name := r.Form.Get("username")
	parser := &jwt.Parser{}
	parser.ParseUnverified(token,jwt.MapClaims{})
}
