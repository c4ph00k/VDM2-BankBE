package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"VDM2-BankBE/internal/model"
	"VDM2-BankBE/internal/repository"
)

// DefaultAccountService implements AccountService
type DefaultAccountService struct {
	accountRepo repository.AccountRepository
	redisClient CacheClient
}

// NewAccountService creates a new account service
func NewAccountService(
	accountRepo repository.AccountRepository,
	redisClient CacheClient,
) AccountService {
	return &DefaultAccountService{
		accountRepo: accountRepo,
		redisClient: redisClient,
	}
}

// Create creates a new account for a user
func (s *DefaultAccountService) Create(ctx context.Context, userID uuid.UUID) (*model.Account, error) {
	account := &model.Account{
		ID:        uuid.New(),
		UserID:    userID,
		Balance:   decimal.NewFromInt(0),
		Currency:  "EUR",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.accountRepo.Create(ctx, account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create account")
	}

	return account, nil
}

// GetByID retrieves an account by ID
func (s *DefaultAccountService) GetByID(ctx context.Context, id uuid.UUID) (*model.Account, error) {
	account, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	return account, nil
}

// GetByUserID retrieves an account by user ID
func (s *DefaultAccountService) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Account, error) {
	account, err := s.accountRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account by user ID")
	}

	return account, nil
}

// GetBalance retrieves an account's balance, using Redis cache when available
func (s *DefaultAccountService) GetBalance(ctx context.Context, accountID uuid.UUID) (decimal.Decimal, error) {
	// Try to get balance from cache first
	balance, err := s.redisClient.GetBalanceCache(ctx, accountID)
	if err == nil {
		// Cache hit
		return balance, nil
	}

	// Cache miss, get from DB
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "failed to get account balance")
	}

	// Update cache
	if err := s.redisClient.SetBalanceCache(ctx, accountID, account.Balance); err != nil {
		// Just log the error, but don't fail the request
		// In a real app, you'd use a logger here
		_ = err
	}

	return account.Balance, nil
}
