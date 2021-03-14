package user

import (
	"fmt"
	"net/http"
	"time"

	"DStorage/format/json"
	"github.com/gin-gonic/gin"

	"DStorage/encrypt"
	"DStorage/service"
)

const (
	// 用于加密的盐值(自定义)
	pwdSalt = "*#890"
)

// SignUpHandlerGet : 处理用户注册请求
func SignUpHandlerGet(c *gin.Context) {
	c.Redirect(http.StatusFound, "http://"+c.Request.Host+"/static/view/signup.html")
}

// SignUpHandlerPost : 处理用户注册请求
func SignUpHandlerPost(c *gin.Context) {
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")

	if len(username) < 3 || len(passwd) < 5 {
		c.JSON(http.StatusOK,
			gin.H{
				"msg": "Invalid parameter",
			})
		return
	}

	// 对密码进行加盐及取Sha1值加密
	encPasswd := encrypt.Sha1([]byte(passwd + pwdSalt))
	// 将用户信息注册到用户表中
	suc := service.UserSignUp(username, encPasswd)
	if suc {
		c.JSON(http.StatusOK,
			gin.H{
				"exception": 0,
				"msg":       "注册成功",
				"data":      nil,
				"forward":   "/model_user/signin",
			})
	} else {
		c.JSON(http.StatusOK,
			gin.H{
				"exception": 0,
				"msg":       "注册失败",
				"data":      nil,
			})
	}
}

// SignInHandler : 处理用户注册请求
func SignInHandlerGet(c *gin.Context) {
	c.Redirect(http.StatusFound, "http://"+c.Request.Host+"/static/view/signin.html")
}

// DoSignInHandler : 登录接口
func SignInHandlerPost(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	encPasswd := encrypt.Sha1([]byte(password + pwdSalt))

	// 1. 校验用户名及密码
	pwdChecked := service.UserSignIn(username, encPasswd)
	if !pwdChecked {
		c.JSON(http.StatusOK,
			gin.H{
				"exception": 0,
				"msg":       "密码校验失败",
				"data":      nil,
			})
		return
	}

	// 2. 生成访问凭证(token)
	token := GenToken(username)
	upRes := service.UpdateToken(username, token)
	if !upRes {
		c.JSON(http.StatusOK,
			gin.H{
				"exception": 0,
				"msg":       "登录失败",
				"data":      nil,
			})
		return
	}

	// 3. 登录成功后重定向到首页
	resp := json.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + c.Request.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	c.Data(http.StatusOK, "octet-stream", resp.JSONBytes())
}

// InfoUserHandler ： 查询用户信息
func InfoUserHandler(c *gin.Context) {
	// 1. 解析请求参数
	username := c.Request.FormValue("username")
	token := c.Request.FormValue("token")

	// TODO 2. 验证token是否有效
	isValidToken := encrypt.IsTokenValid(token)
	if !isValidToken {
		c.JSON(http.StatusForbidden,
			gin.H{})
		return
	}

	// 3. 查询用户信息
	user, err := service.GetUserInfo(username)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{})
		return
	}

	// 4. 组装并且响应用户数据
	resp := json.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	c.Data(http.StatusOK, "octet-stream", resp.JSONBytes())
}

// GenToken : 生成token
func GenToken(username string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := encrypt.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// ExistsUserHandler ： 查询用户是否存在
func ExistsUserHandler(c *gin.Context) {
	// 1. 解析请求参数
	username := c.Request.FormValue("username")

	// 3. 查询用户信息
	exists, err := service.UserExist(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"msg": "server error",
			})
	} else {
		c.JSON(http.StatusOK,
			gin.H{
				"msg":    "ok",
				"exists": exists,
			})
	}
}
