package handler

import (
	cacheConn "distributedCloudStorage/cache/conn"
	"distributedCloudStorage/common"
	dbConn "distributedCloudStorage/db/conn"
	"distributedCloudStorage/model"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

//init info struct
type MultUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadId   string
	ChunkSize  int
	ChunkCount int
}

var chunkSize = 5 * 1024 * 1024 //5MB
var err error
// init Block upload
func InitBlockUpload(c *gin.Context) {
	var (
		username string
		filehash string
		filesize int
	)
	_ = c.Request.ParseForm()
	username = c.Request.Form.Get("username")
	filehash = c.Request.Form.Get("filehash")
	if filesize, err = strconv.Atoi(c.Request.Form.Get("filesize")); err != nil {
		c.JSON(400, gin.H{
			"code":    -1,
			"message": "params filesize invalid",
			"data":    "",
		})
	}
	rConn := cacheConn.GetPool().Get()
	defer rConn.Close()
	chunkCount := float64(filesize) / float64(chunkSize)
	chunkCount = math.Ceil(chunkCount)
	multInfo := MultUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadId:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  chunkSize,
		ChunkCount: int(chunkCount),
	}
	if _, err = rConn.Do("HMSET", "MP_"+multInfo.UploadId, "chunkcount", multInfo.ChunkCount, "filehash", multInfo.FileHash, "filesize", multInfo.FileSize); err != nil {
		log.Println("redis do error : ", err.Error())
		c.JSON(500, gin.H{
			"code":    -1,
			"message": "Internal system error",
			"data":    "",
		})
	}
	c.JSON(200, gin.H{
		"code":    0,
		"message": "OK",
		"data":    multInfo,
	})
}

//mult upload
func BlockUpload(c *gin.Context) {
	var (
		upid string
		//username   string
		chunkIndex string
		file       *os.File
		n          int
		err        error
	)
	_ = c.Request.ParseForm()
	upid = c.Request.Form.Get("uploadid")
	//username = c.Request.Form.Get("username")
	chunkIndex = c.Request.Form.Get("index")
	rConn := cacheConn.GetPool().Get()
	defer rConn.Close()
	filePath := common.FileStoreTmp + upid + "/" + chunkIndex
	err_dir := os.MkdirAll(path.Dir(filePath), 0744)
	if file, err = os.Create(filePath); err_dir != nil || err != nil {
		log.Println("os create file error : ", err_dir.Error(), "", err.Error())
		c.JSON(500, gin.H{
			"code":    -1,
			"message": "Internal system error",
			"data":    "",
		})
	}
	defer file.Close()
	buf := make([]byte, 1024*1024)
	for {
		n, err = c.Request.Body.Read(buf)
		_, _ = file.Write(buf[:n])
		if err != nil {
			break
		}
	}
	if _, err = rConn.Do("HSET", "MP_"+upid, "chkidx_"+chunkIndex, 1); err != nil {
		log.Println("redis do error : ", err.Error())
		c.JSON(500, gin.H{
			"code":    -1,
			"message": "Internal system error",
			"data":    "",
		})
	}
	c.JSON(200, gin.H{
		"code":    0,
		"message": "ok",
		"data":    "",
	})
}

//merge upload
func CompleteUpload(c *gin.Context) {
	var (
		upid     string
		username string
		filehash string
		filename string
		filesize int
		data     []interface{}
	)
	_ = c.Request.ParseForm()
	upid = c.Request.Form.Get("uploadid")
	username = c.Request.Form.Get("username")
	filehash = c.Request.Form.Get("filehash")
	if filesize, err = strconv.Atoi(c.Request.Form.Get("filesize")); err != nil {
		c.JSON(400, gin.H{
			"code":    -1,
			"message": "params filesize invalid",
			"data":    "",
		})
	}
	filename = c.Request.Form.Get("filename")
	rConn := cacheConn.GetPool().Get()
	defer rConn.Close()
	if data, err = redis.Values(rConn.Do("HGETALL", "MP_"+upid)); err != nil {
		log.Println("redis do error : ", err.Error())
		c.JSON(500, gin.H{
			"code":    -1,
			"message": "Internal system error",
			"data":    "",
		})
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		c.JSON(400, gin.H{
			"code":    -2,
			"message": "Invalid request",
			"data":    "",
		})
	}

	//TODO merge block upload file

	fileMeta := model.NewFile()
	fileMeta.FileSize = int64(filesize)
	fileMeta.FileSha1 = filehash
	fileMeta.FileName = filename
	fileMeta.Location = ""
	txn, _ := dbConn.GetDb().Begin()
	flag := fileMeta.Save(txn)
	flagUser := fileMeta.SaveUserFile(txn, username)
	if !flag || !flagUser {
		_ = txn.Rollback()
	}
	c.JSON(200, gin.H{
		"code":    0,
		"message": "OK",
		"data":    "",
	})
}
