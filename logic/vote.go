package logic

import (
	"errors"
	"strconv"

	"bluebell/dao/redis"
)

// 投票业务错误
var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrRepeatedVote   = errors.New("不允许重复投票")
)

// VoteForPost 为帖子投票，逻辑层负责 int64→string 转换和 redis 错误映射
func VoteForPost(userID, postID int64, direction int8) error {
	err := redis.VoteForPost(strconv.FormatInt(userID, 10), strconv.FormatInt(postID, 10), float64(direction))
	if errors.Is(err, redis.ErrPostNotExist) {
		return ErrPostNotExist
	}
	if errors.Is(err, redis.ErrVoteTimeExpire) {
		return ErrVoteTimeExpire
	}
	if errors.Is(err, redis.ErrRepeatedVote) {
		return ErrRepeatedVote
	}
	return err
}
