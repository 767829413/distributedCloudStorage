package db

import (
	"database/sql"
	"distributedCloudStorage/common"
	"distributedCloudStorage/db/conn"
	"log"
)

type FileMetaInfo struct {
	FileSha1 string `sql:"file_sha1"`
	FileName string `sql:"file_name"`
	FileSize int64  `sql:"file_size"`
	FileAddr string `sql:"file_addr"`
	Status   int    `sql:"status"`
	CreateAt string `sql:"create_at"`
	UpdateAt string `sql:"update_at"`
}

//file meta info to database
func (fileInfo *FileMetaInfo) OnFileUploadFinished() bool {
	var (
		err    error
		stmt   *sql.Stmt
		result sql.Result
	)
	txn, _ := conn.MysqlConn().Begin()
	if stmt, err = txn.Prepare("insert ignore into `tbl_file` (`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values (?,?,?,?,?)"); err != nil {
		log.Println("Fail to prepare insert statement,err: ", err.Error())
		return false
	}

	defer stmt.Close()
	if result, err = stmt.Exec(fileInfo.FileSha1, fileInfo.FileName, fileInfo.FileSize, fileInfo.FileAddr, common.FileStateAvailable); err != nil {
		log.Println("Exec insert data err: ", err.Error())
		return false
	}
	if num, err := result.RowsAffected(); err != nil {
		return false
	} else {
		if num <= 0 {
			log.Println("File repeated insertion, hash: ", fileInfo.FileSha1)
		}
		_ = txn.Commit()
		return true
	}
}

func (fileInfo *FileMetaInfo) GetFileMetaInfo(fileHash string) (err error) {
	var (
		stmt *sql.Stmt
	)
	if stmt, err = conn.MysqlConn().Prepare("select `file_sha1`,`file_name`,`file_size`,`file_addr`,`status`,`create_at`,`update_at` from `tbl_file` where `file_sha1` = ? and status = ? limit 1"); err != nil {
		log.Println("Fail to prepare select statement,err: ", err.Error())
		return
	}
	defer stmt.Close()
	if err = stmt.QueryRow(fileHash, common.FileStateAvailable).Scan(&fileInfo.FileSha1, &fileInfo.FileName, &fileInfo.FileSize, &fileInfo.FileAddr, &fileInfo.Status, &fileInfo.CreateAt, &fileInfo.UpdateAt); err != nil {
		log.Println("Fail to get query row,err: ", err.Error())
		return
	}
	return
}

func (fileInfo *FileMetaInfo) UpdateFileMetaInfo() bool {
	var (
		err    error
		stmt   *sql.Stmt
		result sql.Result
	)
	txn, _ := conn.MysqlConn().Begin()
	if stmt, err = txn.Prepare("update `tbl_file` set `file_name` = ?,`file_size` = ?,`file_addr` = ? where `file_sha1` = ? and status = ? limit 1"); err != nil {
		log.Println("Fail to prepare update statement,err: ", err.Error())
		return false
	}

	defer stmt.Close()
	if result, err = stmt.Exec(fileInfo.FileName, fileInfo.FileSize, fileInfo.FileAddr, fileInfo.FileSha1, common.FileStateAvailable); err != nil {
		log.Println("Exec update data err: ", err.Error())
		return false
	}
	if num, err := result.RowsAffected(); err != nil {
		return false
	} else {
		if num <= 0 {
			log.Println("File repeated update, hash: ", fileInfo.FileSha1)
		}
		_ = txn.Commit()
		return true
	}
}

func (fileInfo *FileMetaInfo) DeleteFileMetaInfo(fileHash string) bool {
	var (
		err    error
		stmt   *sql.Stmt
		result sql.Result
	)
	txn, _ := conn.MysqlConn().Begin()
	if stmt, err = txn.Prepare("update `tbl_file`  set status = ? where `file_sha1` = ? and status = ? limit 1"); err != nil {
		log.Println("Fail to prepare soft delete statement,err: ", err.Error())
		return false
	}

	defer stmt.Close()
	if result, err = stmt.Exec(common.FileStateDeleted, fileHash, common.FileStateAvailable); err != nil {
		log.Println("Exec soft delete data err: ", err.Error())
		return false
	}
	if num, err := result.RowsAffected(); err != nil {
		return false
	} else {
		if num <= 0 {
			log.Println("File repeated soft delete, hash: ", fileHash)
			return false
		}
		_ = txn.Commit()
		return true
	}
}
