package model

import (
	"database/sql"
	"distributedCloudStorage/db"
)

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
func (fileMeta *File) Save(txn *sql.Tx) bool {
	fileInfo := fileMeta.getDbFile()
	return fileInfo.Save(txn)
}

// update meta info to db
func (fileMeta *File) Update(txn *sql.Tx) bool {
	fileInfo := fileMeta.getDbFile()
	return fileInfo.Update(txn)
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
func (fileMeta *File) Delete(txn *sql.Tx, fileSha1 string) bool {
	fileInfo := &db.File{}
	return fileInfo.Delete(txn, fileSha1)
}

func (fileMeta *File) getDbFile() (fileInfo *db.File) {
	fileInfo = db.NewFile()
	fileInfo.FileSha1 = fileMeta.FileSha1
	fileInfo.FileName = fileMeta.FileName
	fileInfo.FileAddr = fileMeta.Location
	fileInfo.FileSize = fileMeta.FileSize
	return
}

//Save user file
func (fileMeta *File) SaveUserFile(txn *sql.Tx, name string) bool {
	userFile := db.NewUserFile()
	userFile.UserName = name
	userFile.FileSha1 = fileMeta.FileSha1
	userFile.FileName = fileMeta.FileName
	userFile.FileSize = fileMeta.FileSize
	userFile.UploadAt = fileMeta.UploadAt
	return userFile.Save(txn)
}
