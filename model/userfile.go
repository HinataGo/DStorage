package model

import (
	"fmt"
	"time"

	db "DStorage/db/mysql"
)

// UserFile : 用户文件表结构体
type UserFile struct {
	UserName       string
	FileHash       string
	FileName       string
	FileSize       int64
	UploadDate     string
	LastUpdateDate string
}

// OnUserFileUploadFinished : 更新用户文件表
func OnUserFileUploadFinished(username, fileHash, fileName string, fileSize int64) bool {
	stmt, err := db.Conn().Prepare(
		"insert ignore into storage_user_file (`user_name`,`file_sha1`,`file_name`," +
			"`file_size`,`upload_date`) values (?,?,?,?,?)")
	if err != nil {
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, fileHash, fileName, fileSize, time.Now())
	if err != nil {
		return false
	}
	return true
}

// QueryUserFileMetas : 批量获取用户文件信息 username 用户名 limit查询条目限制(解释所谓的批量)
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := db.Conn().Prepare(
		"select file_sha1,file_name,file_size,upload_date," +
			"last_update_date from storage_user_file where user_name=? limit ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		return nil, err
	}

	var userFiles []UserFile
	for rows.Next() {
		uFile := UserFile{}
		err = rows.Scan(&uFile.FileHash, &uFile.FileName, &uFile.FileSize,
			&uFile.UploadDate, &uFile.LastUpdateDate)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		userFiles = append(userFiles, uFile)
	}
	return userFiles, nil
}
