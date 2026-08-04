package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"bluebell/logic"
	"bluebell/models"
)

// LoginResponse 登录成功返回的响应体
type LoginResponse struct {
	Token string `json:"token"`
}

// SignUpHandler 处理注册请求
// @Summary 用户注册
// @Description 用户注册接口，创建新账号
// @Tags 用户相关接口
// @Accept json
// @Produce json
// @Param object body models.ParamSignUp true "注册参数"
// @Success 200 {object} controller.Response
// @Router /signup [post]
func SignUpHandler(c *gin.Context) {
	// 1.获取参数和参数校验
	var p models.ParamSignUp //param是值类型后续传指针
	if err := c.ShouldBindJSON(&p); err != nil {
		// 输出日志
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//翻译错误
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": Translate(err),
		// })
		ResponseErrorWithMsg(c, CodeInvalidParams, Translate(err))
		return
	}

	// 手动对参数进行校验
	// if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 {
	// 	zap.L().Error("请求参数有误")
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"msg": "请求参数有误",
	// 	})
	// 	return
	// }

	// 2.业务处理
	if err := logic.SignUp(&p); err != nil {
		zap.L().Error("SignUp failed", zap.Error(err))
		// 区分业务错误和系统错误
		if errors.Is(err, logic.ErrUserExist) {
			// c.JSON(http.StatusOK, gin.H{"msg": err.Error()})
			ResponseError(c, CodeUserExist)
		} else if errors.Is(err, logic.ErrPasswordMismatch) {
			// c.JSON(http.StatusOK, gin.H{"msg": err.Error()})
			ResponseError(c, CodePasswordMismatch)
		} else {
			// c.JSON(http.StatusOK, gin.H{"msg": "服务器内部错误"})
			ResponseError(c, CodeServerBusy)
		}
		return
	}

	// 3.返回响应
	// c.JSON(http.StatusOK, gin.H{
	// 	"msg": "注册成功",
	// })
	ResponseSuccess(c, nil)
}

// LoginHandler 处理登录请求
// @Summary 用户登录
// @Description 用户登录接口，登录成功后返回 JWT Token
// @Tags 用户相关接口
// @Accept json
// @Produce json
// @Param object body models.ParamLogin true "登录参数"
// @Success 200 {object} controller.Response{data=controller.LoginResponse}
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	//1.获取参数和参数校验
	p := new(models.ParamLogin) //new返回指针，后续logic就直接传p
	if err := c.ShouldBindJSON(p); err != nil {
		// 输出日志
		zap.L().Error("Login with invalid param", zap.Error(err))
		//翻译错误
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": Translate(err),
		// })
		ResponseErrorWithMsg(c, CodeInvalidParams, Translate(err))
		return
	}
	//2.业务处理
	token, err := logic.Login(p)
	if err != nil {
		zap.L().Error("Login failed", zap.Error(err))
		// 区分业务错误和系统错误
		if errors.Is(err, logic.ErrInvalidCredentials) {
			// c.JSON(http.StatusOK, gin.H{"msg": "用户名或密码错误"})
			ResponseError(c, CodeInvalidCredentials)
		} else {
			// c.JSON(http.StatusOK, gin.H{"msg": "服务器内部错误"})
			ResponseError(c, CodeServerBusy)
		}
		return
	}
	//3.返回响应
	// c.JSON(http.StatusOK, gin.H{
	// 	"msg": "登录成功",
	// })
	ResponseSuccess(c, gin.H{"token": token})

}
