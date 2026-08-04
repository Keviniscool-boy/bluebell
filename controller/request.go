package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// 定义错误
var (
	ErrUserNotLogin = errors.New("用户未登录")
)

// getCurrentUserID 从上下文中获取当前请求用户的ID
func getCurrentUserID(c *gin.Context) (int64, error) {
	uid, ok := c.Get("userID")
	if !ok {
		return 0, ErrUserNotLogin
	}
	userID, ok := uid.(int64)
	if !ok {
		return 0, ErrUserNotLogin
	}
	return userID, nil
}
