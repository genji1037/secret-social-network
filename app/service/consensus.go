package service

import (
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"secret-social-network/app/dgraph"
	"secret-social-network/app/storage"
)

func CreateConsensusRelation(order storage.ConsensusOrder) {
	// TODO: gracefully shutdown

	// link at d-graph
	value1, value2 := decimal.Zero, decimal.Zero
	if order.Value1 != nil {
		value1 = *order.Value1
	}
	if order.Value2 != nil {
		value2 = *order.Value2
	}
	value := value1.Add(value2)
	valueF, _ := value.Float64()
	err := dgraph.User{}.LinkOrAdd(order.AppID, order.OpenID1, order.OpenID2, valueF)
	ok := true
	if err != nil {
		log.WithFields(log.Fields{
			"app_id":   order.AppID,
			"open_id1": order.OpenID1,
			"open_id2": order.OpenID2,
			"value":    valueF,
		}).Errorf("link or add failed: %s", err.Error())
		ok = false
	}
	CommitLink(order, ok)
}
