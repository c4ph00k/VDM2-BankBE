package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"VDM2-BankBE/internal/model"
	"VDM2-BankBE/internal/util"
)

// GormAccountRepository implements AccountRepository using GORM
type GormAccountRepository struct {
	db *gorm.DB
}

// NewGormAccountRepository creates a new account repository with GORM
func NewGormAccountRepository(db *gorm.DB) AccountRepository {
	return &GormAccountRepository{db: db}
}

// Create inserts a new account into the database
func (r *GormAccountRepository) Create(ctx context.Context, account *model.Account) error {
	if account.ID == uuid.Nil {
		account.ID = uuid.New()
	}

	err := r.db.WithContext(ctx).Create(account).Error
	if err != nil {
		return errors.Wrap(err, "failed to create account")
	}

	return nil
}

// GetByID retrieves an account by ID
func (r *GormAccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Account, error) {
	var account model.Account

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.NewNotFoundError("account not found")
		}
		return nil, errors.Wrap(err, "failed to get account by ID")
	}

	return &account, nil
}

// GetByUserID retrieves an account by user ID
func (r *GormAccountRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Account, error) {
	var account model.Account

	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.NewNotFoundError("account not found")
		}
		return nil, errors.Wrap(err, "failed to get account by user ID")
	}

	return &account, nil
}

// UpdateBalance updates an account's balance
func (r *GormAccountRepository) UpdateBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal) error {
	// Use a transaction to ensure consistency
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "failed to begin transaction")
	}

	// Get the current account with locking to prevent race conditions
	var account model.Account
	err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", id).First(&account).Error
	if err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return util.NewNotFoundError("account not found")
		}
		return errors.Wrap(err, "failed to get account for balance update")
	}

	// Update the balance
	account.Balance = account.Balance.Add(amount)

	// Check for negative balance
	if account.Balance.LessThan(decimal.Zero) {
		tx.Rollback()
		return util.NewBadRequestError("insufficient funds")
	}

	// Save the updated balance
	err = tx.Save(&account).Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to update account balance")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// Delete deletes an account from the database
func (r *GormAccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&model.Account{}, "id = ?", id).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete account")
	}

	return nil
}
