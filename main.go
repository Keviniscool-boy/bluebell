package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bluebell/controller"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/routes"
	"bluebell/settings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// @title bluebell 项目接口文档
// @version 1.0
// @description 基于 Gin + MySQL + Redis 的社区 Web 应用，提供用户、社区、帖子、投票等接口
// @termsOfService http://swagger.io/terms/

// @contact.name bluebell
// @contact.url https://github.com/Keviniscool-boy?tab=repositories
// @contact.email raotao29@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description 请求头中携带 Bearer Token，例如：Bearer &lt;token&gt;
func main() {
	// 1.加载配置文件
	if err := settings.Init(); err != nil {
		fmt.Printf("settings.Init() failed, err:%v\n", err)
		return
	}

	// 2.初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("logger.Init() failed, err:%v\n", err)
		return
	}
	defer logger.Close()

	// 3.初始化数据库
	if err := mysql.Init(); err != nil {
		fmt.Printf("mysql.Init() failed, err:%v\n", err)
		return
	}
	defer mysql.Close()

	// 4.初始化 Redis
	if err := redis.Init(); err != nil {
		fmt.Printf("redis.Init() failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	// 5.初始化雪花算法
	if err := snowflake.Init(settings.Conf.Snowflake.StartTime, settings.Conf.Snowflake.MachineID); err != nil {
		fmt.Printf("snowflake.Init() failed, err:%v\n", err)
		return
	}

	// 6.初始化 validator 翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("controller.InitTrans() failed, err:%v\n", err)
		return
	}

	// 7.注册路由
	r := routes.Setup()

	// 8.启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		zap.L().Info("server starting", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("server startup failed", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}
	logger.Info("server exited")
}
