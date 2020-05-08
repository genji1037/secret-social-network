package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"time"
)

type ConsensusLinkLogAction string
type ConsensusLinkLogState string

const (
	ConsensusLinkLogActionLink   ConsensusLinkLogAction = "link"
	ConsensusLinkLogActionUnlink ConsensusLinkLogAction = "unlink"

	ConsensusLinkLogStateWait   ConsensusLinkLogState = "wait"
	ConsensusLinkLogStateCommit ConsensusLinkLogState = "commit"
	ConsensusLinkLogStateFailed ConsensusLinkLogState = "failed"
)

type ConsensusLinkLog struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	OpenID1   string                 `gorm:"not null"`
	OpenID2   string                 `gorm:"not null"`
	Action    ConsensusLinkLogAction `gorm:"not null"`
	Value     *decimal.Decimal       `sql:"type:decimal(32,16);not null"`
	AppID     string                 `gorm:"not null"`
	State     ConsensusLinkLogState  `gorm:"not null;default:'wait'"`
	// todo add failed reason
}

func (c *ConsensusLinkLog) Create() error {
	return gormDb.Create(c).Error
}

func (c *ConsensusLinkLog) CreateTx(tx *gorm.DB) error {
	return tx.Create(c).Error
}

func (c *ConsensusLinkLog) ChState(new ConsensusLinkLogState) error {
	return gormDb.Model(c).Where("id = ?", c.ID).UpdateColumn("state", new).Error
}
