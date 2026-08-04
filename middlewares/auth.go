package middlewares

import (
	"bluebell/controller"
	"bluebell/pkg/jwt"
	"github.com/gin-gonic/gin"
	"strings"
)

// jwt 中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带 Token 有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设 Token 放在 Header 的 Authorization 中，并使用 Bearer 开头
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}

		// 按空格分割，取 Bearer 后面的 token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			controller.ResponseError(c, controller.CodeInvalidAuth)
			c.Abort()
			return
		}
		tokenString := parts[1]

		// 解析 Token
		mc, err := jwt.ParseToken(tokenString)
		if err != nil {
			controller.ResponseError(c, controller.CodeInvalidAuth)
			c.Abort()
			return
		}

		// 将解析后的用户信息存储到上下文中
		c.Set("userID", mc.UserId)
		c.Set("userName", mc.UserName)

		c.Next()
	}
}
