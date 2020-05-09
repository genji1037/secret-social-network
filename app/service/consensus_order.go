package service

import (
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"secret-social-network/app/client"
	"secret-social-network/app/config"
	"secret-social-network/app/model"
	"secret-social-network/app/storage"
	"secret-social-network/app/util"
	"strconv"
)

// CreateConsensusOrder create consensus link order.
func CreateConsensusOrder(appID, openID1, openID2 string, value1, value2 decimal.Decimal) (*model.ConsensusOrder, error) {
	cfg := config.GetServe().Consensus
	orderID := strconv.Itoa(int(util.Generate()))

	// apply trade_no at open platform
	var err error
	var tradeNo1, tradeNo2 string
	if value1.GreaterThan(decimal.Zero) {
		args := client.ApplyPaymentArgs{
			AppID:   appID,
			OpenID:  openID1,
			OrderID: "A" + orderID,
			Token:   cfg.Token,
			Amount:  value1,
			Remark:  cfg.PaymentRemark,
		}
		tradeNo1, err = client.ApplyPayment(args)
		if err != nil {
			log.WithField("args", args).Errorf("create consensus order1, apply payment failed: %s", err.Error())
			return nil, err
		}
	}

	if value2.GreaterThan(decimal.Zero) {
		args := client.ApplyPaymentArgs{
			AppID:   appID,
			OpenID:  openID1,
			OrderID: "B" + orderID,
			Token:   cfg.Token,
			Amount:  value1,
			Remark:  cfg.PaymentRemark,
		}
		tradeNo2, err = client.ApplyPayment(args)
		if err != nil {
			log.WithField("args", args).Errorf("create consensus order2, apply payment failed: %s", err.Error())
			return nil, err
		}
	}

	// persist
	order := storage.ConsensusOrder{
		OrderID:     orderID,
		OpenID1:     openID1,
		OpenID2:     openID2,
		Value1:      &value1,
		Value2:      &value2,
		TradeNo1:    tradeNo1,
		TradeNo2:    tradeNo2,
		AppID:       appID,
		LinkState:   model.ConsensusOrderLinkStateWait,
		UnlinkState: model.ConsensusOrderUnlinkStateNone,
	}
	if err := order.Create(); err != nil {
		return nil, err
	}

	orderResp := order.Set()
	return &orderResp, nil
}

func CommitLink(order storage.ConsensusOrder, ok bool) {
	state := model.ConsensusOrderLinkStateDone
	if !ok {
		state = model.ConsensusOrderLinkStateFailed
	}
	rowAffected, err := order.ChLinkState(model.ConsensusOrderLinkStateConfirmed, state)
	if err != nil || rowAffected == 0 {
		log.Errorf("commit link order %d [%s] failed", order.OrderID, state)
	}
}

func CommitUnlink(ids []uint, ok bool) {
	state := model.ConsensusOrderUnlinkStateDone
	if !ok {
		state = model.ConsensusOrderUnlinkStateFailed
	}
	err := storage.ConsensusOrder{}.BatchChUnlinkState(nil, ids, state)
	if err != nil {
		log.Errorf("commit unlink order %+v [%s] failed", ids, state)
	}
}
