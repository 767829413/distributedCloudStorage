package test

import (
	"distributedCloudStorage/common"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestMergeFile(t *testing.T) {
	uploadId := "admin15e738d106611838"
	fileName := "C:/Users/NUC/Desktop/" + "Camera_Roll1.rar"
	log.Println("error : ", MergeFile(uploadId, fileName))
}

func MergeFile(uploadId string, fileName string) (err error) {
	if fileHd, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm); err != nil {
		log.Fatal(err.Error())
	} else {
		fileDir := common.FileStoreTmp + uploadId + "/"
		filepath.Walk(fileDir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				fileData, err := ioutil.ReadFile(path)
				if err != nil {
					log.Fatal(err.Error())
				}
				fileHd.Write(fileData)
			}
			return err
		})
		defer fileHd.Close()
	}
	return
}
