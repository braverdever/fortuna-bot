package models

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID             uuid.UUID `gorm:"id"`
	ChatId         int       `gorm:"chat_id"`
	ChainScanLabel string    `gorm:"chain_scan_label"`
	AccountWorth   int       `gorm:"account_worth"`
	PrivateKey     string    `gorm:"private_key"`
	Address        string    `gorm:"wallet_address"`
	Createdate     time.Time `gorm:"column:create_date;type:timestamp"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp"`
}
