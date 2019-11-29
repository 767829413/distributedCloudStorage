package handler

import (
	"distributedCloudStorage/meta"
	"distributedCloudStorage/util"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

//file upload
func Upload(w http.ResponseWriter, r *http.Request) {
	var (
		fileMeta   *meta.Meta
		file       multipart.File
		fileHeader *multipart.FileHeader
		osFile     *os.File
		data       []byte

		err error
	)
	switch r.Method {
	case http.MethodPost: //接收文件流并存储到本地目录
		if file, fileHeader, err = r.FormFile("file"); err != nil {
			log.Println("get file fail : ", err.Error())
			return
		}
		defer file.Close()
		fileMeta = &meta.Meta{
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
		fileMeta.AddInfoDb()
		http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	case http.MethodGet: //返回上传html页面
		if data, err = ioutil.ReadFile("./static/view/index.html"); err != nil {
			log.Println("reade static file err : ", err.Error())
			_, _ = io.WriteString(w, "internel server error")
		}
		_, _ = io.WriteString(w, string(data))
	}
}

//file upload success
func UploadSuccess(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "Upload File Success")
}

//get file meta info
func GetMeta(w http.ResponseWriter, r *http.Request) {
	var (
		fileMeta *meta.Meta
		data     []byte
		err      error
	)
	_ = r.ParseForm()
	filehash := r.FormValue("filehash")

	fileMeta = &meta.Meta{}
	if err = fileMeta.GetInfoDb(filehash); err != nil {
		log.Println("get file meta err: ", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if data, err = json.Marshal(fileMeta); err != nil {
		log.Println("Marshal Meta fail : ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

//file download
func DownLoad(w http.ResponseWriter, r *http.Request) {
	var (
		fileMeta *meta.Meta
		file     *os.File
		data     []byte
		err      error
	)
	_ = r.ParseForm()
	filehash := r.Form.Get("filehash")
	fileMeta = &meta.Meta{}
	if err = fileMeta.GetInfoDb(filehash); err != nil {
		log.Println("get file meta err: ", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if file, err = os.Open(fileMeta.Location); err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	defer file.Close()

	if data, err = ioutil.ReadAll(file); err != nil {
		log.Println("read file fail :", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-disposition", "attachment;filename=\""+fileMeta.FileName+"\"")
	_, _ = w.Write(data)
}

//update file meta info
func MetaUpdata(w http.ResponseWriter, r *http.Request) {
	var (
		data     []byte
		err      error
		fileMeta *meta.Meta
	)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	_ = r.ParseForm()
	opType := r.PostFormValue("op")
	filehash := r.PostFormValue("filehash")
	newFileName := r.PostFormValue("filename")
	fileMeta = &meta.Meta{}
	if err = fileMeta.GetInfoDb(filehash); err != nil {
		log.Println("get file meta err: ", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch opType {
	case "0":
		fileMeta.FileName = newFileName
		if flag := fileMeta.UpdateInfoDb(); !flag {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if data, err = json.Marshal(fileMeta); err != nil {
		log.Println("Marshal Meta fail: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

//delete file and file meta info
func Delete(w http.ResponseWriter, r *http.Request) {
	var (
		fileMeta *meta.Meta
		flag     bool
		err      error
	)
	_ = r.ParseForm()
	filehash := r.PostFormValue("filehash")
	//hard delete
	fileMeta = &meta.Meta{}
	if err = fileMeta.GetInfoDb(filehash); err != nil {
		log.Println("get file meta err: ", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_ = os.Remove(fileMeta.Location)
	//soft delete
	if flag = fileMeta.DeleteInfoDb(filehash); !flag {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
