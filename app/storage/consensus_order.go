package storage

import (
	"github.com/shopspring/decimal"
	"secret-social-network/app/model"
	"time"
)

type ConsensusOrder struct {
	ID           uint `gorm:"primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	OrderID      string           `gorm:"not null" sql:"index"`
	OpenID1      string           `gorm:"not null"`
	OpenID2      string           `gorm:"not null"`
	Value1       *decimal.Decimal `sql:"type:decimal(32,16);not null"`
	Value2       *decimal.Decimal `sql:"type:decimal(32,16);not null"`
	AppID        string           `gorm:"not null"`
	TradeNo1     string
	TradeNo2     string
	State        model.ConsensusOrderState `gorm:"not null"`
	FailedReason string
	Unlinked     bool `gorm:"not null;default:false"`
}

func (c *ConsensusOrder) Create() error {
	return gormDb.Create(c).Error
}

func (c *ConsensusOrder) Query() error {
	return gormDb.Model(c).Where(c).Last(c).Error
}

func (c *ConsensusOrder) ChState(old, new model.ConsensusOrderState) (int64, error) {
	db := gormDb.Model(c).
		Where("id = ? and state = ?", c.ID, old).
		UpdateColumn("state", new)
	return db.RowsAffected, db.Error
}

// LinkedList list all confirm and not unlinked orders.
func (ConsensusOrder) LinkedList(appID string, openID1, openID2 string) ([]ConsensusOrder, error) {
	rs := make([]ConsensusOrder, 0)
	err := gormDb.Model(new(ConsensusOrder)).
		Where("((open_id1 = ? and open_id2 = ?) or (open_id2 = ? and open_id1 = ?)) and unlinked = false").
		Scan(&rs).Error
	return rs, err
}

func (c *ConsensusOrder) Set() model.ConsensusOrder {
	return model.ConsensusOrder{
		CreatedAt: c.CreatedAt.Unix(),
		UpdatedAt: c.UpdatedAt.Unix(),
		OrderID:   c.OrderID,
		OpenID1:   c.OpenID1,
		OpenID2:   c.OpenID2,
		Value1:    c.Value1,
		Value2:    c.Value2,
		AppID:     c.AppID,
		TradeNo1:  c.TradeNo1,
		TradeNo2:  c.TradeNo2,
	}
}
