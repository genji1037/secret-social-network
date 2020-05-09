package server

import (
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"secret-social-network/app/server/handler"
)

// Run runs REST server.
func Run(addr string) error {
	r := gin.New()
	ginpprof.Wrap(r)

	ssns := r.Group("/ssns")
	{
		consensus := ssns.Group("/consensus")
		// 建立共识关系
		consensus.PUT("/relation", handler.ConsensusLink)
		// 确认共识关系
		consensus.POST("/:orderID/confirm", handler.ConsensusConfirm)
		// 断开共识关系
		consensus.DELETE("/relation", handler.ConsensusUnlink)
	}

	//r.POST("/init", handler.Init)
	//r.GET("/propagate", handler.Propagate)

	return r.Run(addr)
}
