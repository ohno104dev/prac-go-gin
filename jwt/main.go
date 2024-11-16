package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/ohno104dev/prac-go-gin/self_jwt"
)

var defaultHeader = jwt.JwtHeader{
	Algo: "HS256",
	Type: "JWT",
}

type User struct {
	Name string `json:"name"`
	Role string `json:"role"`
	Vip  bool   `json:"vip"`
}

func login(ctx *gin.Context) {
	user := User{
		Name: "小歐",
		Role: "一般會員",
		Vip:  false,
	}

	header := defaultHeader
	payload := jwt.JwtPayload{
		Issue:      "yoyo.com",
		IssueAt:    time.Now().Unix(),
		Expiration: time.Now().Add(3 * 24 * time.Hour).Unix(),
		UserDefined: map[string]any{
			"name": user.Name,
			"role": user.Role,
			"vip":  user.Vip,
		},
	}

	if token, err := jwt.GenJwt(header, payload); err != nil {
		log.Println("failed to generate jwt token:", err)
		ctx.String(http.StatusInternalServerError, "jwt token error")
	} else {
		ctx.String(http.StatusOK, token)
	}
}

func myHomePage(ctx *gin.Context) {
	// some code here
	ctx.String(http.StatusOK, fmt.Sprintf("歡迎 %s, 這是你的個人主頁", ctx.GetString("name")))
}

func postVideo(ctx *gin.Context) {
	// some code here
	ctx.String(http.StatusOK, "影片發布成功")
}

func main() {
	engine := gin.Default()

	engine.GET("/login", login)
	engine.GET("/home", jwtAuthMW(), myHomePage)
	engine.GET("/post", jwtAuthMW(), postVideo)
	engine.Run("localhost:8000")
}
