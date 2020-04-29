package server

import (
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	_ "net/http/pprof"
	"secret-social-network/server/handler"
)

func Run(addr string) {
	r := gin.New()
	ginpprof.Wrap(r)

	r.POST("/init", handler.Init)
	r.POST("/link", handler.Link)
	r.GET("/query", handler.FindAll)
	r.GET("/propagate", handler.Propagate)

	r.Run(addr)
}
