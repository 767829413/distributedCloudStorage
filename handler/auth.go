package handler

import (
	"distributedCloudStorage/model"
	"github.com/dgrijalva/jwt-go/request"
	"net/http"
	"strings"
)

func Token(httpFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			tokenString string
			err         error
		)
		_ = r.ParseForm()
		tokenString = r.Form.Get("token")
		if tokenString == "" {
			if tokenString, err = request.AuthorizationHeaderExtractor.ExtractToken(r); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			tokenString = strings.Replace(tokenString, "Bearer ", "", -1)
		}
		name := r.Form.Get("username")
		user := model.NewUser("", "")
		if flag := user.CheckToken(name, tokenString); !flag {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		httpFunc(w, r)
	}
}
