package service

import (
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"secret-social-network/app/dgraph"
	"secret-social-network/app/model"
	"secret-social-network/app/storage"
)

func CreateConsensusRelation(order storage.ConsensusOrder) {
	// TODO: gracefully shutdown
	// link log
	value1, value2 := decimal.Zero, decimal.Zero
	if order.Value1 != nil {
		value1 = *order.Value1
	}
	if order.Value2 != nil {
		value2 = *order.Value2
	}
	value := value1.Add(value2)
	linkLog := storage.ConsensusLinkLog{
		OpenID1: order.OpenID1,
		OpenID2: order.OpenID2,
		Action:  storage.ConsensusLinkLogActionLink,
		Value:   &value,
		AppID:   order.AppID,
	}
	if err := linkLog.Create(); err != nil {
		log.WithField("linkLog", linkLog).Errorf("create link log failed: %s", err.Error())
		markOrder(order, model.ConsensusOrderStateFailed)
		return
	}
	markOrder(order, model.ConsensusOrderStateDone)

	// link at d-graph
	valueF, _ := value.Float64()
	err := dgraph.User{}.LinkOrAdd(order.AppID, order.OpenID1, order.OpenID2, valueF)
	if err != nil {
		log.WithFields(log.Fields{
			"app_id":   order.AppID,
			"open_id1": order.OpenID1,
			"open_id2": order.OpenID2,
			"value":    valueF,
		}).Errorf("link or add failed: %s", err.Error())
		CommitLog(linkLog, false)
		return
	}
	CommitLog(linkLog, true)
}
