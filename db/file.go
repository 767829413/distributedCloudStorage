package db

import (
	"database/sql"
	"distributedCloudStorage/db/conn"
	"log"
)

//file meta info to database
func OnFileUploadFinished(fileSha1 string, fileName string, fileSize int64, location string) bool {
	var (
		err    error
		stmt   *sql.Stmt
		result sql.Result
	)
	txn, _ := conn.MysqlConn().Begin()
	if stmt, err = txn.Prepare("insert ignore into `tbl_file` (`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values (?,?,?,?,1)"); err != nil {
		log.Println("Fail to prepare statement,err: ", err.Error())
		return false
	}
	//if stmt, err = db.Prepare("insert ignore into `tbl_file` (`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values (?,?,?,?,1)"); err != nil {
	//	log.Println("Fail to prepare statement,err: ", err.Error())
	//	return false
	//}
	defer stmt.Close()
	if result, err = stmt.Exec(fileSha1, fileName, fileSize, location); err != nil {
		log.Println("Exec insert data err: ", err.Error())
		return false
	}
	if num, err := result.RowsAffected(); err != nil {
		return false
	} else {
		if num <= 0 {
			log.Println("File repeated insertion, hash: ", fileSha1)
		}
		_ = txn.Commit()
		return true
	}
}
