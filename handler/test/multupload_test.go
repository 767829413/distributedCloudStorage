package test

import (
	"bufio"
	"bytes"
	"fmt"
	jsonit "github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"testing"
)

func multipartUpload(filename string, targetURL string, chunkSize int) error {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	index := 0

	buf := make([]byte, chunkSize) //每次读取chunkSize大小的内容
	wg := sync.WaitGroup{}
	for {
		n, err := bfRd.Read(buf)
		fmt.Println("ttttt : ", n)
		if n <= 0 {
			break
		}
		index++
		wg.Add(1)
		bufCopied := make([]byte, 5*1048576)
		copy(bufCopied, buf)

		go func(b []byte, curIdx int) {
			fmt.Printf("upload_size: %d\n", len(b))

			resp, err := http.Post(
				targetURL+"&index="+strconv.Itoa(curIdx),
				"multipart/form-data",
				bytes.NewReader(b))
			if err != nil {
				fmt.Println(err)
			}

			body, er := ioutil.ReadAll(resp.Body)
			fmt.Printf("body : %+v error : %+v\n", string(body), er)
			defer resp.Body.Close()
			defer wg.Done()
		}(bufCopied[:n], index)

		//遇到任何错误立即返回，并忽略 EOF 错误信息
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err.Error())
			}
		}
	}
	wg.Wait()
	return nil
}

func TestBlockUpload(t *testing.T) {
	username := "admin"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVfYXQiOjE1NzgyODgzMDcsInVzZXJfbmFtZSI6ImFkbWluIn0.lcZp5GnM2twFU4Wpv1fdg6NnrtXbzH2l19quFfpaS34"
	filehash := "1ba51899128dd9aa87a0e28d9a2d4a3b7595333a"

	// 1. 请求初始化分块上传接口
	resp, err := http.PostForm(
		"http://localhost:8080/file/mpupload/init",
		url.Values{
			"username": {username},
			"token":    {token},
			"filehash": {filehash},
			"filesize": {"42623945"},
		})
	if err != nil {
		t.Error(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err.Error())
	}
	uploadID := jsonit.Get(body, "data").Get("UploadId").ToString()
	chunkSize := jsonit.Get(body, "data").Get("ChunkSize").ToInt()
	t.Logf("uploadid: %s  chunksize: %d\n", uploadID, chunkSize)

	// 3. 请求分块上传接口
	filename := "C:/Users/NUC/Desktop/Camera_Roll.rar"
	tURL := "http://localhost:8080/file/mpupload/uppart?" +
		"username=admin&token=" + token + "&uploadid=" + uploadID
	multipartUpload(filename, tURL, chunkSize)
	// 4. 请求分块完成接口
	resp, err = http.PostForm(
		"http://localhost:8080/file/mpupload/complete",
		url.Values{
			"username": {username},
			"token":    {token},
			"filehash": {filehash},
			"filesize": {"42623945"},
			"filename": {"Camera_Roll.rar"},
			"uploadid": {uploadID},
		})

	if err != nil {
		t.Fatal(err.Error())
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("complete result: %s\n", string(body))
}


