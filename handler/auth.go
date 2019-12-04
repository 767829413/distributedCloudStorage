package handler

import (
	"distributedCloudStorage/model"
	"net/http"
)

func Token(httpFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		//if tokenString, err = request.AuthorizationHeaderExtractor.ExtractToken(r); err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	return
		//}
		//tokenString = strings.Replace(tokenString, "Bearer ", "", -1)
		name := r.Form.Get("username")
		tokenString := r.Form.Get("token")
		user := model.NewUser("", "")
		if flag := user.CheckToken(name, tokenString); !flag {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		httpFunc(w, r)
	}
}
