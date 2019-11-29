package meta

import "distributedCloudStorage/db"

//File Meta information struct
type Meta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

// add meta info to db
func (fileMeta *Meta) AddInfoDb() bool {
	fileInfo := &db.FileMetaInfo{
		FileSha1: fileMeta.FileSha1,
		FileName: fileMeta.FileName,
		FileAddr: fileMeta.Location,
		FileSize: fileMeta.FileSize,
	}
	return fileInfo.OnFileUploadFinished()
}

// update meta info to db
func (fileMeta *Meta) UpdateInfoDb() bool {
	fileInfo := &db.FileMetaInfo{
		FileSha1: fileMeta.FileSha1,
		FileName: fileMeta.FileName,
		FileAddr: fileMeta.Location,
		FileSize: fileMeta.FileSize,
	}
	return fileInfo.UpdateFileMetaInfo()
}

//get meta info for db
func (fileMeta *Meta) GetInfoDb(fileSha1 string) (err error) {
	fileInfo := &db.FileMetaInfo{}
	if err = fileInfo.GetFileMetaInfo(fileSha1); err != nil {
		return
	}
	fileMeta.FileSize = fileInfo.FileSize
	fileMeta.FileName = fileInfo.FileName
	fileMeta.FileSha1 = fileInfo.FileSha1
	fileMeta.Location = fileInfo.FileAddr
	fileMeta.UploadAt = fileInfo.UpdateAt
	return
}

//remove *Meta PS: Thread-safe operation, map is Non-thread safe
func (fileMeta *Meta) DeleteInfoDb(fileSha1 string) bool {
	fileInfo := &db.FileMetaInfo{}
	return fileInfo.DeleteFileMetaInfo(fileSha1)
}
