package service_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"VDM2-BankBE/internal/model"
	repmocks "VDM2-BankBE/internal/repository/mocks"
	"VDM2-BankBE/internal/service"
	servicemocks "VDM2-BankBE/internal/service/mocks"
	"VDM2-BankBE/internal/util"
)

func TestTransferService_Transfer(t *testing.T) {
	t.Parallel()

	fromAccountID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440300")
	toAccountID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440301")

	startFromBalance, _ := decimal.NewFromString("100.00")
	startToBalance, _ := decimal.NewFromString("50.00")
	amount, _ := decimal.NewFromString("25.00")

	tests := []struct {
		name       string
		amount     decimal.Decimal
		buildMocks func(ctrl *gomock.Controller) (
			*repmocks.MockTransferRepository,
			*repmocks.MockAccountRepository,
			*repmocks.MockMovementRepository,
			*servicemocks.MockCacheClient,
			*servicemocks.MockTxDB,
		)
		assert func(t *testing.T, got *model.Transfer, err error)
	}{
		{
			name:   "amount must be > 0",
			amount: decimal.Zero,
			buildMocks: func(ctrl *gomock.Controller) (
				*repmocks.MockTransferRepository,
				*repmocks.MockAccountRepository,
				*repmocks.MockMovementRepository,
				*servicemocks.MockCacheClient,
				*servicemocks.MockTxDB,
			) {
				return repmocks.NewMockTransferRepository(ctrl),
					repmocks.NewMockAccountRepository(ctrl),
					repmocks.NewMockMovementRepository(ctrl),
					servicemocks.NewMockCacheClient(ctrl),
					servicemocks.NewMockTxDB(ctrl)
			},
			assert: func(t *testing.T, got *model.Transfer, err error) {
				if got != nil {
					t.Fatalf("expected nil transfer, got %+v", got)
				}
				apiErr, ok := err.(*util.APIError)
				if !ok || apiErr.Code != 400 {
					t.Fatalf("expected 400 APIError, got %#v", err)
				}
			},
		},
		{
			name:   "insufficient funds returns 400",
			amount: amount,
			buildMocks: func(ctrl *gomock.Controller) (
				*repmocks.MockTransferRepository,
				*repmocks.MockAccountRepository,
				*repmocks.MockMovementRepository,
				*servicemocks.MockCacheClient,
				*servicemocks.MockTxDB,
			) {
				transferRepo := repmocks.NewMockTransferRepository(ctrl)
				accountRepo := repmocks.NewMockAccountRepository(ctrl)
				movementRepo := repmocks.NewMockMovementRepository(ctrl)
				cache := servicemocks.NewMockCacheClient(ctrl)
				txdb := servicemocks.NewMockTxDB(ctrl)

				accountRepo.EXPECT().GetByID(gomock.Any(), fromAccountID).Return(&model.Account{ID: fromAccountID, Balance: decimal.NewFromInt(1)}, nil)
				accountRepo.EXPECT().GetByID(gomock.Any(), toAccountID).Return(&model.Account{ID: toAccountID, Balance: startToBalance}, nil)

				return transferRepo, accountRepo, movementRepo, cache, txdb
			},
			assert: func(t *testing.T, got *model.Transfer, err error) {
				if got != nil {
					t.Fatalf("expected nil transfer, got %+v", got)
				}
				apiErr, ok := err.(*util.APIError)
				if !ok || apiErr.Code != 400 || apiErr.Message != "insufficient funds" {
					t.Fatalf("unexpected error: %#v", err)
				}
			},
		},
		{
			name:   "success runs transaction and updates cache",
			amount: amount,
			buildMocks: func(ctrl *gomock.Controller) (
				*repmocks.MockTransferRepository,
				*repmocks.MockAccountRepository,
				*repmocks.MockMovementRepository,
				*servicemocks.MockCacheClient,
				*servicemocks.MockTxDB,
			) {
				transferRepo := repmocks.NewMockTransferRepository(ctrl)
				accountRepo := repmocks.NewMockAccountRepository(ctrl)
				movementRepo := repmocks.NewMockMovementRepository(ctrl)
				cache := servicemocks.NewMockCacheClient(ctrl)
				txdb := servicemocks.NewMockTxDB(ctrl)

				fromAccount := &model.Account{ID: fromAccountID, Balance: startFromBalance}
				toAccount := &model.Account{ID: toAccountID, Balance: startToBalance}

				accountRepo.EXPECT().GetByID(gomock.Any(), fromAccountID).Return(fromAccount, nil)
				accountRepo.EXPECT().GetByID(gomock.Any(), toAccountID).Return(toAccount, nil)

				txdb.EXPECT().
					Transaction(gomock.Any()).
					DoAndReturn(func(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
						return fc(&gorm.DB{})
					})

				transferRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, tr *model.Transfer) error {
					if tr.Status != "pending" {
						t.Fatalf("unexpected transfer status: %q", tr.Status)
					}
					if !tr.Amount.Equal(amount) {
						t.Fatalf("unexpected transfer amount: %s", tr.Amount.String())
					}
					return nil
				})

				movementRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(2).Return(nil)
				accountRepo.EXPECT().UpdateBalance(gomock.Any(), fromAccountID, amount.Neg()).Return(nil)
				accountRepo.EXPECT().UpdateBalance(gomock.Any(), toAccountID, amount).Return(nil)
				transferRepo.EXPECT().UpdateStatus(gomock.Any(), uint64(0), "completed", gomock.Any()).Return(nil)

				cache.EXPECT().SetBalanceCache(gomock.Any(), fromAccountID, startFromBalance.Sub(amount)).Return(nil)
				cache.EXPECT().SetBalanceCache(gomock.Any(), toAccountID, startToBalance.Add(amount)).Return(nil)

				updated := &model.Transfer{
					ID:          0,
					FromAccount: fromAccountID,
					ToAccount:   toAccountID,
					Amount:      amount,
					Status:      "completed",
					InitiatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				}
				transferRepo.EXPECT().GetByID(gomock.Any(), uint64(0)).Return(updated, nil)

				return transferRepo, accountRepo, movementRepo, cache, txdb
			},
			assert: func(t *testing.T, got *model.Transfer, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got == nil || got.Status != "completed" {
					t.Fatalf("unexpected transfer: %+v", got)
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			transferRepo, accountRepo, movementRepo, cache, txdb := tc.buildMocks(ctrl)
			svc := service.NewTransferService(transferRepo, accountRepo, movementRepo, cache, txdb)

			got, err := svc.Transfer(context.Background(), fromAccountID, toAccountID, tc.amount, "desc")
			tc.assert(t, got, err)
		})
	}
}

