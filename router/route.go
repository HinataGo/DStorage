package router

import (
	"DStorage/api/download"
	"DStorage/api/fast_upload"
	"DStorage/api/file"
	"DStorage/api/part_upload"
	"DStorage/api/upload"
	"DStorage/api/user"
	"DStorage/handler"
	"github.com/gin-gonic/gin"
)

// Router : 路由表配置
func Router() *gin.Engine {
	// gin framework, 包括Logger, Recovery
	router := gin.Default()

	// 处理静态资源
	router.Static("/static/", "./static")

	// 不需要经过验证就能访问的接口
	router.GET("/model_user/signup", user.SignUpHandlerGet)
	router.POST("/model_user/signup", user.SignUpHandlerPost)

	router.GET("/model_user/signin", user.SignInHandlerGet)
	router.POST("/model_user/signin", user.SignInHandlerPost)

	// 加入中间件，用于校验token的拦截器
	router.Use(handler.HTTPInterceptor())

	// Use之后的所有handler都会经过拦截器进行token校验

	// 文件CRUD接口
	router.GET("/file/upload", upload.GetUploadHandler)
	router.POST("/file/upload", upload.PostUploadHandler)
	router.GET("/file/upload/suc", upload.SucUploadHandler)
	// 查询文件接口
	router.POST("/file/query", file.FileQueryHandler)
	// 文件删除
	router.POST("/file/delete", file.FileDeleteHandler)
	// 元数据获取
	router.POST("/file/meta", file.GetFileMetaHandler)
	// 元数据更新
	router.POST("/file/update", file.FileMetaUpdateHandler)

	// 文件下载服务
	router.GET("/file/download", download.DownloadHandler)
	router.POST("/file/downloadurl", download.DownloadURLHandler)

	// 秒传接口
	router.POST("/file/fastupload", fast_upload.FastUploadHandler)

	// 分块上传接口
	router.POST("/file/mpupload/init", part_upload.InitialMultipartUploadHandler)
	router.POST("/file/mpupload/uppart", part_upload.PartHandler)
	router.POST("/file/mpupload/complete", part_upload.CompleteUploadHandler)

	// 用户相关接口
	router.POST("/user/info", user.InfoUserHandler)

	return router
}
