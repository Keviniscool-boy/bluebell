package models

// User 数据库用户模型
type User struct {
	UserID   int64  `json:"user_id,string" db:"user_id"`
	Username string `json:"username" db:"username"`
	Password string `json:"-" db:"password"`
}
