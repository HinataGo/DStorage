package file

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"DStorage/meta"
	"DStorage/service"
	"github.com/gin-gonic/gin"
)

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
