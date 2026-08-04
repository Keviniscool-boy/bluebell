package redis

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// 投票算法常量
const (
	oneWeekInSeconds = 7 * 24 * 3600 // 投票有效期：发帖一周内允许投票
	scorePerVote     = 432           // 每一票的分数：432 = 86400/200，即约200票可抵消一天的帖子年龄
)

// 投票业务错误
var (
	ErrPostNotExist   = errors.New("帖子不存在")
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrRepeatedVote   = errors.New("不允许重复投票")
)

// CreatePostTimeAndScore 创建帖子时把发帖时间与初始分数写入zset（TxPipeline原子执行）
func CreatePostTimeAndScore(postID string, createTime float64) error {
	ctx := context.Background()
	pipeline := rdb.TxPipeline()
	pipeline.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), redis.Z{Score: createTime, Member: postID})
	pipeline.ZAdd(ctx, getRedisKey(KeyPostScoreZSet), redis.Z{Score: createTime, Member: postID})
	if _, err := pipeline.Exec(ctx); err != nil {
		zap.L().Error("redis CreatePostTimeAndScore pipeline failed",
			zap.String("postID", postID), zap.Error(err))
		return err
	}
	return nil
}

// VoteForPost 为帖子投票（value: 1赞成 / -1反对 / 0取消）
func VoteForPost(userID, postID string, value float64) error {
	ctx := context.Background()

	// 1. 判断帖子是否存在（是否已在redis登记发帖时间）
	postTime, err := rdb.ZScore(ctx, getRedisKey(KeyPostTimeZSet), postID).Result()
	if errors.Is(err, redis.Nil) {
		return ErrPostNotExist
	}
	if err != nil {
		return err
	}
	// 2. 判断是否在投票有效期内（发帖一周内）
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	// 3. 获取该用户对同一帖子的旧投票值（从未投票视为0）
	oldValue, err := rdb.ZScore(ctx, getRedisKey(KeyPostVotedPF)+postID, userID).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	// 这一次投票和之前一致则视为重复投票
	if oldValue == value {
		return ErrRepeatedVote
	}

	// 4. 计算分数变化（从赞改为踩会抵消两票的分数）
	var dir float64
	if value > oldValue {
		dir = 1
	} else {
		dir = -1
	}
	diff := math.Abs(value - oldValue)
	scoreDelta := dir * diff * scorePerVote

	// 5. TxPipeline 原子更新帖子分数与用户的投票记录
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(ctx, getRedisKey(KeyPostScoreZSet), scoreDelta, postID)
	if value == 0 {
		// 取消投票：移除该用户的投票记录
		pipeline.ZRem(ctx, getRedisKey(KeyPostVotedPF)+postID, userID)
	} else {
		pipeline.ZAdd(ctx, getRedisKey(KeyPostVotedPF)+postID, redis.Z{Score: value, Member: userID})
	}
	if _, err := pipeline.Exec(ctx); err != nil {
		zap.L().Error("redis VoteForPost pipeline failed",
			zap.String("postID", postID), zap.String("userID", userID), zap.Error(err))
		return err
	}
	return nil
}

// GetPostIDsInOrder 按分数从大到小（score=热度）或按时间（time=最新）返回帖子ID列表（分页）
func GetPostIDsInOrder(order string, page, size int64) ([]string, error) {
	ctx := context.Background()
	key := getRedisKey(KeyPostTimeZSet)
	if order == "score" {
		key = getRedisKey(KeyPostScoreZSet)
	}
	// ZRangeArgs{Rev: true} 等价于 ZREVRANGE，避免使用已废弃的 ZRevRange
	start := (page - 1) * size
	stop := start + size - 1
	ids, err := rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   key,
		Start: start,
		Stop:  stop,
		Rev:   true,
	}).Result()
	if err != nil {
		zap.L().Error("redis GetPostIDsInOrder failed", zap.String("order", order), zap.Error(err))
		return nil, err
	}
	return ids, nil
}
