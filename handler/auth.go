package handler

import (
	"distributedCloudStorage/model"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Token(c *gin.Context) {
	var (
		tokenString string
		err         error
	)
	_ = c.Request.ParseForm()
	tokenString = c.Request.Form.Get("token")
	if tokenString == "" {
		if tokenString, err = request.AuthorizationHeaderExtractor.ExtractToken(c.Request); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", -1)
	}
	name := c.Request.Form.Get("username")
	user := model.NewUser("", "")
	if flag := user.CheckToken(name, tokenString); !flag {
		c.Writer.WriteHeader(http.StatusForbidden)
		return
	}
}
