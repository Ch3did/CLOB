package services

import (
	"context"
	"fmt"
	"math"

	"gorm.io/gorm"

	"CLOB/models"
)

type MatchingEngine struct {
	db      *gorm.DB
	balance *BalanceService
}

func NewMatchingEngine(db *gorm.DB) *MatchingEngine {
	return &MatchingEngine{
		db:      db,
		balance: NewBalanceService(db),
	}
}

func (e *MatchingEngine) MatchOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	var match models.Order

	var oppositeSide string
	var priceOp string
	var orderClause string

	if order.Side == models.SideBuy {
		oppositeSide = string(models.SideSell)
		priceOp = "<="
		orderClause = "price ASC, created_at ASC"
	} else {
		oppositeSide = string(models.SideBuy)
		priceOp = ">="
		orderClause = "price DESC, created_at ASC"
	}

	q := e.db.WithContext(ctx).
		Where("instrument = ? AND side = ? AND id != ?", order.Instrument, oppositeSide, order.ID).
		Where(fmt.Sprintf("price %s ?", priceOp), order.Price).
		Order(orderClause)

	if err := q.First(&match).Error; err != nil {
		return nil, err
	}

	return &match, nil
}

func (e *MatchingEngine) CreateTrade(ctx context.Context, incoming *models.Order, matched *models.Order) error {

	var buyOrder, sellOrder *models.Order
	if incoming.Side == models.SideBuy {
		buyOrder = incoming
		sellOrder = matched
	} else {
		buyOrder = matched
		sellOrder = incoming
	}

	price := matched.Price

	tradeQty := math.Min(buyOrder.Quantity, sellOrder.Quantity)

	return e.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		trade := &models.Trade{
			Instrument:  buyOrder.Instrument,
			BuyOrderID:  buyOrder.ID,
			SellOrderID: sellOrder.ID,
			Price:       price,
			Quantity:    tradeQty,
		}

		if err := tx.Create(trade).Error; err != nil {
			return err
		}

		base, quote, err := splitInstrument(buyOrder.Instrument)
		if err != nil {
			return err
		}

		totalQuote := tradeQty * price

		if err := e.balance.LockToTrade(ctx, buyOrder.AccountID, quote, totalQuote); err != nil {
			return err
		}

		if err := e.balance.Credit(ctx, buyOrder.AccountID, base, tradeQty); err != nil {
			return err
		}

		if err := e.balance.LockToTrade(ctx, sellOrder.AccountID, base, tradeQty); err != nil {
			return err
		}

		if err := e.balance.Credit(ctx, sellOrder.AccountID, quote, totalQuote); err != nil {
			return err
		}

		// calcular remaining
		buyRemaining := buyOrder.Quantity - tradeQty
		sellRemaining := sellOrder.Quantity - tradeQty

		if buyRemaining <= 0 {
			if err := tx.Model(buyOrder).Update("status", "filled").Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(buyOrder).Updates(map[string]interface{}{
				"quantity": buyRemaining,
				"status":   "partial",
			}).Error; err != nil {
				return err
			}
		}

		if sellRemaining <= 0 {
			if err := tx.Model(sellOrder).Update("status", "filled").Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(sellOrder).Updates(map[string]interface{}{
				"quantity": sellRemaining,
				"status":   "partial",
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
