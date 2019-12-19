package db

import (
	"database/sql"
	"distributedCloudStorage/common"
	"distributedCloudStorage/db/conn"
	"log"
)

type File struct {
	FileSha1 string `sql:"file_sha1"`
	FileName string `sql:"file_name"`
	FileSize int64  `sql:"file_size"`
	FileAddr string `sql:"file_addr"`
	Status   int    `sql:"status"`
	CreateAt string `sql:"create_at"`
	UpdateAt string `sql:"update_at"`
}

const fileTable = "tbl_file"

func NewFile() *File {
	return &File{}
}

//file meta info to database
func (fileInfo *File) Save(txn *sql.Tx) bool {
	return conn.Exec(txn, "insert ignore into `"+fileTable+"` (`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values (?,?,?,?,?)", fileInfo.FileSha1, fileInfo.FileName, fileInfo.FileSize, fileInfo.FileAddr, common.FileStateAvailable)
}

func (fileInfo *File) Get(fileHash string) (err error) {
	var (
		row *sql.Row
	)
	if row, _, err = conn.Get(conn.QueryGet, "select `file_sha1`,`file_name`,`file_size`,`file_addr`,`status`,`create_at`,`update_at` from `"+fileTable+"` where `file_sha1` = ? and status = ? limit 1", fileHash, common.FileStateAvailable); err != nil {
		log.Println(fileTable, " ", err.Error())
		return
	}
	if err = row.Scan(&fileInfo.FileSha1, &fileInfo.FileName, &fileInfo.FileSize, &fileInfo.FileAddr, &fileInfo.Status, &fileInfo.CreateAt, &fileInfo.UpdateAt); err != nil {
		log.Println(fileTable, " ", err.Error())
		return
	}
	return
}

func (fileInfo *File) Update(txn *sql.Tx) bool {

	return conn.Exec(txn, "update `"+fileTable+"` set `file_name` = ?,`file_size` = ?,`file_addr` = ? where `file_sha1` = ? and status = ? limit 1", fileInfo.FileName, fileInfo.FileSize, fileInfo.FileAddr, fileInfo.FileSha1, common.FileStateAvailable)
}

func (fileInfo *File) Delete(txn *sql.Tx, fileHash string) bool {
	return conn.Exec(txn, "update `"+fileTable+"`  set status = ? where `file_sha1` = ? and status = ? limit 1", common.FileStateDeleted, fileHash, common.FileStateAvailable)
}
