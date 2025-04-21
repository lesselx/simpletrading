package memory

import (
	"simpletrading/tradeservice/internal/domain"

	"gorm.io/gorm"
)

type TradeRepository struct {
	db *gorm.DB
}

func NewTradeRepository(db *gorm.DB) *TradeRepository {
	return &TradeRepository{db: db}
}

func (r *TradeRepository) Insert(trade *domain.Trade) error {
	return r.db.Create(trade).Error
}

// func (r *TradeRepository) GetLowestPriceInLast24Hours() (float64, error) {
// 	var trade domain.Trade
// 	threshold := time.Now().Add(-24 * time.Hour)

// 	err := r.db.
// 		Where("created_at >= ?", threshold).
// 		Order("price asc").
// 		First(&trade).Error

// 	if err != nil {
// 		return 0, err
// 	}

// 	return trade.Price, nil
// }
