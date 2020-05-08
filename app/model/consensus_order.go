package model

import "github.com/shopspring/decimal"

type ConsensusOrderState string

/**
create order --> wait
wait --confirmed from third--> confirmed (temporary state)
confirmed --create relation log ok--> done
confirmed --create relation log failed--> failed
*/
const (
	ConsensusOrderStateWait      ConsensusOrderState = "wait"      // 待支付
	ConsensusOrderStateConfirmed ConsensusOrderState = "confirmed" // 第三方已经确认
	ConsensusOrderStateDone      ConsensusOrderState = "done"      // 关系建立完成
	ConsensusOrderStateFailed    ConsensusOrderState = "failed"    // 关系建立失败
	ConsensusOrderStateUnlink    ConsensusOrderState = "unlink"    // 已经解除
)

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
