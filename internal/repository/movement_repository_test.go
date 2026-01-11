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

func TestGormMovementRepository_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	accountID := uuid.MustParse("550e8400-e29b-41d4-a716-446655441000")
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		setupSQL  func(m sqlmock.Sqlmock)
		exec      func(repo repository.MovementRepository) error
		assertErr func(t *testing.T, err error)
	}{
		{
			name: "success",
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectQuery(`INSERT INTO "movements" .*RETURNING "id","occurred_at"`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "occurred_at"}).AddRow(uint64(1), now))
				m.ExpectCommit()
			},
			exec: func(repo repository.MovementRepository) error {
				mv := &model.Movement{
					AccountID:   accountID,
					Amount:      mustDecimal(t, "10.00"),
					Type:        "credit",
					Description: "desc",
					OccurredAt:  now,
				}
				return repo.Create(ctx, mv)
			},
			assertErr: func(t *testing.T, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},
		{
			name: "db error wraps create failure",
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectQuery(`INSERT INTO "movements" .*RETURNING "id","occurred_at"`).
					WillReturnError(errors.New("db err"))
				m.ExpectRollback()
			},
			exec: func(repo repository.MovementRepository) error {
				mv := &model.Movement{
					AccountID:   accountID,
					Amount:      mustDecimal(t, "10.00"),
					Type:        "credit",
					Description: "desc",
					OccurredAt:  now,
				}
				return repo.Create(ctx, mv)
			},
			assertErr: func(t *testing.T, err error) {
				if err == nil || !regexp.MustCompile(`failed to create movement`).MatchString(err.Error()) {
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
			repo := repository.NewGormMovementRepository(dbm.DB)
			err := tc.exec(repo)
			tc.assertErr(t, err)

			if err := dbm.Mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet sqlmock expectations: %v", err)
			}
		})
	}
}

func TestGormMovementRepository_GetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	accountID := uuid.MustParse("550e8400-e29b-41d4-a716-446655441010")
	movementID := uint64(123)

	tests := []struct {
		name      string
		setupSQL  func(m sqlmock.Sqlmock)
		assertErr func(t *testing.T, mv *model.Movement, err error)
	}{
		{
			name: "success",
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT .* FROM "movements" WHERE id = \$1 ORDER BY .* LIMIT \$2`).
					WithArgs(movementID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "amount", "type", "description", "occurred_at"}).
						AddRow(movementID, accountID, "10.00", "credit", "desc", now))
			},
			assertErr: func(t *testing.T, mv *model.Movement, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if mv == nil || mv.ID != movementID {
					t.Fatalf("unexpected movement: %+v", mv)
				}
			},
		},
		{
			name: "not found maps to APIError 404",
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT .* FROM "movements" WHERE id = \$1 ORDER BY .* LIMIT \$2`).
					WithArgs(movementID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			assertErr: func(t *testing.T, mv *model.Movement, err error) {
				if mv != nil {
					t.Fatalf("expected nil movement, got %+v", mv)
				}
				var apiErr *util.APIError
				if !errors.As(err, &apiErr) {
					t.Fatalf("expected APIError, got %#v", err)
				}
				if apiErr.Code != 404 || apiErr.Message != "movement not found" {
					t.Fatalf("unexpected APIError: %+v", apiErr)
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
			repo := repository.NewGormMovementRepository(dbm.DB)
			mv, err := repo.GetByID(ctx, movementID)
			tc.assertErr(t, mv, err)

			if err := dbm.Mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet sqlmock expectations: %v", err)
			}
		})
	}
}

func TestGormMovementRepository_GetByAccountID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	accountID := uuid.MustParse("550e8400-e29b-41d4-a716-446655441020")
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	params := &util.PaginationParams{Page: 1, Limit: 10}

	tests := []struct {
		name      string
		setupSQL  func(m sqlmock.Sqlmock)
		assertErr func(t *testing.T, mvs []*model.Movement, count int, err error)
	}{
		{
			name: "success count + select",
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT count\(\*\) FROM "movements" WHERE account_id = \$1`).
					WithArgs(accountID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(2)))

				m.ExpectQuery(`SELECT .* FROM "movements" WHERE account_id = \$1 ORDER BY occurred_at DESC LIMIT \$2`).
					WithArgs(accountID, params.Limit).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "amount", "type", "description", "occurred_at"}).
						AddRow(uint64(1), accountID, "10.00", "credit", "a", now).
						AddRow(uint64(2), accountID, "5.00", "debit", "b", now))
			},
			assertErr: func(t *testing.T, mvs []*model.Movement, count int, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if count != 2 {
					t.Fatalf("unexpected count: %d", count)
				}
				if len(mvs) != 2 {
					t.Fatalf("unexpected len: %d", len(mvs))
				}
			},
		},
		{
			name: "count error wraps",
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT count\(\*\) FROM "movements" WHERE account_id = \$1`).
					WithArgs(accountID).
					WillReturnError(errors.New("count err"))
			},
			assertErr: func(t *testing.T, mvs []*model.Movement, count int, err error) {
				if err == nil || !regexp.MustCompile(`failed to count movements`).MatchString(err.Error()) {
					t.Fatalf("expected wrapped error, got: %v", err)
				}
			},
		},
		{
			name: "select error wraps",
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT count\(\*\) FROM "movements" WHERE account_id = \$1`).
					WithArgs(accountID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(2)))

				m.ExpectQuery(`SELECT .* FROM "movements" WHERE account_id = \$1 ORDER BY occurred_at DESC LIMIT \$2`).
					WithArgs(accountID, params.Limit).
					WillReturnError(errors.New("select err"))
			},
			assertErr: func(t *testing.T, mvs []*model.Movement, count int, err error) {
				if err == nil || !regexp.MustCompile(`failed to get movements by account ID`).MatchString(err.Error()) {
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
			repo := repository.NewGormMovementRepository(dbm.DB)
			mvs, count, err := repo.GetByAccountID(ctx, accountID, params)
			tc.assertErr(t, mvs, count, err)

			if err := dbm.Mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet sqlmock expectations: %v", err)
			}
		})
	}
}

func mustDecimal(t *testing.T, s string) decimal.Decimal {
	t.Helper()
	d, err := decimal.NewFromString(s)
	if err != nil {
		t.Fatalf("bad decimal %q: %v", s, err)
	}
	return d
}

