package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"mall-destroyer/logger"
	"net/http"
	"secret-social-network/app/model"
	"secret-social-network/app/server/respond"
	"secret-social-network/app/service"
	"secret-social-network/app/storage"
	"sync/atomic"
)

// ConsensusConfirm is consensus confirm handler.
func ConsensusConfirm(c *gin.Context) {
	orderID := c.Param("orderID")

	// query order
	order := storage.ConsensusOrder{
		OrderID: orderID,
	}
	if err := order.Query(); err != nil {
		if err == gorm.ErrRecordNotFound {
			respond.Error(c, http.StatusBadRequest, respond.BadRequest("order not found"))
			return
		}
		respond.Error(c, http.StatusInternalServerError, respond.InternalServerError)
		return
	}

	// validate state
	if order.LinkState != model.ConsensusOrderLinkStateWait {
		logger.Warnf("[REST] trying confirm consensus order %s that already confirmed", order.OrderID)
		respond.Error(c, http.StatusBadRequest, respond.AlreadyConfirmed)
		return
	}

	// confirm
	atomic.AddInt64(&service.InflightConsRelation, 1)
	rowAffected, err := order.ChLinkState(model.ConsensusOrderLinkStateWait, model.ConsensusOrderLinkStateConfirmed)
	if err != nil {
		atomic.AddInt64(&service.InflightConsRelation, -1)
		respond.Error(c, http.StatusInternalServerError, respond.InternalServerError)
		return
	}
	if rowAffected == 0 { // confirmed by other thread
		atomic.AddInt64(&service.InflightConsRelation, -1)
		logger.Warnf("[REST] trying confirm consensus order %s that already confirmed by other parallel thread", order.OrderID)
		respond.Error(c, http.StatusBadRequest, respond.AlreadyConfirmed)
		return
	}

	respond.Success(c, nil)

	// create relation asynchronously
	go service.CreateConsensusRelation(order)

}
