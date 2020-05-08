package respond

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Data: data,
	})
}

func Error(c *gin.Context, statusCode int, err BizError) {
	c.JSON(statusCode, Response{
		Code: err.Code,
		Msg:  err.Msg,
	})
}
