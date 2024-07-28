package mw

import (
	"errors"
	"fast/pkg/db"
	"fast/utils"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type UserInfo struct {
	db.User
}

type CustomClaims struct {
	UserInfo
	jwt.StandardClaims
}

// MySecret 密钥
var secret = utils.GetRandomString(32)

// OutTime 超时时间（分钟）
const OutTime = time.Duration(60 * 24 * 7)

// GenToken 创建Token
func GenToken(user UserInfo) (string, error) {
	claim := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Minute * OutTime)), // 过期时间
			Issuer:    "",                                            // 签发人
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析token
func ParseToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println(" token parse err:", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新token
func RefreshToken(tokenStr string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = jwt.At(time.Now().Add(time.Minute * OutTime))
		return GenToken(claims.UserInfo)
	}

	return "", errors.New("cloudn't handle this token")
}
