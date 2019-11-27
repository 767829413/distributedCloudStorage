package handler

import (
	"distributedCloudStorage/meta"
	"distributedCloudStorage/util"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	var (
		fileMeta   *meta.FileMeta
		file       multipart.File
		fileHeader *multipart.FileHeader
		osFile     *os.File
		data       []byte
		err        error
	)
	switch r.Method {
	case http.MethodPost: //接收文件流并存储到本地目录
		if file, fileHeader, err = r.FormFile("file"); err != nil {
			log.Println("get file fail : ", err.Error())
			return
		}
		defer file.Close()
		fileMeta = &meta.FileMeta{
			FileName: fileHeader.Filename,
			Location: "C:/Users/gogo/src/distributedCloudStorage/tmp/" + fileHeader.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		if osFile, err = os.Create(fileMeta.Location); err != nil {
			log.Println("create file fail : ", err.Error())
			return
		}
		defer osFile.Close()
		if fileMeta.FileSize, err = io.Copy(osFile, file); err != nil {
			log.Println("save file fail : ", err.Error())
			return
		}
		_, _ = osFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(osFile)
		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	case http.MethodGet: //返回上传html页面
		if data, err = ioutil.ReadFile("./static/view/index.html"); err != nil {
			log.Println("reade static file err : ", err.Error())
			_, _ = io.WriteString(w, "internel server error")
		}
		_, _ = io.WriteString(w, string(data))
	}
}

func UploadSuccess(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "Upload File Success")
}
