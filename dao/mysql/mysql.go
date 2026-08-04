package mysql

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var db *sqlx.DB

// Init 初始化 MySQL 连接
func Init() error {
	// 读取配置
	host := viper.GetString("mysql.host")
	port := viper.GetInt("mysql.port")
	user := viper.GetString("mysql.user")
	password := viper.GetString("mysql.password")
	dbname := viper.GetString("mysql.dbname")

	// 构建 DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	// 连接数据库
	var err error
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return fmt.Errorf("mysql connect failed: %w", err)
	}

	// 连接池配置
	db.SetMaxOpenConns(100)                 // 最大打开连接数
	db.SetMaxIdleConns(10)                  // 最大空闲连接数
	db.SetConnMaxLifetime(30 * time.Minute) // 连接最大存活时间

	// 测试连接
	if err = db.Ping(); err != nil {
		return fmt.Errorf("mysql ping failed: %w", err)
	}

	zap.L().Info("MySQL connected",
		zap.String("host", host),
		zap.Int("port", port),
		zap.String("dbname", dbname),
	)

	return nil
}

// Close 关闭数据库连接
func Close() {
	if db != nil {
		if err := db.Close(); err != nil {
			zap.L().Error("mysql close failed", zap.Error(err))
		}
	}
}

// GetDB 获取数据库实例
func GetDB() *sqlx.DB {
	return db
}
