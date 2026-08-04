package models

// 定义一个结构体，用于接收前端传来的参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// 登录请求
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 投票请求
type VoteData struct {
	PostID    int64 `json:"post_id,string" binding:"required"`       // 帖子id（雪花id字符串）
	Direction int8  `json:"direction,string" binding:"oneof=1 0 -1"` // 投票方向 1赞成 0取消 -1反对（不能加 required，0 是合法值）
}

// ParamPostList 帖子列表请求参数
// page 从 1 开始，size 为每页数量，未传时由 handler 填充默认值
type ParamPostList struct {
	Page  int64  `json:"page" form:"page"`
	Size  int64  `json:"size" form:"size"`
	Order string `json:"order" form:"order"` // 排序方式：time 按时间 | score 按热度
}
