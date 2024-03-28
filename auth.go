package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pc-phone-go/conf"
	"pc-phone-go/entity"
	"pc-phone-go/funcs/logger"
	"strings"
)

const BearerSchema = "Bearer "

// UseAuth 验证请求的中间件
//
// 验证信息可以在请求头`Authorization`中，也可以在查询字符串`auth`中
//
// 中间件 https://www.alexedwards.net/blog/making-and-using-middleware
func UseAuth(c *gin.Context) {
	// 不是 /api/ 的请求，直接下一步
	if !strings.HasPrefix(c.FullPath(), "/api/") {
		c.Next()
		return
	}

	var token string
	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		if len(authHeader) < len(BearerSchema) {
			abortAuth(c, fmt.Sprintf("非法的授权码'%s'", authHeader))
			return
		}

		// Trim Bearer prefix to get the token
		token = authHeader[len(BearerSchema):]
	} else if authStr := c.Query("auth"); authStr != "" {
		token = authStr
	} else {
		abortAuth(c, "没有授权码")
		return
	}

	logger.Info.Printf("访问'%s'的完整信息：'%s'\n", c.FullPath(), token)

	// 验证通过，继续下一步
	if token == conf.Conf.Auth {
		logger.Info.Printf("已通过授权，继续下一步 '%s'\n", c.FullPath())
		c.Next()
		return
	}

	// 没有匹配到有效的验证码，禁止访问
	abortAuth(c, fmt.Sprintf("无效的授权码'%s'", token))
}

// 拒绝访问
func abortAuth(c *gin.Context, msg string) {
	logger.Warn.Printf("'%s' 拒绝访问: %s\n", c.FullPath(), msg)

	c.AbortWithStatusJSON(http.StatusOK, entity.Rest{Code: 10000, Msg: fmt.Sprintf("拒绝访问: %s", msg)})
}

// 生成随机、唯一的 Bearer Authorization
//
// n 表示长度。如 generateRandomToken(32)
func generateRandomToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
