package service

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"VDM2-BankBE/internal/model"
	"VDM2-BankBE/internal/repository"
	"VDM2-BankBE/internal/util"
	"VDM2-BankBE/pkg/cache"
)

// DefaultMovementService implements MovementService
type DefaultMovementService struct {
	movementRepo repository.MovementRepository
	accountRepo  repository.AccountRepository
	redisClient  *cache.RedisClient
}

// NewMovementService creates a new movement service
func NewMovementService(
	movementRepo repository.MovementRepository,
	accountRepo repository.AccountRepository,
	redisClient *cache.RedisClient,
) MovementService {
	return &DefaultMovementService{
		movementRepo: movementRepo,
		accountRepo:  accountRepo,
		redisClient:  redisClient,
	}
}

// Create creates a new movement and updates the account balance
func (s *DefaultMovementService) Create(
	ctx context.Context,
	accountID uuid.UUID,
	amount decimal.Decimal,
	movementType string,
	description string,
) (*model.Movement, error) {
	// Validate movement type
	if movementType != "credit" && movementType != "debit" {
		return nil, util.NewBadRequestError("movement type must be 'credit' or 'debit'")
	}

	// Get the account to verify it exists
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	// Create the movement
	movement := &model.Movement{
		AccountID:   accountID,
		Amount:      amount,
		Type:        movementType,
		Description: description,
		OccurredAt:  time.Now(),
	}

	// Update the account balance
	balanceChange := amount
	if movementType == "debit" {
		balanceChange = amount.Neg()
	}

	// Update balance in DB
	err = s.accountRepo.UpdateBalance(ctx, accountID, balanceChange)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update account balance")
	}

	// Create movement in DB
	err = s.movementRepo.Create(ctx, movement)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create movement")
	}

	// Update balance cache
	newBalance := account.Balance.Add(balanceChange)
	_ = s.redisClient.SetBalanceCache(ctx, accountID, newBalance)

	return movement, nil
}

// GetByID retrieves a movement by ID
func (s *DefaultMovementService) GetByID(ctx context.Context, id uint64) (*model.Movement, error) {
	movement, err := s.movementRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get movement")
	}

	return movement, nil
}

// GetByAccountID retrieves movements for an account with pagination
func (s *DefaultMovementService) GetByAccountID(ctx context.Context, accountID uuid.UUID, page, limit int) (*util.PaginatedResponse, error) {
	// Check if account exists
	_, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	// Create pagination params
	params, err := util.NewPaginationParams(strconv.Itoa(page), strconv.Itoa(limit))
	if err != nil {
		return nil, errors.Wrap(err, "invalid pagination parameters")
	}

	// Get movements
	movements, count, err := s.movementRepo.GetByAccountID(ctx, accountID, params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get movements")
	}

	// Create paginated response
	response := util.NewPaginatedResponse(movements, params, count)
	return response, nil
}
