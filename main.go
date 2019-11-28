package main

import (
	"distributedCloudStorage/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.Upload)
	http.HandleFunc("/file/upload/success", handler.UploadSuccess)
	http.HandleFunc("/file/meta", handler.GetMeta)
	http.HandleFunc("/file/download", handler.DownLoad)
	http.HandleFunc("/file/update", handler.MetaUpdata)
	http.HandleFunc("/file/delete", handler.Delete)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("start server fail: ", err.Error())
	}
}
