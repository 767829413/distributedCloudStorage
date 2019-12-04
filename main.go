package main

import (
	"distributedCloudStorage/common"
	"distributedCloudStorage/handler"
	"log"
	"net/http"
)

func main() {
	//file relation api
	http.HandleFunc("/file/upload", handler.Token(handler.Upload))
	http.HandleFunc("/file/upload/success", handler.UploadSuccess)
	http.HandleFunc("/file/meta", handler.GetMeta)
	http.HandleFunc("/file/query", handler.Token(handler.FileQuery))
	http.HandleFunc("/file/download", handler.DownLoad)
	http.HandleFunc("/file/update", handler.MetaUpdata)
	http.HandleFunc("/file/delete", handler.Delete)

	//user relation api
	http.HandleFunc("/user/signup", handler.Signup)
	http.HandleFunc("/user/signin", handler.SignIn)
	http.HandleFunc("/user/info", handler.Token(handler.Info))

	// static file
	http.Handle("/", http.FileServer(http.Dir(common.StaticFileDir)))
	http.Handle("/signin.html", http.FileServer(http.Dir(common.StaticFileDir+"/view")))
	http.Handle("/home.html", http.FileServer(http.Dir(common.StaticFileDir+"/view")))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("start server fail: ", err.Error())
	}
}
