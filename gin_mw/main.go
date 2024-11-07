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
		ctx.Next()
		elapsed := time.Since(begin)
		log.Printf("%s use time %d ms\n", ctx.Request.URL.Path, elapsed.Milliseconds())
	}
}

var limitCh = make(chan struct{}, 10)

func LimitMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limitCh <- struct{}{}
		ctx.Next()
		<-limitCh
	}
}

func main() {
	engine := gin.Default()
	engine.Use(LimitMW()) //global middlewares
	engine.GET("/boy", timeMW(), func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "gin, hi boy")
	})

	engine.Run("127.0.0.1:8000")
}
