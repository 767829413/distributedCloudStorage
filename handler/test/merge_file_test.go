package test

import (
	"distributedCloudStorage/handler"
	"log"
	"testing"
)

func TestMergeFile(t *testing.T) {
	uploadId := "admin15e73f2838b4bfd0"
	fileName := "C:/Users/NUC/Desktop/" + "Camera_Roll1.rar"
	log.Println("error : ", handler.MergeFile(uploadId, fileName))
}
