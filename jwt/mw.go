package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	jwt "github.com/ohno104dev/prac-go-gin/self_jwt"
)

func jwtAuthMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		_, payload, err := jwt.VerifyJwt(token)
		if err != nil {
			ctx.String(http.StatusForbidden, "auth failed")
			ctx.Abort()
			return
		}

		for k, v := range payload.UserDefined {
			fmt.Println(k+":", v)
			ctx.Set(k, v)
		}
		ctx.Next()
	}
}
