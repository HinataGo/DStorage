package service

import (
	"database/sql"
	"fmt"

	db "DStorage/db/mysql"
)

// TableFile : 文件表结构体
type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// FileUploadFinished : 文件上传完成,保存meta
func FileUploadFinished(fileHash string, fileName string, fileSize int64, fileAddr string) bool {

	stmt, err := db.Conn().Prepare("insert ignore into storage_file" +
		"(`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values(?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failed to prepare statement,err:" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(fileHash, fileName, fileSize, fileAddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Println("File with hash been upload before", fileHash)
		}
		return true
	}
	return false
}

// GetFileMeta : 从mysql获取文件元信息
// fix method
func GetFileMeta(fileHash string) (*TableFile, error) {
	stmt, err := db.Conn().Prepare(
		"select file_sha1,file_addr,file_name,file_size from tbl_file " +
			"where file_sha1=? and status=1 limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	file := TableFile{}
	err = stmt.QueryRow(fileHash).Scan(
		&file.FileHash, &file.FileAddr, &file.FileName, &file.FileSize)
	if err != nil {
		if err == sql.ErrNoRows {
			// 查不到对应记录， 返回参数及错误均为nil
			return nil, nil
		}
		fmt.Println(err.Error())
		return nil, err
	}
	return &file, nil
}
