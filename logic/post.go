package logic

import (
	"errors"
	"strconv"
	"time"

	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

// 业务错误定义
var ErrPostNotExist = errors.New("帖子不存在")

func CreatePost(p *models.Post) (err error) {
	// 生成post_id 雪花算法生成id
	p.ID = snowflake.GenerateID()

	// 保存到数据库
	if err = mysql.CreatePost(p); err != nil {
		return err
	}
	// 同步把发帖时间和初始分数写入redis，供投票与热度排序使用
	return redis.CreatePostTimeAndScore(strconv.FormatInt(p.ID, 10), float64(time.Now().Unix()))
}

// GetPostDetail 根据帖子id查询帖子详情
func GetPostDetail(id int64) (data *models.ApiPostDetail, err error) {
	post, err := mysql.GetPostDetailByID(id)
	if err != nil {
		if errors.Is(err, mysql.ErrPostNotExist) {
			return nil, ErrPostNotExist
		}
		return nil, err
	}
	// 根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		return nil, err
	}
	// 根据社区id查询社区信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		if errors.Is(err, mysql.ErrCommunityNotExist) {
			return nil, ErrCommunityNotExist
		}
		return nil, err
	}
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return data, nil
}

func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	// 查询帖子列表（分页）
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			return nil, err
		}
		// 根据社区id查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			if errors.Is(err, mysql.ErrCommunityNotExist) {
				return nil, ErrCommunityNotExist
			}
			return nil, err
		}
		data = append(data, &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		})
	}
	return data, nil
}

// GetPostList2 从redis获取帖子id（按score热度或time时间排序，支持分页），再逐个组装详情返回
func GetPostList2(page, size int64, order string) (data []*models.ApiPostDetail, err error) {
	ids, err := redis.GetPostIDsInOrder(order, page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(ids))
	for _, id := range ids {
		pid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			zap.L().Warn("invalid post id in redis", zap.String("id", id))
			continue
		}
		post, err := GetPostDetail(pid)
		if err != nil {
			// 帖子可能在redis中但已从mysql删除，跳过即可
			zap.L().Warn("get post detail failed", zap.Int64("pid", pid), zap.Error(err))
			continue
		}
		data = append(data, post)
	}
	return data, nil
}
