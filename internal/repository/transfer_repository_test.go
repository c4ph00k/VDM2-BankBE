package repository_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"VDM2-BankBE/internal/model"
	"VDM2-BankBE/internal/repository"
	"VDM2-BankBE/internal/testutil"
	"VDM2-BankBE/internal/util"
)

func TestGormTransferRepository_GetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	id := uint64(10)
	fromID := uuid.MustParse("550e8400-e29b-41d4-a716-446655441100")
	toID := uuid.MustParse("550e8400-e29b-41d4-a716-446655441101")

	tests := []struct {
		name      string
		setupSQL  func(m sqlmock.Sqlmock)
		assertErr func(t *testing.T, tr *model.Transfer, err error)
	}{
		{
			name: "success",
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT .* FROM "transfers" WHERE id = \$1 ORDER BY .* LIMIT \$2`).
					WithArgs(id, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "from_account", "to_account", "amount", "status", "initiated_at", "completed_at"}).
						AddRow(id, fromID, toID, "25.00", "completed", now, nil))
			},
			assertErr: func(t *testing.T, tr *model.Transfer, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tr == nil || tr.ID != id {
					t.Fatalf("unexpected transfer: %+v", tr)
				}
			},
		},
		{
			name: "not found maps to APIError 404",
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT .* FROM "transfers" WHERE id = \$1 ORDER BY .* LIMIT \$2`).
					WithArgs(id, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			assertErr: func(t *testing.T, tr *model.Transfer, err error) {
				if tr != nil {
					t.Fatalf("expected nil transfer, got %+v", tr)
				}
				var apiErr *util.APIError
				if !errors.As(err, &apiErr) {
					t.Fatalf("expected APIError, got %#v", err)
				}
				if apiErr.Code != 404 || apiErr.Message != "transfer not found" {
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
			repo := repository.NewGormTransferRepository(dbm.DB)
			tr, err := repo.GetByID(ctx, id)
			tc.assertErr(t, tr, err)

			if err := dbm.Mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet sqlmock expectations: %v", err)
			}
		})
	}
}

func TestGormTransferRepository_GetByAccountID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	accountID := uuid.MustParse("550e8400-e29b-41d4-a716-446655441120")
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	params := &util.PaginationParams{Page: 1, Limit: 10}

	tests := []struct {
		name      string
		setupSQL  func(m sqlmock.Sqlmock)
		assertErr func(t *testing.T, trs []*model.Transfer, count int, err error)
	}{
		{
			name: "success count + select",
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT count\(\*\) FROM "transfers" WHERE from_account = \$1 OR to_account = \$2`).
					WithArgs(accountID, accountID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(1)))

				m.ExpectQuery(`SELECT .* FROM "transfers" WHERE from_account = \$1 OR to_account = \$2 ORDER BY initiated_at DESC LIMIT \$3`).
					WithArgs(accountID, accountID, params.Limit).
					WillReturnRows(sqlmock.NewRows([]string{"id", "from_account", "to_account", "amount", "status", "initiated_at", "completed_at"}).
						AddRow(uint64(1), accountID, uuid.MustParse("550e8400-e29b-41d4-a716-446655441121"), "25.00", "completed", now, nil))
			},
			assertErr: func(t *testing.T, trs []*model.Transfer, count int, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if count != 1 || len(trs) != 1 {
					t.Fatalf("unexpected result: count=%d len=%d", count, len(trs))
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
			repo := repository.NewGormTransferRepository(dbm.DB)
			trs, count, err := repo.GetByAccountID(ctx, accountID, params)
			tc.assertErr(t, trs, count, err)

			if err := dbm.Mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet sqlmock expectations: %v", err)
			}
		})
	}
}

func TestGormTransferRepository_UpdateStatus(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uint64(99)

	tests := []struct {
		name      string
		status    string
		completed *string
		setupSQL  func(m sqlmock.Sqlmock)
		assertErr func(t *testing.T, err error)
	}{
		{
			name:      "invalid RFC3339 returns error and does not hit DB",
			status:    "completed",
			completed: ptr("not-a-time"),
			setupSQL:  func(m sqlmock.Sqlmock) {},
			assertErr: func(t *testing.T, err error) {
				if err == nil || !regexp.MustCompile(`invalid completedAt time format`).MatchString(err.Error()) {
					t.Fatalf("expected invalid time error, got: %v", err)
				}
			},
		},
		{
			name:      "nil completedAt updates status only",
			status:    "failed",
			completed: nil,
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(`UPDATE "transfers" SET .*status.* WHERE id = \$`).
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
			name:      "valid completedAt updates status and completed_at",
			status:    "completed",
			completed: ptr("2026-01-01T00:00:00Z"),
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(`UPDATE "transfers" SET .*completed_at.*status.* WHERE id = \$`).
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
			name:      "db error wraps update failure",
			status:    "failed",
			completed: nil,
			setupSQL: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(`UPDATE "transfers" SET .*status.* WHERE id = \$`).
					WillReturnError(errors.New("db err"))
				m.ExpectRollback()
			},
			assertErr: func(t *testing.T, err error) {
				if err == nil || !regexp.MustCompile(`failed to update transfer status`).MatchString(err.Error()) {
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
			repo := repository.NewGormTransferRepository(dbm.DB)
			err := repo.UpdateStatus(ctx, id, tc.status, tc.completed)
			tc.assertErr(t, err)

			if err := dbm.Mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet sqlmock expectations: %v", err)
			}
		})
	}
}

func ptr(s string) *string { return &s }

