package mw

import (
	"fast/web/base"
	"strings"
)

func JWTAuth(r *base.Ctx) {
	// 获取token
	authHeader := r.GetHeader("Authorization")

	if authHeader == "" {
		r.ErrorWithCode(401, "请传递正确的验证头信息")
		r.Abort()
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)

	if !(len(parts) == 2 && parts[0] == "Bearer") {
		r.ErrorWithCode(401, "请传递正确的验证头信息")
		r.Abort()
		return
	}

	// token解析
	claims, err := ParseToken(parts[1])
	if err != nil {
		r.ErrorWithCode(401, "请传递正确的验证头信息")
		r.Abort()
		return
	}
	// 存入上下文
	r.Set("user", claims.UserInfo.User)
	r.Next()
}
