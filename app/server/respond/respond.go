package respond

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response is ssns common respond.
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Success is ssns common success respond.
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Data: data,
	})
}

// Error is ssns common failed respond.
func Error(c *gin.Context, statusCode int, err BizError) {
	c.JSON(statusCode, Response{
		Code: err.Code,
		Msg:  err.Msg,
	})
}
