package meta

/*
文件元数据处理:

*/

import (
	model "DStorage/model"
)

// FileMeta:文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta:新增/更新文件元信息
func UpdateFileMeta(fMeta FileMeta) {
	fileMetas[fMeta.FileSha1] = fMeta
}

// UpdateFileMetaDB:新增/更新文件元信息到mysql中
func UpdateFileMetaDB(fMeta FileMeta) bool {
	return model.FileUploadFinished(fMeta.FileSha1, fMeta.FileName, fMeta.FileSize, fMeta.Location)
}

// GetFileMetaDB:从mysql获取文件元信息
func GetFileMetaDB(fileSha1 string) (*FileMeta, error) {
	tFile, err := model.GetFileMeta(fileSha1)
	if tFile == nil || err != nil {
		return nil, err
	}
	fMeta := FileMeta{
		FileSha1: tFile.FileHash,
		FileName: tFile.FileName.String,
		FileSize: tFile.FileSize.Int64,
		Location: tFile.FileAddr.String,
	}
	return &fMeta, nil
}

// GetFileMeta:通过Sha1获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// RemoveFileMeta : 删除元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
