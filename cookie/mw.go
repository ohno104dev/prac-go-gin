package main

import (
	"fmt"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

func authMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId := genSessionId(ctx)
		exists := false

		for _, cookie := range ctx.Request.Cookies() {
			if cookie.Name == authCookie && cookie.Value == sessionId {
				exists = true
				break
			}
		}

		if !exists {
			ctx.Redirect(http.StatusMovedPermanently, GO_TO_LOGIN)
			ctx.Abort()
		}

		if v, ok := userInfos.Load(sessionId); !ok {
			fmt.Printf("session id %s 不存在\n", sessionId)
			ctx.String(http.StatusUnauthorized, "身份驗證失敗")
			ctx.Abort()
		} else {
			var user User
			sonic.Unmarshal(v.([]byte), &user)
			ctx.Set("name", user.Name)
			ctx.Set("role", user.Role)
			ctx.Set("vip", user.Vip)
		}
	}
}
