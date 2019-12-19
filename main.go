package main

import (
	"distributedCloudStorage/common"
	"distributedCloudStorage/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	// static file
	//r.StaticFile("/signin.html", common.StaticFileDir+"/view/signin.html")
	//r.StaticFile("/home.html", common.StaticFileDir+"/view/home.html")
	r.Static("/static", common.StaticFileDir)

	r.Any("/file/upload/success", handler.UploadSuccess)
	r.Any("/user/signup", handler.Signup)
	r.Any("/user/signin", handler.SignIn)
	token := r.Group("/")
	token.Use(handler.Token)
	{
		token.Any("/file/upload", handler.Upload)
		token.POST("/file/meta", handler.GetMeta)
		token.POST("/file/query", handler.FileQuery)
		token.POST("/file/download", handler.DownLoad)
		token.POST("/file/update", handler.MetaUpdata)
		token.POST("/file/delete", handler.Delete)
		token.POST("/file/fastupload", handler.FastUpload)
		token.POST("/user/info", handler.Info)
	}
	_ = r.Run(":8080")
}
