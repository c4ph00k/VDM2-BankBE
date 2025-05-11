package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"VDM2-BankBE/internal/model"
	"VDM2-BankBE/internal/util"
)

// GormUserRepository implements UserRepository using GORM
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new user repository with GORM
func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

// Create inserts a new user into the database
func (r *GormUserRepository) Create(ctx context.Context, user *model.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *GormUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.NewNotFoundError("user not found")
		}
		return nil, errors.Wrap(err, "failed to get user by ID")
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.NewNotFoundError("user not found")
		}
		return nil, errors.Wrap(err, "failed to get user by email")
	}

	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *GormUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.NewNotFoundError("user not found")
		}
		return nil, errors.Wrap(err, "failed to get user by username")
	}

	return &user, nil
}

// Update updates a user in the database
func (r *GormUserRepository) Update(ctx context.Context, user *model.User) error {
	err := r.db.WithContext(ctx).Save(user).Error
	if err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	return nil
}

// Delete deletes a user from the database
func (r *GormUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	return nil
}
