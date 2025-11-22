package services

// import (
// 	"context"
// 	"fmt"
// 	"strings"

// 	"gorm.io/gorm"

// 	"CLOB/models"
// )

// type MatchingEngine struct {
// 	db      *gorm.DB
// 	balance *BalanceService
// }

// func NewMatchingEngine(db *gorm.DB) *MatchingEngine {
// 	return &MatchingEngine{
// 		db:      db,
// 		balance: NewBalanceService(db),
// 	}
// }

// // BTC/BRL ou BTC-BRL → base=BTC, quote=BRL
// func splitInstrument(instr string) (base, quote string, err error) {
// 	parts := strings.Split(instr, "/")
// 	if len(parts) != 2 {
// 		parts = strings.Split(instr, "-")
// 	}
// 	if len(parts) != 2 {
// 		return "", "", fmt.Errorf("invalid instrument: %s", instr)
// 	}
// 	return parts[0], parts[1], nil
// }

// // PlaceOrder insere a ordem no livro, tenta casar e gera Trades + atualização de Balance.
// func (e *MatchingEngine) PlaceOrder(
// 	ctx context.Context,
// 	order *models.Order,
// ) ([]models.Trade, error) {
// 	var trades []models.Trade

// 	err := e.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		// 1. Persiste a ordem nova
// 		if err := tx.Create(order).Error; err != nil {
// 			return err
// 		}

// 		base, quote, err := splitInstrument(order.Instrument)
// 		if err != nil {
// 			return err
// 		}

// 		// 2. Buscar ordens opostas compatíveis (contra-parte)
// 		var counterparts []models.Order

// 		switch order.Side {
// 		case models.SideBuy:
// 			// Buy casa com sells com preço <= buy.Price
// 			if err := tx.
// 				Where("instrument = ? AND side = ? AND price <= ? AND quantity > 0",
// 					order.Instrument, models.SideSell, order.Price).
// 				Order("price ASC, created_at ASC").
// 				Find(&counterparts).Error; err != nil {
// 				return err
// 			}

// 		case models.SideSell:
// 			// Sell casa com buys com preço >= sell.Price
// 			if err := tx.
// 				Where("instrument = ? AND side = ? AND price >= ? AND quantity > 0",
// 					order.Instrument, models.SideBuy, order.Price).
// 				Order("price DESC, created_at ASC").
// 				Find(&counterparts).Error; err != nil {
// 				return err
// 			}

// 		default:
// 			return fmt.Errorf("invalid side")
// 		}

// 		remaining := order.Quantity

// 		for i := range counterparts {
// 			if remaining <= 0 {
// 				break
// 			}

// 			resting := &counterparts[i]
// 			if resting.ID == order.ID {
// 				continue
// 			}

// 			availableQty := resting.Quantity
// 			if availableQty <= 0 {
// 				continue
// 			}

// 			execQty := remaining
// 			if execQty > availableQty {
// 				execQty = availableQty
// 			}

// 			// Regra simples: trade sai no preço da ordem "resting" (maker)
// 			execPrice := resting.Price

// 			var buyOrder, sellOrder *models.Order
// 			if order.Side == models.SideBuy {
// 				buyOrder, sellOrder = order, resting
// 			} else {
// 				buyOrder, sellOrder = resting, order
// 			}

// 			trade, err := e.executeMatch(
// 				ctx,
// 				tx,
// 				buyOrder,
// 				sellOrder,
// 				execQty,
// 				execPrice,
// 				base,
// 				quote,
// 			)
// 			if err != nil {
// 				return err
// 			}

// 			trades = append(trades, *trade)

// 			// Atualiza as quantities (estamos usando Quantity como "restante")
// 			remaining -= execQty
// 			resting.Quantity -= execQty

// 			if err := tx.Save(resting).Error; err != nil {
// 				return err
// 			}
// 		}

// 		// Atualiza a ordem de entrada com a quantidade restante
// 		order.Quantity = remaining
// 		if err := tx.Save(order).Error; err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	return trades, err
// }

// // executeMatch: cria Trade e atualiza Balance de comprador e vendedor.
// func (e *MatchingEngine) executeMatch(
// 	ctx context.Context,
// 	tx *gorm.DB,
// 	buyOrder *models.Order,
// 	sellOrder *models.Order,
// 	qty float64,
// 	price float64,
// 	base string,
// 	quote string,
// ) (*models.Trade, error) {
// 	// 1. Criar Trade
// 	trade := &models.Trade{
// 		Instrument:  buyOrder.Instrument,
// 		BuyOrderID:  buyOrder.ID,
// 		SellOrderID: sellOrder.ID,
// 		Price:       price,
// 		Quantity:    qty,
// 	}

// 	if err := tx.WithContext(ctx).Create(trade).Error; err != nil {
// 		return nil, err
// 	}

// 	total := qty * price

// 	// 2. Atualizar balances segundo o exemplo do enunciado:
// 	//    comprador recebe base, paga quote
// 	//    vendedor entrega base, recebe quote

// 	// Comprador: +base, -quote
// 							  (ctx , accountID uint, asset string, amount float64
// 	if err := e.balance.Credit(ctx, buyOrder.AccountID, base, qty, tx); err != nil {
// 		return nil, err
// 	}
// 	if err := e.balance.Debit(ctx, buyOrder.AccountID, quote, total, tx); err != nil {
// 		return nil, err
// 	}

// 	// Vendedor: -base, +quote
// 	if err := e.balance.Debit(ctx, sellOrder.AccountID, base, qty, tx); err != nil {
// 		return nil, err
// 	}
// 	if err := e.balance.Credit(ctx, sellOrder.AccountID, quote, total, tx); err != nil {
// 		return nil, err
// 	}

// 	return trade, nil
// }
