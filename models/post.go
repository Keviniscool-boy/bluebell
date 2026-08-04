package models

import (
	"encoding/json"
	"time"
)

type Post struct {
	ID          int64     `json:"id,string" db:"post_id"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	AuthorID    int64     `json:"author_id,string" db:"author_id"`
	CommunityID int64     `json:"community_id,string" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

// ApiPostDetail 帖子详情接口返回模型（Post + 作者名 + 社区详情）
type ApiPostDetail struct {
	AuthorName       string `json:"author_name"`
	*Post                   // 嵌入结构体
	*CommunityDetail        // 嵌入结构体
}

// MarshalJSON 手动组装帖子详情 JSON。
// Post 与 CommunityDetail 内嵌后存在同名字段（id、create_time），
// 直接序列化时这些冲突字段会被 encoding/json 整体丢弃，
// 导致帖子 id、创建时间丢失。这里用同级直接字段覆盖内嵌的同名字段，
// 同时保证 id 以字符串输出，避免 JS 对超过 2^53 的雪花 ID 丢失精度。
func (a *ApiPostDetail) MarshalJSON() ([]byte, error) {
	post := a.Post
	if post == nil {
		post = &Post{}
	}
	type Alias ApiPostDetail // 定义别名，避免递归调用本方法
	return json.Marshal(&struct {
		*Alias
		ID         int64     `json:"id,string"`
		CreateTime time.Time `json:"create_time"`
	}{
		Alias:      (*Alias)(a),
		ID:         post.ID,
		CreateTime: post.CreateTime,
	})
}
