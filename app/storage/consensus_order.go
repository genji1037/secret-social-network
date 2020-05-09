package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"secret-social-network/app/model"
	"time"
)

// ConsensusOrder represent mysql table data model.
type ConsensusOrder struct {
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	OrderID     string           `gorm:"not null" sql:"index"`
	OpenID1     string           `gorm:"not null"`
	OpenID2     string           `gorm:"not null"`
	UID1        string           `gorm:"not null" sql:"index"`
	UID2        string           `gorm:"not null" sql:"index"`
	Value1      *decimal.Decimal `sql:"type:decimal(32,16);not null"`
	Value2      *decimal.Decimal `sql:"type:decimal(32,16);not null"`
	AppID       string           `gorm:"not null"`
	TradeNo1    string
	TradeNo2    string
	LinkState   model.ConsensusOrderLinkState `gorm:"not null"`
	LinkedAt    *time.Time
	UnlinkState model.ConsensusOrderUnlinkState `gorm:"not null"`
	UnlinkedAt  *time.Time
}

// Create create consensus order record.
func (c *ConsensusOrder) Create() error {
	return gormDb.Create(c).Error
}

// Query query a consensus order by provide condition.
func (c *ConsensusOrder) Query() error {
	return gormDb.Model(c).Where(c).Last(c).Error
}

// ChLinkState change link state, and return the changed record number.
func (c *ConsensusOrder) ChLinkState(old, new model.ConsensusOrderLinkState) (int64, error) {
	db := gormDb.Model(c).
		Where("id = ? and link_state = ?", c.ID, old)

	if new == model.ConsensusOrderLinkStateConfirmed {
		db = db.Update(map[string]interface{}{
			"link_state": new,
			"linked_at":  time.Now(),
		})
	} else {
		db = db.UpdateColumn("link_state", new)
	}

	return db.RowsAffected, db.Error
}

// BatchChUnlinkState batch change unlink state.
func (ConsensusOrder) BatchChUnlinkState(tx *gorm.DB, ids []uint, state model.ConsensusOrderUnlinkState) error {
	db := gormDb
	if tx != nil {
		db = tx
	}
	db = db.Model(new(ConsensusOrder)).
		Where("id in (?)", ids)

	if state == model.ConsensusOrderUnlinkStateConfirmed {
		db = db.Update(map[string]interface{}{
			"unlink_state": state,
			"unlinked_at":  time.Now(),
		})
	} else {
		db = db.UpdateColumn("unlink_state", state)
	}

	return db.Error
}

// LinkedListForUpdate list all confirm and not unlinked orders.
func (ConsensusOrder) LinkedListForUpdate(tx *gorm.DB, appID, UID1, UID2 string) ([]ConsensusOrder, error) {
	rs := make([]ConsensusOrder, 0)
	err := tx.Model(new(ConsensusOrder)).
		Where("app_id = ? and ((uid1 = ? and uid2 = ?) or (uid2 = ? and uid1 = ?)) and link_state <> ? and unlink_state = ?",
			appID, UID1, UID2, UID2, UID1, model.ConsensusOrderLinkStateWait, model.ConsensusOrderUnlinkStateNone).
		Set("gorm:query_option", "FOR UPDATE").
		Scan(&rs).Error
	return rs, err
}

// Set generate consensus order data model from mysql table data model.
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
