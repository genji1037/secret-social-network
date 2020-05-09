package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"net/http"
	"secret-social-network/app/server/respond"
	"secret-social-network/app/service"
)

type linkRequest struct {
	AppID   string          `json:"app_id" binding:"required"`
	OpenID1 string          `json:"open_id1" binding:"required,uuid"`
	OpenID2 string          `json:"open_id2" binding:"required,uuid"`
	Value1  decimal.Decimal `json:"value1" binding:"required,gte=0"`
	Value2  decimal.Decimal `json:"value2" binding:"required,gte=0"`
}

// ConsensusLink consensus link handler.
func ConsensusLink(c *gin.Context) {
	req := linkRequest{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}

	// validate
	if req.Value1.Equal(decimal.Zero) && req.Value2.Equal(decimal.Zero) {
		respond.Error(c, http.StatusBadRequest, respond.BadRequest("both two values are zero"))
		return
	}

	order, err := service.CreateConsensusOrder(req.AppID, req.OpenID1, req.OpenID2, req.Value1, req.Value2)
	if err != nil {
		log.WithField("args", req).Errorf("[REST] create consensus order failed: %s", err.Error())
		respond.Error(c, http.StatusInternalServerError, respond.InternalServerError)
		return
	}

	respond.Success(c, order)
}
