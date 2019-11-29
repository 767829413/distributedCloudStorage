package model

import "distributedCloudStorage/db"

//File File information struct
type File struct {
	FileSha1 string `json:"file_sha1"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	Location string `json:"location"`
	UploadAt string `json:"uploadat"`
}

func NewFile() *File {
	return &File{}
}

// add meta info to db
func (fileMeta *File) Add() bool {
	fileInfo := fileMeta.getDbFile()
	return fileInfo.Save()
}

// update meta info to db
func (fileMeta *File) Update() bool {
	fileInfo := fileMeta.getDbFile()
	return fileInfo.Update()
}

//get meta info for db
func (fileMeta *File) Get(fileSha1 string) (err error) {
	fileInfo := &db.File{}
	if err = fileInfo.Get(fileSha1); err != nil {
		return
	}
	fileMeta.FileSize = fileInfo.FileSize
	fileMeta.FileName = fileInfo.FileName
	fileMeta.FileSha1 = fileInfo.FileSha1
	fileMeta.Location = fileInfo.FileAddr
	fileMeta.UploadAt = fileInfo.UpdateAt
	return
}

//remove *File PS: Thread-safe operation, map is Non-thread safe
func (fileMeta *File) Delete(fileSha1 string) bool {
	fileInfo := &db.File{}
	return fileInfo.Delete(fileSha1)
}

func (fileMeta *File) getDbFile() (fileInfo *db.File) {
	fileInfo = db.NewFile()
	fileInfo.FileSha1 = fileMeta.FileSha1
	fileInfo.FileName = fileMeta.FileName
	fileInfo.FileAddr = fileMeta.Location
	fileInfo.FileSize = fileMeta.FileSize
	return
}
