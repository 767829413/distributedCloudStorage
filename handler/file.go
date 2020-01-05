package handler

import (
	"distributedCloudStorage/common"
	"distributedCloudStorage/db"
	"distributedCloudStorage/db/conn"
	"distributedCloudStorage/model"
	"distributedCloudStorage/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

//file upload
func Upload(c *gin.Context) {
	var (
		fileMeta   *model.File
		file       multipart.File
		fileHeader *multipart.FileHeader
		osFile     *os.File
		data       []byte

		err error
	)
	switch c.Request.Method {
	case http.MethodPost: //接收文件流并存储到本地目录
		if file, fileHeader, err = c.Request.FormFile("file"); err != nil {
			log.Println("get file fail : ", err.Error())
			return
		}
		defer file.Close()
		name := c.Request.Form.Get("username")

		fileMeta = model.NewFile()
		fileMeta.FileName = fileHeader.Filename
		fileMeta.Location = common.FileStoreTmp + fileHeader.Filename
		fileMeta.UploadAt = time.Now().Format("2006-01-02 15:04:05")

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
		txn, _ := conn.GetDb().Begin()
		flag := fileMeta.Save(txn)
		flagUser := fileMeta.SaveUserFile(txn, name)
		if !flag || !flagUser {
			_ = txn.Rollback()
		}
		_ = txn.Commit()
		http.Redirect(c.Writer, c.Request, "/file/upload/success", http.StatusFound)
	case http.MethodGet: //返回上传html页面
		if data, err = ioutil.ReadFile(common.StaticFileDir + "/view/index.html"); err != nil {
			log.Println("reade static file err : ", err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = c.Writer.Write(data)
	default:
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

//Fast upload file
func FastUpload(c *gin.Context) {
	var (
		err error
	)

	c.Writer.Header().Set("Content-Type", "application/json")
	name := c.Request.Form.Get("username")
	filehash := c.Request.Form.Get("filehash")
	fileName := c.Request.Form.Get("filename")
	fileSize, err := strconv.Atoi(c.Request.Form.Get("filesize"))
	fileMeta := model.NewFile()
	fileMeta.FileSha1 = filehash
	fileMeta.FileName = fileName
	fileMeta.FileSize = int64(fileSize)
	if err = fileMeta.Get(filehash); err != nil {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败,访问普通上传接口",
		}
		c.Writer.WriteHeader(http.StatusNotFound)
		_, _ = c.Writer.Write(resp.JSONBytes())
		return
	}
	txn, _ := conn.GetDb().Begin()
	if flag := fileMeta.SaveUserFile(txn, name); !flag {
		_ = txn.Rollback()
		resp := util.RespMsg{
			Code: -2,
			Msg:  "秒传失败",
		}
		c.Writer.WriteHeader(http.StatusOK)
		_, _ = c.Writer.Write(resp.JSONBytes())
		return
	}
	resp := util.RespMsg{
		Code: 0,
		Msg:  "秒传成功",
	}
	_ = txn.Commit()
	c.Writer.WriteHeader(http.StatusOK)
	_, _ = c.Writer.Write(resp.JSONBytes())
}

//file upload success
func UploadSuccess(c *gin.Context) {
	_, _ = io.WriteString(c.Writer, "Upload File Success")
}

//get file meta info
func GetMeta(c *gin.Context) {
	var (
		fileMeta *model.File
		data     []byte
		err      error
	)
	_ = c.Request.ParseForm()
	filehash := c.Request.FormValue("filehash")

	fileMeta = model.NewFile()
	if err = fileMeta.Get(filehash); err != nil {
		log.Println("get file meta err: ", err.Error())
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	if data, err = json.Marshal(fileMeta); err != nil {
		log.Println("Marshal File fail : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	_, _ = c.Writer.Write(data)
}

//Get user file info
func FileQuery(c *gin.Context) {
	var (
		userFiles []*db.UserFile
		data      []byte
		err       error
	)
	limit, _ := strconv.Atoi(c.Request.Form.Get("limit"))
	name := c.Request.Form.Get("username")
	user := model.NewUser(name, "")
	if userFiles, err = user.GetUserFiles(0, limit); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if data, err = json.Marshal(userFiles); err != nil {
		log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	_, _ = c.Writer.Write(data)
}

//file download
func DownLoad(c *gin.Context) {
	var (
		fileMeta *model.File
		file     *os.File
		data     []byte
		err      error
	)
	_ = c.Request.ParseForm()
	filehash := c.Request.Form.Get("filehash")
	fileMeta = model.NewFile()
	if err = fileMeta.Get(filehash); err != nil {
		log.Println("get file meta err: ", err.Error())
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	if file, err = os.Open(fileMeta.Location); err != nil {
		c.Writer.WriteHeader(http.StatusNotFound)
	}
	defer file.Close()

	if data, err = ioutil.ReadAll(file); err != nil {
		log.Println("read file fail :", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}
	c.Writer.Header().Set("Content-Type", "application/octect-stream")
	c.Writer.Header().Set("Content-disposition", "attachment;filename=\""+fileMeta.FileName+"\"")
	_, _ = c.Writer.Write(data)
}

//update file meta info
func MetaUpdata(c *gin.Context) {
	var (
		data     []byte
		err      error
		fileMeta *model.File
	)
	if c.Request.Method != http.MethodPost {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	_ = c.Request.ParseForm()
	opType := c.Request.PostFormValue("op")
	filehash := c.Request.PostFormValue("filehash")
	newFileName := c.Request.PostFormValue("filename")
	fileMeta = model.NewFile()
	if err = fileMeta.Get(filehash); err != nil {
		log.Println("get file meta err: ", err.Error())
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	switch opType {
	case "0":
		fileMeta.FileName = newFileName
		txn, _ := conn.GetDb().Begin()
		if flag := fileMeta.Update(txn); !flag {
			_ = txn.Rollback()
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_ = txn.Commit()
	default:
		c.Writer.WriteHeader(http.StatusForbidden)
		return
	}
	if data, err = json.Marshal(fileMeta); err != nil {
		log.Println("Marshal File fail: ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	_, _ = c.Writer.Write(data)
}

//delete file and file meta info
func Delete(c *gin.Context) {
	var (
		fileMeta *model.File
		flag     bool
		err      error
	)
	_ = c.Request.ParseForm()
	filehash := c.Request.PostFormValue("filehash")
	//hard delete
	fileMeta = model.NewFile()
	if err = fileMeta.Get(filehash); err != nil {
		log.Println("get file meta err: ", err.Error())
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	_ = os.Remove(fileMeta.Location)
	//soft delete
	txn, _ := conn.GetDb().Begin()
	if flag = fileMeta.Delete(txn, filehash); !flag {
		_ = txn.Rollback()
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = txn.Commit()
	c.Writer.WriteHeader(http.StatusOK)
	return
}
