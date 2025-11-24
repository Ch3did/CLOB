package services

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"CLOB/models"
)

type BalanceService struct {
	db *gorm.DB
}

func NewBalanceService(db *gorm.DB) *BalanceService {
	return &BalanceService{db: db}
}

func (s *BalanceService) LockToTrade(
	ctx context.Context,
	accountID uint,
	asset string,
	amount float64,
) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		bal, err := s.getOrCreate(ctx, tx, accountID, asset)
		if err != nil {
			return err
		}

		if bal.Locked < amount {
			return fmt.Errorf("insufficient locked balance for account %d, asset %s", accountID, asset)
		}

		bal.Locked -= amount

		return tx.Save(bal).Error
	})
}

func (s *BalanceService) getOrCreate(
	ctx context.Context,
	tx *gorm.DB,
	accountID uint,
	asset string,
) (*models.Balance, error) {
	var bal models.Balance

	err := tx.WithContext(ctx).
		Where("account_id = ? AND asset = ?", accountID, asset).
		First(&bal).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		bal = models.Balance{
			AccountID: accountID,
			Asset:     asset,
			Available: 0,
			Locked:    0,
		}
		if err := tx.Create(&bal).Error; err != nil {
			return nil, err
		}
		return &bal, nil
	}
	if err != nil {
		return nil, err
	}

	return &bal, nil
}

func (s *BalanceService) Credit(
	ctx context.Context,
	accountID uint,
	asset string,
	amount float64,
) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		bal, err := s.getOrCreate(ctx, tx, accountID, asset)
		if err != nil {
			return err
		}

		bal.Available += amount
		return tx.Save(bal).Error
	})
}

func (s *BalanceService) Debit(
	ctx context.Context,
	accountID uint,
	asset string,
	amount float64,
) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		bal, err := s.getOrCreate(ctx, tx, accountID, asset)
		if err != nil {
			return err
		}

		if bal.Available < amount {
			return fmt.Errorf("insufficient available balance for account %d, asset %s", accountID, asset)
		}

		bal.Available -= amount
		return tx.Save(bal).Error
	})
}

func (s *BalanceService) ListByAccount(
	ctx context.Context,
	accountID uint,
) ([]models.Balance, error) {
	var balances []models.Balance
	err := s.db.WithContext(ctx).
		Where("account_id = ?", accountID).
		Find(&balances).Error
	return balances, err
}

func (s *BalanceService) GetOne(
	ctx context.Context,
	accountID uint,
	asset string,
) (*models.Balance, error) {
	var bal models.Balance
	err := s.db.WithContext(ctx).
		Where("account_id = ? AND asset = ?", accountID, asset).
		First(&bal).Error
	if err != nil {
		return nil, err
	}
	return &bal, nil
}
