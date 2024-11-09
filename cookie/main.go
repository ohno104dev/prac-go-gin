package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

var (
	userInfos sync.Map // 模擬緩存
)

const (
	authCookie  = "auth"
	GO_TO_LOGIN = "http://localhost:8000/go_to_login"
)

func genSessionId(ctx *gin.Context) string {
	return base64.StdEncoding.EncodeToString([]byte(ctx.Request.RemoteAddr))
}

type User struct {
	Name string `json:"name"`
	Role string `json:"role"`
	Vip  bool   `json:"vip"`
}

func login(ctx *gin.Context) {
	sessionId := genSessionId(ctx)
	user := User{
		Name: "小歐",
		Role: "一般會員",
		Vip:  false,
	}
	userInfo, _ := sonic.Marshal(&user)
	userInfos.Store(sessionId, userInfo)

	ctx.SetCookie(
		authCookie,  // cookie name
		sessionId,   // cookie value
		3000,        // maxAge
		"/",         // path, cookie存放目錄
		"localhost", // 從屬domain
		false,       // 是否https only
		true,        // 是否只允許http情求(js不行)
	)
	fmt.Printf("set cookie %s = %s to client\n", authCookie, sessionId)
	ctx.String(http.StatusOK, "登入成功")
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
	cwd, _ := os.Getwd()
	fmt.Println("Current Working Directory:", cwd)
	engine := gin.Default()

	engine.LoadHTMLFiles("./templates/go_to_login.html")
	// engine.LoadHTMLGlob("./templates/*.html") // for all
	engine.GET("/go_to_login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "go_to_login.html", gin.H{})
	})

	engine.GET("/login", login)
	engine.GET("/home", authMW(), myHomePage)
	engine.GET("/post", authMW(), postVideo)
	engine.Run("localhost:8000")
}
