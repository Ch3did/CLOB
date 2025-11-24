package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"CLOB/models"
)

type OrderService struct {
	db     *gorm.DB
	engine *MatchingEngine
}

func NewOrderService(db *gorm.DB, engine *MatchingEngine) *OrderService {
	return &OrderService{
		db:     db,
		engine: engine,
	}
}

func splitInstrument(instr string) (string, string, error) {
	parts := strings.Split(instr, "/")
	if len(parts) != 2 {
		parts = strings.Split(instr, "-")
	}
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid instrument: %s", instr)
	}
	return parts[0], parts[1], nil
}

// CreateOrder: valida side, trava saldo e cria ordem, depois tenta casar.
func (s *OrderService) CreateOrder(
	ctx context.Context,
	accountID uint,
	instrument string,
	side models.Side,
	price float64,
	quantity float64,
) (*models.Order, *models.Order, error) {
	if side != models.SideBuy && side != models.SideSell {
		return nil, nil, fmt.Errorf("side must be BUY or SELL")
	}

	base, quote, err := splitInstrument(instrument)
	if err != nil {
		return nil, nil, err
	}

	var order models.Order

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var asset string
		var required float64

		if side == models.SideBuy {
			asset = quote
			required = price * quantity
		} else {
			asset = base
			required = quantity
		}

		var bal models.Balance
		if err := tx.
			Where("account_id = ? AND asset = ?", accountID, asset).
			First(&bal).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("insufficient balance")
			}
			return err
		}

		if bal.Available < required {
			return fmt.Errorf("insufficient balance")
		}

		bal.Available -= required
		bal.Locked += required
		if err := tx.Save(&bal).Error; err != nil {
			return err
		}

		order = models.Order{
			AccountID:  accountID,
			Instrument: instrument,
			Side:       side,
			Price:      price,
			Quantity:   quantity,
			Status:     "open",
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, nil, err
	}

	// tenta casar
	match, err := s.engine.MatchOrder(ctx, &order)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &order, nil, nil
		}
		return &order, nil, err
	}

	if err := s.engine.CreateTrade(ctx, &order, match); err != nil {
		return &order, match, err
	}

	return &order, match, nil
}

func (s *OrderService) CancelOrder(
	ctx context.Context,
	accountID uint,
	orderID uint,
) error {
	var order models.Order
	if err := s.db.WithContext(ctx).
		Where("id = ? AND account_id = ?", orderID, accountID).
		First(&order).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("order not found")
		}
		return err
	}

	var tradeCount int64
	if err := s.db.WithContext(ctx).
		Model(&models.Trade{}).
		Where("buy_order_id = ? OR sell_order_id = ?", order.ID, order.ID).
		Count(&tradeCount).Error; err != nil {

		return err
	}

	if tradeCount > 0 {
		return fmt.Errorf("order already has executed trades and cannot be canceled")
	}

	base, quote, err := splitInstrument(order.Instrument)
	if err != nil {
		return err
	}

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var asset string
		var amount float64

		if order.Side == models.SideBuy {
			asset = quote
			amount = order.Price * order.Quantity
		} else if order.Side == models.SideSell {
			asset = base
			amount = order.Quantity
		} else {
			return fmt.Errorf("invalid side on order")
		}

		var bal models.Balance
		if err := tx.
			Where("account_id = ? AND asset = ?", accountID, asset).
			First(&bal).Error; err != nil {

			return err
		}

		if bal.Locked < amount {
			return fmt.Errorf("locked balance too low for cancel")
		}

		bal.Locked -= amount
		bal.Available += amount
		if err := tx.Save(&bal).Error; err != nil {
			return err
		}

		if err := tx.Delete(&order).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
