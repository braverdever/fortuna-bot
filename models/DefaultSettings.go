package models

import (
	"time"

	"github.com/google/uuid"
)

type DefaultSettings struct {
	ID              uuid.UUID `gorm:"id"`
	Slippage        int       `gorm:"slippage"`
	SellGweiExtra   float32   `gorm:"sell_gwei_extra"`
	ApproveGwei     float32   `gorm:"approve_gwei"`
	BuyTax          float32   `gorm:"buy_tax"`
	SellTax         float32   `gorm:"sell_tax"`
	MinLiquidity    int       `gorm:"min_liquidity"`
	AlphaMode       bool      `gorm:"alpha_mode"`
	MultitxOrRevert bool      `gorm:"multitx_or_revert"`
	AntiRug         bool      `gorm:"anti_rug"`
	Createdate      time.Time `gorm:"column:create_date;type:timestamp"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:timestamp"`
}
