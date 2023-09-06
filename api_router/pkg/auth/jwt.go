package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

type Claims struct {
	UserId int64 `json:"user_id"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 24 * 30 // 设置过期时间

var Secret = []byte(viper.GetString("server.jwtSecret")) // 设置密码，配置文件中读取

// GenerateToken 签发Token
func GenerateToken(userId int64) (string, error) {
	now := time.Now()
	// 创建一个自己的声明
	claims := Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(TokenExpireDuration).Unix(), // 过期时间： 当前时间 + 过期时间
			Issuer:    "admin",                             // 签发人
		},
	}
	// 创建签名对象
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 获取token
	token, err := tokenClaims.SignedString(Secret)

	return token, err
}

// ParseToken 解析token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, errors.New("invalid token")
}
