package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"VDM2-BankBE/internal/model"
	"VDM2-BankBE/internal/util"
)

// AuthService defines methods for authentication
type AuthService interface {
	SignUp(ctx context.Context, email, username, firstName, lastName, fiscalCode, password string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	GoogleAuth(ctx context.Context) (string, string, error)
	GoogleCallback(ctx context.Context, code, state string) (string, error)
	VerifyToken(ctx context.Context, token string) (*model.User, error)
}

// AccountService defines methods for account operations
type AccountService interface {
	Create(ctx context.Context, userID uuid.UUID) (*model.Account, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Account, error)
	GetBalance(ctx context.Context, accountID uuid.UUID) (decimal.Decimal, error)
}

// MovementService defines methods for movement operations
type MovementService interface {
	Create(ctx context.Context, accountID uuid.UUID, amount decimal.Decimal, movementType, description string) (*model.Movement, error)
	GetByID(ctx context.Context, id uint64) (*model.Movement, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID, page, limit int) (*util.PaginatedResponse, error)
}

// TransferService defines methods for transfer operations
type TransferService interface {
	Transfer(ctx context.Context, fromAccountID, toAccountID uuid.UUID, amount decimal.Decimal, description string) (*model.Transfer, error)
	GetByID(ctx context.Context, id uint64) (*model.Transfer, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID, page, limit int) (*util.PaginatedResponse, error)
}

// Service combines all services
type Service struct {
	Auth     AuthService
	Account  AccountService
	Movement MovementService
	Transfer TransferService
}

// NewService creates a new service provider
func NewService(
	authService AuthService,
	accountService AccountService,
	movementService MovementService,
	transferService TransferService,
) *Service {
	return &Service{
		Auth:     authService,
		Account:  accountService,
		Movement: movementService,
		Transfer: transferService,
	}
}
