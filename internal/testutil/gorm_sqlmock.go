package testutil

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormSQLMock struct {
	DB      *gorm.DB
	SQLDB   *sql.DB
	Mock    sqlmock.Sqlmock
	Cleanup func()
}

func NewGormSQLMock(t *testing.T) GormSQLMock {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		DisableAutomaticPing: true,
	})
	if err != nil {
		_ = sqlDB.Close()
		t.Fatalf("failed to open gorm db: %v", err)
	}

	cleanup := func() {
		t.Helper()
		_ = sqlDB.Close()
	}

	return GormSQLMock{
		DB:      gdb,
		SQLDB:   sqlDB,
		Mock:    mock,
		Cleanup: cleanup,
	}
}
