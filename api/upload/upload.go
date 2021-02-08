package upload

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	TStore "DStorage/config/path"
	cfg "DStorage/config/store"
	"DStorage/encrypt"
	"DStorage/meta"
	"DStorage/service"
	"DStorage/store/oss"
	"github.com/gin-gonic/gin"
)

// GetUploadHandler : 处理用户注册请求
func GetUploadHandler(c *gin.Context) {
	data, err := ioutil.ReadFile("./static/view/upload.html")
	if err != nil {
		c.String(404, `网页不存在`)
		return
	}
	c.Header("Content-TypeStore", "text/html; charset=utf-8")
	c.String(200, string(data))
}

// PostUploadHandler ： 处理文件上传
func PostUploadHandler(c *gin.Context) {
	errCode := 0
	defer func() {
		if errCode < 0 {
			c.JSON(http.StatusOK, gin.H{
				"exception": errCode,
				"msg":       "Upload failed",
			})
		}
	}()

	// 接收文件流及存储到本地目录
	file, head, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Printf("Failed to get data, err:%s\n", err.Error())
		errCode = -1
		return
	}
	defer file.Close()

	fileMeta := meta.FileMeta{
		FileName:   head.Filename,
		Location:   "/tmp/" + head.Filename,
		UploadDate: time.Now().Format("2006-01-02 15:04:05"),
	}

	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		fmt.Printf("Failed to create file, err:%s\n", err.Error())
		errCode = -2
		return
	}
	defer newFile.Close()

	fileMeta.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		fmt.Printf("Failed to save data into file, err:%s\n", err.Error())
		errCode = -3
		return
	}

	newFile.Seek(0, 0)
	fileMeta.FileSha1 = encrypt.FileSha1(newFile)

	// 游标重新回到文件头部
	newFile.Seek(0, 0)

	// 文件写入OSS存储
	if TStore.DefaultStoreType == cfg.OSSStore {
		ossPath := "oss/" + fileMeta.FileSha1
		err = oss.Bucket().PutObject(ossPath, newFile)
		if err != nil {
			fmt.Println(err.Error())
			errCode = -4
			return
		}
		fileMeta.Location = ossPath
	}
	// meta.UpdateFileMeta(fileMeta)
	_ = meta.UpdateFileMetaDB(fileMeta)
	// 更新用户文件表记录
	username := c.Request.FormValue("username")
	suc := service.UserFileUploadFinished(username, fileMeta.FileSha1,
		fileMeta.FileName, fileMeta.FileSize)
	if suc {
		c.Redirect(http.StatusFound, "/static/view/home.html")
	} else {
		errCode = -5
	}
}

// SucUploadHandler : 上传已完成
func SucUploadHandler(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"exception": 0,
			"msg":       "Upload Finish!",
		})
}

// GetFileMetaHandler : 获取文件元信息
func GetFileMetaHandler(c *gin.Context) {

	filehash := c.Request.FormValue("filehash")
	// fMeta := meta.GetFileMeta(filehash)
	fMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"exception": -2,
				"msg":       "Upload failed!",
			})
		return
	}

	if fMeta != nil {
		data, err := json.Marshal(fMeta)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{
					"exception": -3,
					"msg":       "Upload failed!",
				})
			return
		}
		c.Data(http.StatusOK, "application/json", data)
	} else {
		c.JSON(http.StatusOK,
			gin.H{
				"exception": -4,
				"msg":       "No sub file",
			})
	}
}

// FileQueryHandler : 查询批量的文件元信息
func FileQueryHandler(c *gin.Context) {

	limitCnt, _ := strconv.Atoi(c.Request.FormValue("limit"))
	username := c.Request.FormValue("username")
	// fileMetas, _ := meta.GetLastFileMetasDB(limitCnt)
	userFiles, err := service.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"exception": -1,
				"msg":       "Query failed!",
			})
		return
	}

	data, err := json.Marshal(userFiles)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"exception": -2,
				"msg":       "Query failed!",
			})
		return
	}
	c.Data(http.StatusOK, "application/json", data)
}

// DownloadHandler : 文件下载接口
func DownloadHandler(c *gin.Context) {
	fsha1 := c.Request.FormValue("filehash")
	fm, _ := meta.GetFileMetaDB(fsha1)

	c.FileAttachment(fm.Location, fm.FileName)
}

// FileMetaUpdateHandler ： 更新元信息接口(重命名)
func FileMetaUpdateHandler(c *gin.Context) {

	opType := c.Request.FormValue("op")
	fileSha1 := c.Request.FormValue("filehash")
	newFileName := c.Request.FormValue("filename")

	if opType != "0" {
		c.Status(http.StatusForbidden)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	// TODO: 更新文件表中的元信息记录

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, data)
}

// FileDeleteHandler : 删除文件及元信息
func FileDeleteHandler(c *gin.Context) {
	fileSha1 := c.Request.FormValue("filehash")

	fMeta := meta.GetFileMeta(fileSha1)
	// 删除文件
	os.Remove(fMeta.Location)
	// 删除文件元信息
	meta.RemoveFileMeta(fileSha1)
	// TODO: 删除表文件信息

	c.Status(http.StatusOK)
}

// DownloadURLHandler : 生成文件的下载地址
func DownloadURLHandler(c *gin.Context) {
	filehash := c.Request.FormValue("filehash")
	// 从文件表查找记录
	row, _ := service.GetFileMeta(filehash)

	// TODO: 判断文件存在OSS，还是在本地
	if strings.HasPrefix(row.FileAddr.String, "/tmp") {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")
		tmpUrl := fmt.Sprintf("http://%s/file/download?filehash=%s&username=%s&token=%s",
			c.Request.Host, filehash, username, token)
		c.Data(http.StatusOK, "octet-stream", []byte(tmpUrl))
	} else if strings.HasPrefix(row.FileAddr.String, "oss/") {
		// TODO:  oss下载url
		signedURL := oss.DownloadURL(row.FileAddr.String)
		c.Data(http.StatusOK, "octet-stream", []byte(signedURL))
	}
}
