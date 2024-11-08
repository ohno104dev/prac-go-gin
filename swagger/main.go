package main

import (
	_ "github.com/ohno104dev/prac-go-gin/swagger/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	engine := gin.Default()
	engine.GET("/swagger/*all", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//

	engine.GET("/get/:id", GetUser)
	engine.POST("/update_user", UpdateUser)

	engine.Run("127.0.0.1:8000")
}
