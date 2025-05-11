package service

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"VDM2-BankBE/internal/model"
	"VDM2-BankBE/internal/repository"
	"VDM2-BankBE/internal/util"
	"VDM2-BankBE/pkg/cache"
)

// DefaultTransferService implements TransferService
type DefaultTransferService struct {
	transferRepo repository.TransferRepository
	accountRepo  repository.AccountRepository
	movementRepo repository.MovementRepository
	redisClient  *cache.RedisClient
	db           *gorm.DB // For transactions
}

// NewTransferService creates a new transfer service
func NewTransferService(
	transferRepo repository.TransferRepository,
	accountRepo repository.AccountRepository,
	movementRepo repository.MovementRepository,
	redisClient *cache.RedisClient,
	db *gorm.DB,
) TransferService {
	return &DefaultTransferService{
		transferRepo: transferRepo,
		accountRepo:  accountRepo,
		movementRepo: movementRepo,
		redisClient:  redisClient,
		db:           db,
	}
}

// Transfer performs a funds transfer between accounts
func (s *DefaultTransferService) Transfer(
	ctx context.Context,
	fromAccountID, toAccountID uuid.UUID,
	amount decimal.Decimal,
	description string,
) (*model.Transfer, error) {
	// Validate amount
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, util.NewBadRequestError("amount must be greater than zero")
	}

	// Check if accounts exist
	fromAccount, err := s.accountRepo.GetByID(ctx, fromAccountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get source account")
	}

	toAccount, err := s.accountRepo.GetByID(ctx, toAccountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get destination account")
	}

	// Check for self-transfer
	if fromAccountID == toAccountID {
		return nil, util.NewBadRequestError("cannot transfer to the same account")
	}

	// Check if source account has sufficient funds
	if fromAccount.Balance.LessThan(amount) {
		return nil, util.NewBadRequestError("insufficient funds")
	}

	// Create transfer record
	transfer := &model.Transfer{
		FromAccount: fromAccountID,
		ToAccount:   toAccountID,
		Amount:      amount,
		Status:      "pending",
		InitiatedAt: time.Now(),
	}

	// Execute transfer in a transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Create transfer record
		if err := s.transferRepo.Create(ctx, transfer); err != nil {
			return errors.Wrap(err, "failed to create transfer record")
		}

		// Create debit movement
		debitMovement := &model.Movement{
			AccountID:   fromAccountID,
			Amount:      amount,
			Type:        "debit",
			Description: description + " (Transfer #" + uintToString(transfer.ID) + ")",
			OccurredAt:  time.Now(),
		}
		if err := s.movementRepo.Create(ctx, debitMovement); err != nil {
			return errors.Wrap(err, "failed to create debit movement")
		}

		// Create credit movement
		creditMovement := &model.Movement{
			AccountID:   toAccountID,
			Amount:      amount,
			Type:        "credit",
			Description: description + " (Transfer #" + uintToString(transfer.ID) + ")",
			OccurredAt:  time.Now(),
		}
		if err := s.movementRepo.Create(ctx, creditMovement); err != nil {
			return errors.Wrap(err, "failed to create credit movement")
		}

		// Update source account balance
		if err := s.accountRepo.UpdateBalance(ctx, fromAccountID, amount.Neg()); err != nil {
			return errors.Wrap(err, "failed to update source account balance")
		}

		// Update destination account balance
		if err := s.accountRepo.UpdateBalance(ctx, toAccountID, amount); err != nil {
			return errors.Wrap(err, "failed to update destination account balance")
		}

		// Update transfer status
		now := time.Now().Format(time.RFC3339)
		if err := s.transferRepo.UpdateStatus(ctx, transfer.ID, "completed", &now); err != nil {
			return errors.Wrap(err, "failed to update transfer status")
		}

		return nil
	})

	if err != nil {
		// If transaction failed, mark transfer as failed
		now := time.Now().Format(time.RFC3339)
		_ = s.transferRepo.UpdateStatus(ctx, transfer.ID, "failed", &now)
		return nil, errors.Wrap(err, "transfer failed")
	}

	// Update balance cache
	_ = s.redisClient.SetBalanceCache(ctx, fromAccountID, fromAccount.Balance.Sub(amount))
	_ = s.redisClient.SetBalanceCache(ctx, toAccountID, toAccount.Balance.Add(amount))

	// Get the updated transfer
	updatedTransfer, err := s.transferRepo.GetByID(ctx, transfer.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get updated transfer")
	}

	return updatedTransfer, nil
}

// GetByID retrieves a transfer by ID
func (s *DefaultTransferService) GetByID(ctx context.Context, id uint64) (*model.Transfer, error) {
	transfer, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transfer")
	}

	return transfer, nil
}

// GetByAccountID retrieves transfers for an account with pagination
func (s *DefaultTransferService) GetByAccountID(ctx context.Context, accountID uuid.UUID, page, limit int) (*util.PaginatedResponse, error) {
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

	// Get transfers
	transfers, count, err := s.transferRepo.GetByAccountID(ctx, accountID, params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transfers")
	}

	// Create paginated response
	response := util.NewPaginatedResponse(transfers, params, count)
	return response, nil
}

// Helper function to convert uint to string
func uintToString(n uint64) string {
	return strconv.FormatUint(n, 10)
}
