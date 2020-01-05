package main

import (
	"distributedCloudStorage/common"
	"distributedCloudStorage/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	// static file
	r.Static("/static", common.StaticFileDir)

	r.Any("/file/upload/success", handler.UploadSuccess)
	r.Any("/user/signup", handler.Signup)
	r.Any("/user/signin", handler.SignIn)

	r.POST("/file/mpupload/init", handler.InitBlockUpload)
	r.POST("/file/mpupload/uppart", handler.BlockUpload)
	r.POST("/file/mpupload/complete", handler.CompleteUpload)
	token := r.Group("/")
	token.Use(handler.Token)
	{
		//file relation
		token.Any("/file/upload", handler.Upload)
		token.POST("/file/meta", handler.GetMeta)
		token.POST("/file/query", handler.FileQuery)
		token.POST("/file/download", handler.DownLoad)
		token.POST("/file/update", handler.MetaUpdata)
		token.POST("/file/delete", handler.Delete)
		token.POST("/file/fastupload", handler.FastUpload)

		//token.POST("/file/mpupload/init", handler.InitBlockUpload)
		//token.OPTIONS("/file/mpupload/uppart", handler.BlockUpload)
		//token.OPTIONS("/file/mpupload/complete", handler.CompleteUpload)
		token.OPTIONS("/file/mpupload/cancel", handler.InitBlockUpload)
		token.OPTIONS("/file/mpupload/status", handler.InitBlockUpload)
		//user relation
		token.POST("/user/info", handler.Info)
	}
	_ = r.Run(":8080")
}
