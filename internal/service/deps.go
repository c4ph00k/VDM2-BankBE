package service

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"VDM2-BankBE/pkg/oauth"
)

// CacheClient represents the cache/store boundary used by services.
// Implemented by `pkg/cache.RedisClient`.
//go:generate mockgen -destination=./mocks/mock_cache_client.go -package=mocks VDM2-BankBE/internal/service CacheClient
type CacheClient interface {
	// Balance cache
	SetBalanceCache(ctx context.Context, accountID uuid.UUID, balance decimal.Decimal) error
	GetBalanceCache(ctx context.Context, accountID uuid.UUID) (decimal.Decimal, error)

	// OAuth state store
	SetOAuthState(ctx context.Context, state string, redirectURL string) error
	GetOAuthState(ctx context.Context, state string) (string, error)
}

// GoogleOAuthClient represents the Google OAuth boundary used by the auth service.
// Implemented by `pkg/oauth.GoogleOAuthClient`.
//go:generate mockgen -destination=./mocks/mock_google_oauth_client.go -package=mocks VDM2-BankBE/internal/service GoogleOAuthClient
type GoogleOAuthClient interface {
	GetAuthURL(state string) string
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*oauth.GoogleUserInfo, error)
}

// TxDB represents the DB transaction boundary used by the transfer service.
// Implemented by `*gorm.DB`.
//go:generate mockgen -destination=./mocks/mock_tx_db.go -package=mocks VDM2-BankBE/internal/service TxDB
type TxDB interface {
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error
}

