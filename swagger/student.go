package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int    `json:"id"`   // 用戶ID
	Name string `json:"name"` // 姓名
	Age  int    `json:"age"`  // 年齡
}

// Param參數類型有: header, body(post參數), query(get參數), path(restful)
// @Summary 取得用戶
// @Produce json
// @Param   id  path     int    true "用戶ID"
// @Success 200 {object} User   "成功"
// @Failure 400 {object} string "參數錯誤"
// @Failure 500 {object} string "內部錯誤"
// @Router  /get/{id} [GET]
func GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if id, err := strconv.Atoi(idStr); err != nil {
		ctx.String(http.StatusBadRequest, "參數錯誤")
	} else {
		ctx.JSON(http.StatusOK, User{ID: id, Name: "小陳", Age: 20})
	}
}

// @Summary 更新用戶
// @Produce json
// @Param   user body     User   true "用戶信息"
// @Success 200  {object} string "更新成功"
// @Failure 400  {object} string "參數錯誤"
// @Failure 500  {object} string "內部錯誤"
// @Router  /update_user [POST]
func UpdateUser(ctx *gin.Context) {
	var user User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.String(http.StatusBadRequest, "參數錯誤")
	} else {
		ctx.String(http.StatusOK, "更新成功")
	}
}
