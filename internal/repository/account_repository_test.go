package repository_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"VDM2-BankBE/internal/model"
	"VDM2-BankBE/internal/repository"
	"VDM2-BankBE/internal/testutil"
	"VDM2-BankBE/internal/util"
)

func TestGormAccountRepository_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name       string
		account    *model.Account
		setupSQL   func(m sqlmock.Sqlmock)
		assertFunc func(t *testing.T, account *model.Account, err error)
	}{
		{
			name:    "success assigns id if nil",
			account: &model.Account{ID: uuid.Nil, UserID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440911"), Currency: "EUR"},
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(`INSERT INTO "accounts" .*`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 1))
				m.ExpectCommit()
			},
			assertFunc: func(t *testing.T, account *model.Account, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if account.ID == uuid.Nil {
					t.Fatalf("expected ID to be set")
				}
			},
		},
		{
			name:    "db error wraps create failure",
			account: &model.Account{ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440912"), UserID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440913"), Currency: "EUR"},
			setupSQL: func(m sqlmock.Sqlmock) {
				baseErr := errors.New("db err")
				m.ExpectBegin()
				m.ExpectExec(`INSERT INTO "accounts" .*`).
					WillReturnError(baseErr)
				m.ExpectRollback()
			},
			assertFunc: func(t *testing.T, account *model.Account, err error) {
				if err == nil {
					t.Fatalf("expected error")
				}
				if !regexp.MustCompile(`failed to create account`).MatchString(err.Error()) {
					t.Fatalf("expected wrapped error, got: %v", err)
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dbm := testutil.NewGormSQLMock(t)
			defer dbm.Cleanup()

			tc.setupSQL(dbm.Mock)
			repo := repository.NewGormAccountRepository(dbm.DB)
			err := repo.Create(ctx, tc.account)

			tc.assertFunc(t, tc.account, err)
			if err := dbm.Mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet sqlmock expectations: %v", err)
			}
		})
	}
}

