package model

import "github.com/shopspring/decimal"

// ConsensusOrderLinkState represent consensus order's link state.
type ConsensusOrderLinkState string

// ConsensusOrderUnlinkState represent consensus order's unlink state.
type ConsensusOrderUnlinkState string

/**
create order --> wait
wait --confirmed from third--> confirmed (temporary state)
confirmed --create relation log ok--> done
confirmed --create relation log failed--> failed
*/
const (
	ConsensusOrderLinkStateWait      ConsensusOrderLinkState = "wait"      // 待支付
	ConsensusOrderLinkStateConfirmed ConsensusOrderLinkState = "confirmed" // 第三方已经确认
	ConsensusOrderLinkStateDone      ConsensusOrderLinkState = "done"      // 关系建立完成
	ConsensusOrderLinkStateFailed    ConsensusOrderLinkState = "failed"    // 关系建立失败

	ConsensusOrderUnlinkStateNone      ConsensusOrderUnlinkState = "none"      // 未解除
	ConsensusOrderUnlinkStateConfirmed ConsensusOrderUnlinkState = "confirmed" // 确认断开
	ConsensusOrderUnlinkStateDone      ConsensusOrderUnlinkState = "done"      // 关系断开完成
	ConsensusOrderUnlinkStateFailed    ConsensusOrderUnlinkState = "failed"    // 关系断开失败

)

// ConsensusOrder represent consensus order data model.
type ConsensusOrder struct {
	CreatedAt int64            `json:"created_at"`
	UpdatedAt int64            `json:"updated_at"`
	OrderID   string           `json:"order_id"`
	OpenID1   string           `json:"open_id1"`
	OpenID2   string           `json:"open_id2"`
	Value1    *decimal.Decimal `json:"value1"`
	Value2    *decimal.Decimal `json:"value2"`
	AppID     string           `json:"app_id"`
	TradeNo1  string           `json:"trade_no1"`
	TradeNo2  string           `json:"trade_no2"`
}
