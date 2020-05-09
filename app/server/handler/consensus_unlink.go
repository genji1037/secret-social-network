package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"net/http"
	"secret-social-network/app/client"
	"secret-social-network/app/dgraph"
	"secret-social-network/app/model"
	"secret-social-network/app/server/respond"
	"secret-social-network/app/service"
	"secret-social-network/app/storage"
)

type consensusUnlinkRequest struct {
	OpenID1 string `json:"open_id1"`
	OpenID2 string `json:"open_id2"`
	AppID   string `json:"app_id"`
}

// ConsensusUnlink represent consensus unlink handler.
func ConsensusUnlink(c *gin.Context) {
	req := consensusUnlinkRequest{}
	if err := c.ShouldBind(&req); err != nil {
		respond.Error(c, http.StatusBadRequest, respond.BadRequest(err.Error()))
		return
	}
	// query all exists consensus order
	tx := storage.TxBegin()
	linkedOrders, err := storage.ConsensusOrder{}.LinkedListForUpdate(tx, req.AppID, req.OpenID1, req.OpenID2)
	if err != nil {
		tx.Rollback()
		respond.Error(c, http.StatusInternalServerError, respond.InternalServerError)
		log.Errorf("get linked list failed: %s", err.Error())
		return
	}

	// prepare unlinked orderIDs and calculate total value
	unlinkedIDs := make([]uint, 0, len(linkedOrders))
	unlinkedOrderIDs := make([]string, 0, len(linkedOrders))
	unlinkedValue := decimal.Zero
	for _, linkedOrder := range linkedOrders {
		unlinkedIDs = append(unlinkedIDs, linkedOrder.ID)
		unlinkedOrderIDs = append(unlinkedOrderIDs, linkedOrder.OrderID)
		unlinkedValue = unlinkedValue.Add(*linkedOrder.Value1)
		unlinkedValue = unlinkedValue.Add(*linkedOrder.Value2)
	}

	if unlinkedValue.LessThanOrEqual(decimal.Zero) {
		tx.Rollback()
		respond.Error(c, http.StatusBadRequest, respond.NoRelation)
		return
	}

	// get uid from open platform.
	UID1, UID2, err := client.GetUID(req.AppID, req.OpenID1, req.OpenID2)
	if err != nil {
		tx.Rollback()
		respond.Error(c, http.StatusInternalServerError, respond.InternalServerError)
		log.Errorf("get uid from open platform failed: %s", err.Error())
		return
	}

	// mark order unlink
	err = storage.ConsensusOrder{}.BatchChUnlinkState(tx, unlinkedIDs, model.ConsensusOrderUnlinkStateConfirmed)
	if err != nil {
		tx.Rollback()
		log.WithField("orderIDs", unlinkedOrderIDs).Warnf("batch unlink failed: %s", err.Error())
		respond.Error(c, http.StatusInternalServerError, respond.InternalServerError)
		return
	}
	tx.Commit()

	// write link log
	value := unlinkedValue.Neg()

	// rewrite link value at d-graph
	valueF, _ := value.Float64()
	err = dgraph.LinkOrAdd(req.AppID, UID1, UID2, valueF)
	ok := true
	if err != nil {
		ok = false
	}
	service.CommitUnlink(unlinkedIDs, ok)

	respond.Success(c, map[string]interface{}{
		"order_ids": unlinkedOrderIDs,
	})
}
