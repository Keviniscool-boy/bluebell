package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// token过期时间
const TokenExpireDuration = time.Hour * 2 //2小时

func signingSecret() ([]byte, error) {
	secret := viper.GetString("jwt.secret")
	if len(secret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must contain at least 32 characters")
	}
	return []byte(secret), nil
}

type MyClaims struct {
	UserId   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(userId int64, userName string) (string, error) {
	secret, err := signingSecret()
	if err != nil {
		return "", err
	}

	//创建一个我们自己的声明
	c := MyClaims{
		UserId:   userId,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), //过期时间
			Issuer:    "bluebell",                                 //签发人
		},
	}
	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	//使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(secret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	secret, err := signingSecret()
	if err != nil {
		return nil, err
	}

	// 解析token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, e error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if token.Valid {
		return mc, nil
	}
	return nil, fmt.Errorf("invalid token")
}
