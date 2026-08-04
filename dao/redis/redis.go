package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var rdb *redis.Client

// Init 初始化 Redis 连接
func Init() error {
	host := viper.GetString("redis.host")
	port := viper.GetInt("redis.port")
	db := viper.GetInt("redis.db")
	password := viper.GetString("redis.password") // 没有密码则为空

	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
		PoolSize: 100, // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试连接
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis connect failed: %w", err)
	}

	zap.L().Info("Redis connected",
		zap.String("host", host),
		zap.Int("port", port),
		zap.Int("db", db),
	)

	return nil
}

// Close 关闭 Redis 连接
func Close() {
	if rdb != nil {
		if err := rdb.Close(); err != nil {
			zap.L().Error("redis close failed", zap.Error(err))
		}
	}
}

// GetClient 获取 Redis 客户端实例
func GetClient() *redis.Client {
	return rdb
}
