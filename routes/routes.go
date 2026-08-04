package routes

import (
	"net/http"

	"bluebell/controller"
	_ "bluebell/docs"
	"bluebell/logger"
	"bluebell/middlewares"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 100)) // 限流中间件，2秒钟最多处理100个请求

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// 注册 Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := r.Group("api/v1")

	// 注册路由
	v1.POST("/signup", controller.SignUpHandler)
	// 登录路由
	v1.POST("/login", controller.LoginHandler)
	//jwt中间件
	v1.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		//如果是已经登录的用户，判断请求头中是否携带了token
		c.String(http.StatusOK, "PONG")
	})
	// 社区相关路由
	v1.GET("/community", controller.CommunityHandler)
	v1.GET("/community/:id", controller.CommunityDetailHandler)

	// 帖子相关路由
	v1.POST("/post", middlewares.JWTAuthMiddleware(), controller.CreatePostHandler)
	v1.GET("/post/:id", controller.GetPostDetailHandler)
	v1.GET("/posts/", controller.GetPostListHandler)

	// 投票相关路由
	v1.POST("/vote", middlewares.JWTAuthMiddleware(), controller.PostVoteController)
	// 按热度/时间排序的帖子列表
	v1.GET("/posts2/", controller.GetPostList2Handler)

	pprof.Register(r) // 注册 pprof 路由

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404 not found",
		})
	})

	return r
}
