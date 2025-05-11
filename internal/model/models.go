package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	Username     string    `gorm:"uniqueIndex;not null" json:"username"`
	FirstName    string    `gorm:"not null" json:"first_name"`
	LastName     string    `gorm:"not null" json:"last_name"`
	FiscalCode   string    `gorm:"uniqueIndex;not null" json:"fiscal_code"`
	PasswordHash string    `gorm:"not null" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Account represents a user's bank account
type Account struct {
	ID        uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
	User      User            `gorm:"foreignKey:UserID" json:"-"`
	Balance   decimal.Decimal `gorm:"type:numeric(18,2);not null;default:0" json:"balance"`
	Currency  string          `gorm:"type:text;not null;default:'EUR'" json:"currency"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// Movement represents a transaction within an account
type Movement struct {
	ID          uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID   uuid.UUID       `gorm:"type:uuid;not null" json:"account_id"`
	Account     Account         `gorm:"foreignKey:AccountID" json:"-"`
	Amount      decimal.Decimal `gorm:"type:numeric(18,2);not null" json:"amount"`
	Type        string          `gorm:"type:text;not null;check:type IN ('credit','debit')" json:"type"`
	Description string          `gorm:"type:text" json:"description"`
	OccurredAt  time.Time       `gorm:"not null;default:now()" json:"occurred_at"`
}

// OAuthToken represents an OAuth token for a user
type OAuthToken struct {
	UserID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"-"`
	Provider     string    `gorm:"type:text;not null" json:"provider"`
	AccessToken  string    `gorm:"type:text;not null" json:"-"`
	RefreshToken string    `gorm:"type:text" json:"-"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// Transfer represents a transfer between two accounts
type Transfer struct {
	ID          uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	FromAccount uuid.UUID       `gorm:"type:uuid;not null" json:"from_account"`
	ToAccount   uuid.UUID       `gorm:"type:uuid;not null" json:"to_account"`
	Amount      decimal.Decimal `gorm:"type:numeric(18,2);not null" json:"amount"`
	Status      string          `gorm:"type:text;not null;check:status IN ('pending','completed','failed')" json:"status"`
	InitiatedAt time.Time       `gorm:"not null;default:now()" json:"initiated_at"`
	CompletedAt *time.Time      `json:"completed_at"`
}

// TableName sets the table names explicitly
func (*User) TableName() string {
	return "users"
}

func (*Account) TableName() string {
	return "accounts"
}

func (*Movement) TableName() string {
	return "movements"
}

func (*OAuthToken) TableName() string {
	return "oauth_tokens"
}

func (*Transfer) TableName() string {
	return "transfers"
}
