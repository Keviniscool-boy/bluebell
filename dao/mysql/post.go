package mysql

import (
	"database/sql"
	"errors"

	"bluebell/models"

	"go.uber.org/zap"
)

var ErrPostNotExist = errors.New("帖子不存在")

func CreatePost(p *models.Post) (err error) {
	sqlStr := "INSERT INTO post(post_id,title,content,author_id,community_id) VALUES (?,?,?,?,?)"
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostDetailByID 根据帖子id查询帖子详情（单表查询，作者/社区信息由 logic 层另查组装）
func GetPostDetailByID(id int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `SELECT post_id, title, content, author_id, community_id, status, create_time
	       FROM post
	       WHERE post_id = ?`
	if err = db.Get(post, sqlStr, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			zap.L().Warn("post not found", zap.Int64("id", id))
			return nil, ErrPostNotExist
		}
		zap.L().Error("query post detail failed", zap.Error(err))
		return nil, err
	}
	return post, nil
}

func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `SELECT post_id, title, content, author_id, community_id, status, create_time
	       FROM post ORDER BY create_time DESC LIMIT ? OFFSET ?`
	posts = make([]*models.Post, 0, size)
	if err = db.Select(&posts, sqlStr, size, (page-1)*size); err != nil {
		zap.L().Error("query post list failed", zap.Error(err))
		return nil, err
	}
	return posts, nil
}
