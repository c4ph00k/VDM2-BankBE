package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"VDM2-BankBE/internal/model"
	"VDM2-BankBE/internal/util"
)

// GormOAuthTokenRepository implements OAuthTokenRepository using GORM
type GormOAuthTokenRepository struct {
	db *gorm.DB
}

// NewGormOAuthTokenRepository creates a new OAuth token repository with GORM
func NewGormOAuthTokenRepository(db *gorm.DB) OAuthTokenRepository {
	return &GormOAuthTokenRepository{db: db}
}

// Create inserts a new OAuth token into the database
func (r *GormOAuthTokenRepository) Create(ctx context.Context, token *model.OAuthToken) error {
	err := r.db.WithContext(ctx).Create(token).Error
	if err != nil {
		return errors.Wrap(err, "failed to create OAuth token")
	}

	return nil
}

// GetByUserID retrieves an OAuth token by user ID
func (r *GormOAuthTokenRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.OAuthToken, error) {
	var token model.OAuthToken

	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.NewNotFoundError("OAuth token not found")
		}
		return nil, errors.Wrap(err, "failed to get OAuth token by user ID")
	}

	return &token, nil
}

// Update updates an OAuth token in the database
func (r *GormOAuthTokenRepository) Update(ctx context.Context, token *model.OAuthToken) error {
	err := r.db.WithContext(ctx).Save(token).Error
	if err != nil {
		return errors.Wrap(err, "failed to update OAuth token")
	}

	return nil
}

// Delete deletes an OAuth token from the database
func (r *GormOAuthTokenRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&model.OAuthToken{}, "user_id = ?", userID).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete OAuth token")
	}

	return nil
}
