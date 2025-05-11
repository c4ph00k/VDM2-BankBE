package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"VDM2-BankBE/internal/model"
	"VDM2-BankBE/internal/util"
)

// GormMovementRepository implements MovementRepository using GORM
type GormMovementRepository struct {
	db *gorm.DB
}

// NewGormMovementRepository creates a new movement repository with GORM
func NewGormMovementRepository(db *gorm.DB) MovementRepository {
	return &GormMovementRepository{db: db}
}

// Create inserts a new movement into the database
func (r *GormMovementRepository) Create(ctx context.Context, movement *model.Movement) error {
	err := r.db.WithContext(ctx).Create(movement).Error
	if err != nil {
		return errors.Wrap(err, "failed to create movement")
	}

	return nil
}

// GetByID retrieves a movement by ID
func (r *GormMovementRepository) GetByID(ctx context.Context, id uint64) (*model.Movement, error) {
	var movement model.Movement

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&movement).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.NewNotFoundError("movement not found")
		}
		return nil, errors.Wrap(err, "failed to get movement by ID")
	}

	return &movement, nil
}

// GetByAccountID retrieves movements for an account with pagination
func (r *GormMovementRepository) GetByAccountID(
	ctx context.Context,
	accountID uuid.UUID,
	params *util.PaginationParams,
) ([]*model.Movement, int, error) {
	var movements []*model.Movement
	var count int64

	// Count total records
	err := r.db.WithContext(ctx).
		Model(&model.Movement{}).
		Where("account_id = ?", accountID).
		Count(&count).Error
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to count movements")
	}

	// Get paginated records
	err = r.db.WithContext(ctx).
		Where("account_id = ?", accountID).
		Order("occurred_at DESC").
		Offset(params.Offset()).
		Limit(params.Limit).
		Find(&movements).Error
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get movements by account ID")
	}

	return movements, int(count), nil
}
