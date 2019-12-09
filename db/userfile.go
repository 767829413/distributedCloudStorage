package db

import (
	"database/sql"
	"distributedCloudStorage/db/conn"
	"log"
)

type UserFile struct {
	FileName   string `sql:"file_name"`
	FileSha1   string `sql:"file_sha1"`
	UserName   string `sql:"user_name"`
	FileSize   int64  `sql:"file_size"`
	Status     int    `sql:"status"`
	LastUpdate string `sql:"last_update"`
	UploadAt   string `sql:"upload_at"`
}

const userFileTable = "tbl_user_file"

func NewUserFile() *UserFile {
	return &UserFile{}
}

//Update user file table
func (userFile *UserFile) Save(txn *sql.Tx) bool {
	return conn.Exec(txn, "insert ignore into `"+userFileTable+"` (`user_name`,`file_sha1`,`file_name`,`file_size`,`upload_at`) values (?,?,?,?,?)", userFile.UserName, userFile.FileSha1, userFile.FileName, userFile.FileSize, userFile.UploadAt)
}

//User file total
func (userFile *UserFile) Count() (num int, err error) {
	var (
		row *sql.Row
	)
	if row, _, err = conn.Get(conn.QueryGet, "select count(1) as num from `"+userTable+"` where `user_name` = ?", userFile.UserName); err != nil {
		log.Println(err.Error())
		return
	}
	if err = row.Scan(&num); err != nil {
		log.Println(err.Error())
		return
	}
	return
}

//User file list
func (userFile *UserFile) List(page int, limit int) (userList []*UserFile, err error) {
	var (
		rows *sql.Rows
	)
	if _, rows, err = conn.Get(conn.QueryList, "select `file_sha1`,`file_name`,`file_size`,`upload_at`,`last_update` from `"+userFileTable+"` where `user_name` = ? limit ?,?", userFile.UserName, page, limit); err != nil {
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		userFileItem := NewUserFile()
		if err = rows.Scan(&userFileItem.FileSha1, &userFileItem.FileName, &userFileItem.FileSize, &userFileItem.UploadAt, &userFileItem.LastUpdate); err != nil {
			log.Println(userFile.UserName, err.Error())
			return
		}
		userList = append(userList, userFileItem)
	}
	return
}
