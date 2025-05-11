package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"VDM2-BankBE/internal/config"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// RedisClient wraps the redis.Client with application-specific methods
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg *config.RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(err, "failed to connect to Redis")
	}

	return &RedisClient{client: client}, nil
}

// Close closes the Redis client
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// SetSession stores a session in Redis
func (r *RedisClient) SetSession(ctx context.Context, sessionID string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal session data")
	}
	return r.client.Set(ctx, "sess:"+sessionID, data, 24*time.Hour).Err()
}

// GetSession retrieves a session from Redis
func (r *RedisClient) GetSession(ctx context.Context, sessionID string, dest interface{}) error {
	data, err := r.client.Get(ctx, "sess:"+sessionID).Bytes()
	if err != nil {
		if err == redis.Nil {
			return errors.New("session not found")
		}
		return errors.Wrap(err, "failed to get session")
	}

	return json.Unmarshal(data, dest)
}

// DeleteSession deletes a session from Redis
func (r *RedisClient) DeleteSession(ctx context.Context, sessionID string) error {
	return r.client.Del(ctx, "sess:"+sessionID).Err()
}

// IncrRateLimit increments the rate limit counter and returns the current count
func (r *RedisClient) IncrRateLimit(ctx context.Context, userID, route string, window time.Duration) (int64, error) {
	key := fmt.Sprintf("rl:%s:%s", userID, route)
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, errors.Wrap(err, "failed to increment rate limit counter")
	}
	if count == 1 {
		r.client.Expire(ctx, key, window)
	}
	return count, nil
}

// SetBalanceCache stores an account balance in the cache
func (r *RedisClient) SetBalanceCache(ctx context.Context, accountID uuid.UUID, balance decimal.Decimal) error {
	key := "acct:balance:" + accountID.String()
	return r.client.Set(ctx, key, balance.String(), 1*time.Minute).Err()
}

// GetBalanceCache retrieves an account balance from the cache
func (r *RedisClient) GetBalanceCache(ctx context.Context, accountID uuid.UUID) (decimal.Decimal, error) {
	key := "acct:balance:" + accountID.String()
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return decimal.Zero, errors.New("balance not found in cache")
		}
		return decimal.Zero, errors.Wrap(err, "failed to get balance from cache")
	}

	balance, err := decimal.NewFromString(val)
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "failed to parse balance")
	}

	return balance, nil
}

// SetOTPCode stores a one-time password/verification code
func (r *RedisClient) SetOTPCode(ctx context.Context, userID uuid.UUID, purpose string, code string) error {
	key := fmt.Sprintf("otp:%s:%s", userID.String(), purpose)
	return r.client.Set(ctx, key, code, 5*time.Minute).Err()
}

// GetOTPCode retrieves a one-time password/verification code
func (r *RedisClient) GetOTPCode(ctx context.Context, userID uuid.UUID, purpose string) (string, error) {
	key := fmt.Sprintf("otp:%s:%s", userID.String(), purpose)
	code, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("OTP code not found or expired")
		}
		return "", errors.Wrap(err, "failed to get OTP code")
	}
	return code, nil
}

// SetOAuthState stores an OAuth state token
func (r *RedisClient) SetOAuthState(ctx context.Context, state string, redirectURL string) error {
	key := "oauth:state:" + state
	return r.client.Set(ctx, key, redirectURL, 10*time.Minute).Err()
}

// GetOAuthState retrieves an OAuth state token
func (r *RedisClient) GetOAuthState(ctx context.Context, state string) (string, error) {
	key := "oauth:state:" + state
	redirectURL, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("OAuth state not found or expired")
		}
		return "", errors.Wrap(err, "failed to get OAuth state")
	}

	// Delete the state after retrieving it
	r.client.Del(ctx, key)

	return redirectURL, nil
}
