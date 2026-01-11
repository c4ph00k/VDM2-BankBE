package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"VDM2-BankBE/internal/model"
	repmocks "VDM2-BankBE/internal/repository/mocks"
	"VDM2-BankBE/internal/service"
	servicemocks "VDM2-BankBE/internal/service/mocks"
)

func TestAccountService_GetBalance(t *testing.T) {
	t.Parallel()

	accountID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440100")
	expectedBalance, _ := decimal.NewFromString("10.50")

	tests := []struct {
		name       string
		buildMocks func(ctrl *gomock.Controller) (*repmocks.MockAccountRepository, *servicemocks.MockCacheClient)
		assert     func(t *testing.T, got decimal.Decimal, err error)
	}{
		{
			name: "cache hit returns cached value",
			buildMocks: func(ctrl *gomock.Controller) (*repmocks.MockAccountRepository, *servicemocks.MockCacheClient) {
				accountRepo := repmocks.NewMockAccountRepository(ctrl)
				cache := servicemocks.NewMockCacheClient(ctrl)
				cache.EXPECT().GetBalanceCache(gomock.Any(), accountID).Return(decimal.NewFromInt(123), nil)
				return accountRepo, cache
			},
			assert: func(t *testing.T, got decimal.Decimal, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !got.Equal(decimal.NewFromInt(123)) {
					t.Fatalf("unexpected balance: got=%s want=%s", got.String(), decimal.NewFromInt(123).String())
				}
			},
		},
		{
			name: "cache miss loads from repo and updates cache",
			buildMocks: func(ctrl *gomock.Controller) (*repmocks.MockAccountRepository, *servicemocks.MockCacheClient) {
				accountRepo := repmocks.NewMockAccountRepository(ctrl)
				cache := servicemocks.NewMockCacheClient(ctrl)

				cache.EXPECT().GetBalanceCache(gomock.Any(), accountID).Return(decimal.Zero, errors.New("cache miss"))
				accountRepo.EXPECT().GetByID(gomock.Any(), accountID).Return(&model.Account{ID: accountID, Balance: expectedBalance}, nil)
				cache.EXPECT().SetBalanceCache(gomock.Any(), accountID, expectedBalance).Return(nil)

				return accountRepo, cache
			},
			assert: func(t *testing.T, got decimal.Decimal, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !got.Equal(expectedBalance) {
					t.Fatalf("unexpected balance: got=%s want=%s", got.String(), expectedBalance.String())
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

			accountRepo, cache := tc.buildMocks(ctrl)
			svc := service.NewAccountService(accountRepo, cache)

			got, err := svc.GetBalance(context.Background(), accountID)
			tc.assert(t, got, err)
		})
	}
}

