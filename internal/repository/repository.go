package repository

import (
	"context"

	"VDM2-BankBE/internal/model"
	"VDM2-BankBE/internal/util"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// AccountRepository defines the interface for account repository operations
type AccountRepository interface {
	Create(ctx context.Context, account *model.Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Account, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// MovementRepository defines the interface for movement repository operations
type MovementRepository interface {
	Create(ctx context.Context, movement *model.Movement) error
	GetByID(ctx context.Context, id uint64) (*model.Movement, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID, params *util.PaginationParams) ([]*model.Movement, int, error)
}

// OAuthTokenRepository defines the interface for OAuth token repository operations
type OAuthTokenRepository interface {
	Create(ctx context.Context, token *model.OAuthToken) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.OAuthToken, error)
	Update(ctx context.Context, token *model.OAuthToken) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

// TransferRepository defines the interface for transfer repository operations
type TransferRepository interface {
	Create(ctx context.Context, transfer *model.Transfer) error
	GetByID(ctx context.Context, id uint64) (*model.Transfer, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID, params *util.PaginationParams) ([]*model.Transfer, int, error)
	UpdateStatus(ctx context.Context, id uint64, status string, completedAt *string) error
}

// Repository provides access to all repositories
type Repository struct {
	User       UserRepository
	Account    AccountRepository
	Movement   MovementRepository
	OAuthToken OAuthTokenRepository
	Transfer   TransferRepository
}

// NewRepository creates a new repository provider
func NewRepository(
	userRepo UserRepository,
	accountRepo AccountRepository,
	movementRepo MovementRepository,
	oauthTokenRepo OAuthTokenRepository,
	transferRepo TransferRepository,
) *Repository {
	return &Repository{
		User:       userRepo,
		Account:    accountRepo,
		Movement:   movementRepo,
		OAuthToken: oauthTokenRepo,
		Transfer:   transferRepo,
	}
}
