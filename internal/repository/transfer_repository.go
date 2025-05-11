package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"VDM2-BankBE/internal/model"
	"VDM2-BankBE/internal/util"
)

// GormTransferRepository implements TransferRepository using GORM
type GormTransferRepository struct {
	db *gorm.DB
}

// NewGormTransferRepository creates a new transfer repository with GORM
func NewGormTransferRepository(db *gorm.DB) TransferRepository {
	return &GormTransferRepository{db: db}
}

// Create inserts a new transfer into the database
func (r *GormTransferRepository) Create(ctx context.Context, transfer *model.Transfer) error {
	err := r.db.WithContext(ctx).Create(transfer).Error
	if err != nil {
		return errors.Wrap(err, "failed to create transfer")
	}

	return nil
}

// GetByID retrieves a transfer by ID
func (r *GormTransferRepository) GetByID(ctx context.Context, id uint64) (*model.Transfer, error) {
	var transfer model.Transfer

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&transfer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.NewNotFoundError("transfer not found")
		}
		return nil, errors.Wrap(err, "failed to get transfer by ID")
	}

	return &transfer, nil
}

// GetByAccountID retrieves transfers for an account with pagination
func (r *GormTransferRepository) GetByAccountID(
	ctx context.Context,
	accountID uuid.UUID,
	params *util.PaginationParams,
) ([]*model.Transfer, int, error) {
	var transfers []*model.Transfer
	var count int64

	// Count total records
	err := r.db.WithContext(ctx).
		Model(&model.Transfer{}).
		Where("from_account = ? OR to_account = ?", accountID, accountID).
		Count(&count).Error
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to count transfers")
	}

	// Get paginated records
	err = r.db.WithContext(ctx).
		Where("from_account = ? OR to_account = ?", accountID, accountID).
		Order("initiated_at DESC").
		Offset(params.Offset()).
		Limit(params.Limit).
		Find(&transfers).Error
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get transfers by account ID")
	}

	return transfers, int(count), nil
}

// UpdateStatus updates a transfer's status
func (r *GormTransferRepository) UpdateStatus(
	ctx context.Context,
	id uint64,
	status string,
	completedAtStr *string,
) error {
	updates := map[string]interface{}{
		"status": status,
	}

	// If completedAt is provided, parse and add it to the updates
	if completedAtStr != nil {
		completedAt, err := time.Parse(time.RFC3339, *completedAtStr)
		if err != nil {
			return errors.Wrap(err, "invalid completedAt time format")
		}
		updates["completed_at"] = completedAt
	}

	err := r.db.WithContext(ctx).
		Model(&model.Transfer{}).
		Where("id = ?", id).
		Updates(updates).Error
	if err != nil {
		return errors.Wrap(err, "failed to update transfer status")
	}

	return nil
}
