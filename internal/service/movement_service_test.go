package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"VDM2-BankBE/internal/model"
	repmocks "VDM2-BankBE/internal/repository/mocks"
	"VDM2-BankBE/internal/service"
	servicemocks "VDM2-BankBE/internal/service/mocks"
	"VDM2-BankBE/internal/util"
)

func TestMovementService_Create(t *testing.T) {
	t.Parallel()

	accountID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440200")
	startBalance, _ := decimal.NewFromString("100.00")
	amount, _ := decimal.NewFromString("10.00")

	tests := []struct {
		name       string
		mType      string
		buildMocks func(ctrl *gomock.Controller) (*repmocks.MockMovementRepository, *repmocks.MockAccountRepository, *servicemocks.MockCacheClient)
		assert     func(t *testing.T, movement *model.Movement, err error)
	}{
		{
			name:  "invalid type returns 400",
			mType: "invalid",
			buildMocks: func(ctrl *gomock.Controller) (*repmocks.MockMovementRepository, *repmocks.MockAccountRepository, *servicemocks.MockCacheClient) {
				return repmocks.NewMockMovementRepository(ctrl), repmocks.NewMockAccountRepository(ctrl), servicemocks.NewMockCacheClient(ctrl)
			},
			assert: func(t *testing.T, movement *model.Movement, err error) {
				if movement != nil {
					t.Fatalf("expected nil movement, got %+v", movement)
				}
				apiErr, ok := err.(*util.APIError)
				if !ok || apiErr.Code != 400 {
					t.Fatalf("expected 400 APIError, got %#v", err)
				}
			},
		},
		{
			name:  "credit updates balance and writes movement",
			mType: "credit",
			buildMocks: func(ctrl *gomock.Controller) (*repmocks.MockMovementRepository, *repmocks.MockAccountRepository, *servicemocks.MockCacheClient) {
				movementRepo := repmocks.NewMockMovementRepository(ctrl)
				accountRepo := repmocks.NewMockAccountRepository(ctrl)
				cache := servicemocks.NewMockCacheClient(ctrl)

				accountRepo.EXPECT().GetByID(gomock.Any(), accountID).Return(&model.Account{ID: accountID, Balance: startBalance}, nil)
				accountRepo.EXPECT().UpdateBalance(gomock.Any(), accountID, amount).Return(nil)
				movementRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, m *model.Movement) error {
					if m.AccountID != accountID {
						t.Fatalf("unexpected movement account id: %s", m.AccountID.String())
					}
					if !m.Amount.Equal(amount) {
						t.Fatalf("unexpected movement amount: %s", m.Amount.String())
					}
					if m.Type != "credit" {
						t.Fatalf("unexpected movement type: %q", m.Type)
					}
					if m.Description != "desc" {
						t.Fatalf("unexpected movement description: %q", m.Description)
					}
					return nil
				})
				cache.EXPECT().SetBalanceCache(gomock.Any(), accountID, startBalance.Add(amount)).Return(nil)

				return movementRepo, accountRepo, cache
			},
			assert: func(t *testing.T, movement *model.Movement, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if movement == nil {
					t.Fatalf("expected movement, got nil")
				}
			},
		},
		{
			name:  "debit updates balance and writes movement",
			mType: "debit",
			buildMocks: func(ctrl *gomock.Controller) (*repmocks.MockMovementRepository, *repmocks.MockAccountRepository, *servicemocks.MockCacheClient) {
				movementRepo := repmocks.NewMockMovementRepository(ctrl)
				accountRepo := repmocks.NewMockAccountRepository(ctrl)
				cache := servicemocks.NewMockCacheClient(ctrl)

				accountRepo.EXPECT().GetByID(gomock.Any(), accountID).Return(&model.Account{ID: accountID, Balance: startBalance}, nil)
				accountRepo.EXPECT().UpdateBalance(gomock.Any(), accountID, amount.Neg()).Return(nil)
				movementRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, m *model.Movement) error {
					if m.Type != "debit" {
						t.Fatalf("unexpected movement type: %q", m.Type)
					}
					return nil
				})
				cache.EXPECT().SetBalanceCache(gomock.Any(), accountID, startBalance.Sub(amount)).Return(nil)

				return movementRepo, accountRepo, cache
			},
			assert: func(t *testing.T, movement *model.Movement, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if movement == nil {
					t.Fatalf("expected movement, got nil")
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

			movementRepo, accountRepo, cache := tc.buildMocks(ctrl)
			svc := service.NewMovementService(movementRepo, accountRepo, cache)

			m, err := svc.Create(context.Background(), accountID, amount, tc.mType, "desc")
			tc.assert(t, m, err)
		})
	}
}

