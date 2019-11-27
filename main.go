package main

import (
	"distributedCloudStorage/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.Upload)
	http.HandleFunc("/file/upload/success", handler.UploadSuccess)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("start server fail: ", err.Error())
	}
}
