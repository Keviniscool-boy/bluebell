package redis

// Redis 相关的 key 常量（使用时用 getRedisKey 拼接 KeyPrefix 得到完整 key）
const (
	KeyPrefix        = "bluebell:"
	KeyPostTimeZSet  = "post:time"   // zset 帖子及发帖时间
	KeyPostScoreZSet = "post:score"  // zset 帖子及投票分数
	KeyPostVotedPF   = "post:voted:" // zset 保存每个帖子被投票的用户及投票类型
)

// getRedisKey 拼接完整的 redis key
func getRedisKey(key string) string {
	return KeyPrefix + key
}
