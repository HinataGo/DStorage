package upload

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	TStore "DStorage/config/path"
	cfg "DStorage/config/store"
	"DStorage/encrypt"
	"DStorage/meta"
	"DStorage/service"
	"DStorage/store/oss"
	"github.com/gin-gonic/gin"
)

// GetUploadHandler : 处理用户上传请求
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
