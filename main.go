package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Student struct {
	Name  string `form:"name" binding:"required"`
	Score int    `form:"score" binding:"gt=0"`

	Enrollment time.Time `form:"enrollment" binding:"required,before_today" time_format:"2006-01-02" time_utc:"8"`
	Graduation time.Time `form:"graduation" binding:"required,gtfield=Enrollment" time_format:"2006-01-02" time_utc:"8"`
}

var beforeToday validator.Func = func(f1 validator.FieldLevel) bool {
	if date, ok := f1.Field().Interface().(time.Time); ok {
		today := time.Now()
		if date.Before(today) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func processErr(err error) string {
	if err == nil {
		return ""
	}

	invalid, ok := err.(*validator.InvalidValidationError)
	if ok {
		return fmt.Sprintf("param error: %v", invalid)
	}

	validationErrs := err.(validator.ValidationErrors)
	msgs := make([]string, 0, 3)
	for _, validationErr := range validationErrs {
		msgs = append(msgs, fmt.Sprintf("field %s 不滿足條件 %s", validationErr.Field(), validationErr.Tag()))
	}

	return strings.Join(msgs, ";")
}

func main() {
	engine := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("before_today", beforeToday)
	}

	engine.GET("/", func(ctx *gin.Context) {
		var stu Student
		if err := ctx.ShouldBind(&stu); err != nil {
			msg := processErr(err)
			ctx.JSON(http.StatusBadRequest, "parse parameter failed. "+msg)
		} else {
			ctx.JSON(http.StatusOK, stu)
		}
	})

	engine.Run("127.0.0.1:8000")
}
