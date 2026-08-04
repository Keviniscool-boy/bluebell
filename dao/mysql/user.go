package mysql

import (
	"database/sql"
	"errors"

	"bluebell/models"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserExist          = errors.New("用户已存在")
)

// CheckUserExist 根据用户名查询用户是否存在
// 参数:
//   - username: 要查询的用户名
//
// 返回:
//   - bool: 用户是否存在 (true 表示存在)
//   - error: 查询过程中产生的错误（无错误时为 nil）
func CheckUserExist(username string) (bool, error) {
	sqlStr := "SELECT COUNT(user_id) FROM user WHERE username = ?"
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return false, err
	}
	return count > 0, nil
}

// InsertUser 向数据库中插入一条用户记录
// InsertUser 向数据库中插入一条用户记录
// 该函数会对传入的 `user.Password` 进行 bcrypt 哈希后再存储到数据库
// 参数:
//   - user: 待插入的用户对象，`User.Password` 为明文密码（函数会替换为哈希值写入）
//
// 返回:
//   - error: 插入过程中的错误，成功时为 nil
func InsertUser(user *models.User) error {
	sqlStr := "INSERT INTO user (user_id, username, password) VALUES (?, ?, ?)"
	// 对密码进行哈希再存储
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.Exec(sqlStr, user.UserID, user.Username, string(hashedPassword))
	if err != nil {
		// 并发注册同名用户时可能撞唯一索引 idx_username（MySQL 错误码 1062），映射为业务错误
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return ErrUserExist
		}
		return err
	}
	return nil
}

// Login 根据用户名查询用户并校验密码
// 说明:
//   - 调用时传入的 `user.Password` 应为登录时提交的明文密码，`user.Username` 为要登录的用户名
//   - 本函数会从数据库读取该用户（并填充 `user.UserID`、`user.Username`、`user.Password`，其中 `user.Password` 为数据库中的哈希值）
//   - 然后用 bcrypt 比较数据库中的哈希与传入的明文密码，校验通过返回 nil，失败返回错误
//
// 参数:
//   - user: 传入包含 `Username` 和明文 `Password` 的用户对象；函数会在成功时填充更多字段
//
// 返回:
//   - error: 查询或校验过程中发生的错误（密码错误时也会返回错误）
func Login(user *models.User) (err error) {
	oPassword := user.Password
	sqlStr := "SELECT user_id, username, password FROM user WHERE username = ?"
	if err = db.Get(user, sqlStr, user.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidCredentials
		}
		return err
	}
	// 判断密码是否正确（user.Password 为数据库中读取的哈希）
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oPassword)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		}
		return err
	}
	return nil

}

// GetUserByID 根据用户id查询用户信息
func GetUserByID(id int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := "SELECT user_id, username FROM user WHERE user_id = ?"
	if err = db.Get(user, sqlStr, id); err != nil {
		return nil, err
	}
	return user, nil
}
