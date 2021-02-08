package download

import (
	"fmt"
	"net/http"
	"strings"

	"DStorage/meta"
	"DStorage/service"
	"DStorage/store/oss"
	"github.com/gin-gonic/gin"
)

// DownloadHandler : 文件下载接口
func DownloadHandler(c *gin.Context) {
	fSha1 := c.Request.FormValue("filehash")
	fm, _ := meta.GetFileMetaDB(fSha1)

	c.FileAttachment(fm.Location, fm.FileName)
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