func TestGormAccountRepository_GetByID_GetByUserID_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	accountID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440920")
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440921")

	type testCase struct {
		name      string
		setupSQL  func(m sqlmock.Sqlmock)
		exec      func(repo repository.AccountRepository) error
		assertErr func(t *testing.T, err error)
	}

	mkAccountRow := func() *sqlmock.Rows {
		return sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "created_at", "updated_at"}).
			AddRow(accountID, userID, "10.00", "EUR", now, now)
	}

	tests := []testCase{
		{
			name: "GetByID success",
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT .* FROM "accounts" WHERE id = \$1 ORDER BY .* LIMIT \$2`).
					WithArgs(accountID, 1).
					WillReturnRows(mkAccountRow())
			},
			exec: func(repo repository.AccountRepository) error {
				got, err := repo.GetByID(ctx, accountID)
				if err != nil {
					return err
				}
				if got.ID != accountID {
					return errors.New("unexpected account id")
				}
				return nil
			},
			assertErr: func(t *testing.T, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},
		{
			name: "GetByID not found maps to APIError 404",
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT .* FROM "accounts" WHERE id = \$1 ORDER BY .* LIMIT \$2`).
					WithArgs(accountID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			exec: func(repo repository.AccountRepository) error {
				_, err := repo.GetByID(ctx, accountID)
				return err
			},
			assertErr: func(t *testing.T, err error) {
				var apiErr *util.APIError
				if !errors.As(err, &apiErr) {
					t.Fatalf("expected APIError, got: %#v", err)
				}
				if apiErr.Code != 404 || apiErr.Message != "account not found" {
					t.Fatalf("unexpected APIError: %+v", apiErr)
				}
			},
		},
		{
			name: "GetByUserID success",
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT .* FROM "accounts" WHERE user_id = \$1 ORDER BY .* LIMIT \$2`).
					WithArgs(userID, 1).
					WillReturnRows(mkAccountRow())
			},
			exec: func(repo repository.AccountRepository) error {
				got, err := repo.GetByUserID(ctx, userID)
				if err != nil {
					return err
				}
				if got.UserID != userID {
					return errors.New("unexpected user id")
				}
				return nil
			},
			assertErr: func(t *testing.T, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},
		{
			name: "Delete success",
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(`DELETE FROM "accounts" WHERE id = \$1`).
					WithArgs(accountID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				m.ExpectCommit()
			},
			exec: func(repo repository.AccountRepository) error {
				return repo.Delete(ctx, accountID)
			},
			assertErr: func(t *testing.T, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dbm := testutil.NewGormSQLMock(t)
			defer dbm.Cleanup()

			tc.setupSQL(dbm.Mock)
			repo := repository.NewGormAccountRepository(dbm.DB)
			err := tc.exec(repo)

			tc.assertErr(t, err)
			if err := dbm.Mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet sqlmock expectations: %v", err)
			}
		})
	}
}

func TestGormAccountRepository_UpdateBalance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	accountID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440930")
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	type testCase struct {
		name      string
		amount    decimal.Decimal
		setupSQL  func(m sqlmock.Sqlmock)
		assertErr func(t *testing.T, err error)
	}

	selectForUpdateRegex := `SELECT .* FROM "accounts" WHERE id = \$1 ORDER BY .* LIMIT 1.*FOR UPDATE`
	selectRegex := `SELECT .* FROM "accounts" WHERE id = \$1 ORDER BY .* LIMIT \$2`
	updateAnyPlaceholderRegex := `UPDATE "accounts" SET .* WHERE "id" = \$\d+`

	tests := []testCase{
		{
			name:   "success BEGIN -> SELECT FOR UPDATE -> UPDATE -> COMMIT",
			amount: decimal.NewFromInt(5),
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectQuery(selectRegex).
					WithArgs(accountID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "created_at", "updated_at"}).
						AddRow(accountID, uuid.MustParse("550e8400-e29b-41d4-a716-446655440931"), "10.00", "EUR", now, now))
				m.ExpectExec(updateAnyPlaceholderRegex).
					WillReturnResult(sqlmock.NewResult(0, 1))
				m.ExpectCommit()
			},
			assertErr: func(t *testing.T, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},
		{
			name:   "begin error returns wrapped error",
			amount: decimal.NewFromInt(1),
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectBegin().WillReturnError(errors.New("begin failed"))
			},
			assertErr: func(t *testing.T, err error) {
				if err == nil || !regexp.MustCompile(`failed to begin transaction`).MatchString(err.Error()) {
					t.Fatalf("expected begin error, got: %v", err)
				}
			},
		},
		{
			name:   "record not found rolls back and returns APIError 404",
			amount: decimal.NewFromInt(1),
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectQuery(selectRegex).
					WithArgs(accountID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
				m.ExpectRollback()
			},
			assertErr: func(t *testing.T, err error) {
				var apiErr *util.APIError
				if !errors.As(err, &apiErr) {
					t.Fatalf("expected APIError, got: %#v", err)
				}
				if apiErr.Code != 404 || apiErr.Message != "account not found" {
					t.Fatalf("unexpected APIError: %+v", apiErr)
				}
			},
		},
		{
			name:   "insufficient funds rolls back and returns APIError 400",
			amount: decimal.RequireFromString("-20.00"),
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectQuery(selectRegex).
					WithArgs(accountID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "created_at", "updated_at"}).
						AddRow(accountID, uuid.MustParse("550e8400-e29b-41d4-a716-446655440932"), "10.00", "EUR", now, now))
				m.ExpectRollback()
			},
			assertErr: func(t *testing.T, err error) {
				var apiErr *util.APIError
				if !errors.As(err, &apiErr) {
					t.Fatalf("expected APIError, got: %#v", err)
				}
				if apiErr.Code != 400 || apiErr.Message != "insufficient funds" {
					t.Fatalf("unexpected APIError: %+v", apiErr)
				}
			},
		},
		{
			name:   "save error rolls back and returns wrapped error",
			amount: decimal.NewFromInt(1),
			setupSQL: func(m sqlmock.Sqlmock) {
				baseErr := errors.New("update failed")
				m.ExpectBegin()
				m.ExpectQuery(selectRegex).
					WithArgs(accountID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "created_at", "updated_at"}).
						AddRow(accountID, uuid.MustParse("550e8400-e29b-41d4-a716-446655440933"), "10.00", "EUR", now, now))
				m.ExpectExec(updateAnyPlaceholderRegex).
					WillReturnError(baseErr)
				m.ExpectRollback()
			},
			assertErr: func(t *testing.T, err error) {
				if err == nil || !regexp.MustCompile(`failed to update account balance`).MatchString(err.Error()) {
					t.Fatalf("expected update error, got: %v", err)
				}
			},
		},
		{
			name:   "commit error returns wrapped error",
			amount: decimal.NewFromInt(1),
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectQuery(selectRegex).
					WithArgs(accountID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "created_at", "updated_at"}).
						AddRow(accountID, uuid.MustParse("550e8400-e29b-41d4-a716-446655440934"), "10.00", "EUR", now, now))
				m.ExpectExec(updateAnyPlaceholderRegex).
					WillReturnResult(sqlmock.NewResult(0, 1))
				m.ExpectCommit().WillReturnError(errors.New("commit failed"))
			},
			assertErr: func(t *testing.T, err error) {
				if err == nil || !regexp.MustCompile(`failed to commit transaction`).MatchString(err.Error()) {
					t.Fatalf("expected commit error, got: %v", err)
				}
			},
		},
	}

	_ = selectForUpdateRegex
	_ = updateAnyPlaceholderRegex
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dbm := testutil.NewGormSQLMock(t)
			defer dbm.Cleanup()

			tc.setupSQL(dbm.Mock)
			repo := repository.NewGormAccountRepository(dbm.DB)
			err := repo.UpdateBalance(ctx, accountID, tc.amount)
			tc.assertErr(t, err)

			if err := dbm.Mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet sqlmock expectations: %v", err)
			}
		})
	}
}

