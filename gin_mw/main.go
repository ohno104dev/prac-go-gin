package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func timeMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		begin := time.Now()
		ctx.Next() // 往下一個執行串列調用
		elapsed := time.Since(begin)
		log.Printf("%s use time %d ms\n", ctx.Request.URL.Path, elapsed.Milliseconds())
		ctx.String(200, " hihihi")
	}
}

var limitCh = make(chan struct{}, 10)

func limitMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limitCh <- struct{}{}
		log.Printf("concurrence %d\n", len(limitCh))
		ctx.Next()
		<-limitCh
	}
}

func boy(ctx *gin.Context) {
	ctx.String(http.StatusOK, "gin, hi boy")
	ctx.Abort() // 中止往後的串列調用
}

func main() {
	engine := gin.Default()
	// engine.Use(LimitMW()) //global middlewares
	engine.GET("/boy", timeMW(), boy, limitMW())

	engine.Run("127.0.0.1:8000")
}
