package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"net/http"
	"secret-social-network/app/client"
	"secret-social-network/app/dgraph"
	"secret-social-network/app/server/respond"
	"secret-social-network/app/service"
	"secret-social-network/app/storage"
)

type ConsensusUnlinkRequest struct {
	OpenID1 string `json:"open_id1"`
	OpenID2 string `json:"open_id2"`
	AppID   string `json:"app_id"`
}

func ConsensusUnlink(c *gin.Context) {
	req := ConsensusUnlinkRequest{}
	if err := c.ShouldBind(&req); err != nil {
		respond.Error(c, http.StatusBadRequest, respond.BadRequest(err.Error()))
		return
	}
	// query all exists consensus order
	linkedOrders, err := storage.ConsensusOrder{}.LinkedList(req.AppID, req.OpenID1, req.OpenID2)
	if err != nil {
		respond.Error(c, http.StatusInternalServerError, respond.InternalServerError)
		log.Errorf("get linked list failed: %s", err.Error())
		return
	}

	// prepare unlinked orderIDs and calculate total value
	unlinkedOrderIDs := make([]string, 0, len(linkedOrders))
	unlinkedValue := decimal.Zero
	for _, linkedOrder := range linkedOrders {
		unlinkedOrderIDs = append(unlinkedOrderIDs, linkedOrder.OrderID)
		unlinkedValue = unlinkedValue.Add(*linkedOrder.Value1)
		unlinkedValue = unlinkedValue.Add(*linkedOrder.Value2)
	}

	if unlinkedValue.LessThanOrEqual(decimal.Zero) {
		respond.Error(c, http.StatusBadRequest, respond.NoRelation)
		return
	}

	tx := storage.TxBegin()
	// mark order unlink
	//todo

	// write link log
	value := unlinkedValue.Neg()
	linkLog := storage.ConsensusLinkLog{
		OpenID1: req.OpenID1,
		OpenID2: req.OpenID2,
		Action:  storage.ConsensusLinkLogActionUnlink,
		Value:   &value,
		AppID:   req.AppID,
		State:   storage.ConsensusLinkLogStateWait,
	}
	if err := linkLog.CreateTx(tx); err != nil {
		respond.Error(c, http.StatusInternalServerError, respond.InternalServerError)
		tx.Rollback()
		return
	}

	tx.Commit()

	// rewrite link value at d-graph
	UID1, UID2, err := client.GetUID(req.AppID, req.OpenID1, req.OpenID2)
	if err != nil {
		respond.Error(c, http.StatusInternalServerError, respond.InternalServerError)
		log.Error("get uid from open platform failed: %s", err.Error())
		return
	}
	valueF, _ := value.Float64()
	err = dgraph.User{}.LinkOrAdd(req.AppID, UID1, UID2, valueF)
	if err != nil {
		service.CommitLog(linkLog, false)
	}
	service.CommitLog(linkLog, true)

	respond.Success(c, map[string]interface{}{
		"order_ids": unlinkedOrderIDs,
	})
}
