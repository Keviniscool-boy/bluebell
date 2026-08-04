package logic

import (
	"errors"

	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

// 业务错误定义
var (
	ErrUserExist          = errors.New("用户已存在")
	ErrPasswordMismatch   = errors.New("两次密码不一致")
	ErrInvalidCredentials = errors.New("用户名或密码错误")
)

// SignUp 用户注册业务逻辑
func SignUp(p *models.ParamSignUp) error {
	// 1. 判断两次密码是否一致
	if p.Password != p.RePassword {
		return ErrPasswordMismatch
	}

	// 2. 查询数据库，查看用户是否存在
	exist, err := mysql.CheckUserExist(p.Username)
	if err != nil {
		return err
	}
	if exist {
		return ErrUserExist
	}

	// 3. 生成 UID
	userID := snowflake.GenerateID()

	// 4. 保存进数据库（密码加密在 mysql 层处理）
	//构造一个 User 实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	if err := mysql.InsertUser(user); err != nil {
		// 并发注册时撞唯一索引，InsertUser 返回 mysql.ErrUserExist，转成业务错误
		if errors.Is(err, mysql.ErrUserExist) {
			return ErrUserExist
		}
		return err
	}
	return nil
}

func Login(p *models.ParamLogin) (token string, err error) {
	//构造一个 User 实例
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 调用 mysql 层的 Login 函数进行登录验证
	if err := mysql.Login(user); err != nil {
		if errors.Is(err, mysql.ErrInvalidCredentials) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}
	//生成jwt
	return jwt.GenToken(user.UserID, user.Username)

}
