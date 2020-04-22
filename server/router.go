package server

import (
	"github.com/gin-gonic/gin"
	"secret-social-network/server/handler"
)

func Run(addr string) {
	r := gin.New()

	r.POST("/init", handler.Init)
	r.POST("/link", handler.Link)
	r.GET("/query", handler.FindAll)

	r.Run(addr)
}
