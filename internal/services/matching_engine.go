package services

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"CLOB/models"
)

type MatchingEngine struct {
	db *gorm.DB
}

func NewMatchingEngine(db *gorm.DB) *MatchingEngine {
	return &MatchingEngine{db: db}
}

// MatchOrder recebe uma ordem já criada e tenta achar uma contraparte compatível.
// Se não encontrar, retorna gorm.ErrRecordNotFound.
func (e *MatchingEngine) MatchOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	var match models.Order

	var oppositeSide string
	var priceOp string
	var orderClause string

	if order.Side == models.SideBuy {
		// BUY casa com SELL onde sell.price <= buy.price
		oppositeSide = string(models.SideSell)
		priceOp = "<="
		orderClause = "price ASC, created_at ASC"
	} else {
		// SELL casa com BUY onde buy.price >= sell.price
		oppositeSide = string(models.SideBuy)
		priceOp = ">="
		orderClause = "price DESC, created_at ASC"
	}

	query := e.db.WithContext(ctx).
		Where("instrument = ? AND side = ? AND id != ?", order.Instrument, oppositeSide, order.ID).
		Where(fmt.Sprintf("price %s ?", priceOp), order.Price).
		Order(orderClause)

	if err := query.First(&match).Error; err != nil {
		return nil, err
	}

	return &match, nil
}
