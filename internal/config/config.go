package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig
	DB       DBConfig
	Redis    RedisConfig
	JWT      JWTConfig
	OAuth    OAuthConfig
	Logging  LoggingConfig
	Security SecurityConfig
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Port    int
	Timeout time.Duration
	Debug   bool
}

// DBConfig holds the database configuration
type DBConfig struct {
	Driver          string
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	AutoMigrate     bool
}

// RedisConfig holds the Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// JWTConfig holds the JWT configuration
type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

// OAuthConfig holds OAuth configuration
type OAuthConfig struct {
	Google OAuthProviderConfig
}

// OAuthProviderConfig holds configuration for an OAuth provider
type OAuthProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	RateLimit RateLimitConfig
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled  bool
	Requests int
	Duration time.Duration
}

// Load loads the configuration from a file
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Set up viper
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")
	viper.SetConfigType("yaml")

	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.timeout", "30s")
	viper.SetDefault("server.debug", true)

	// Read the config
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}

	// Parse the config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, errors.Wrap(err, "failed to parse config")
	}

	// Validate critical config
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// GetDBURL returns the PostgreSQL connection string
func (c *DBConfig) GetDBURL() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// GetRedisAddr returns the Redis address string
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// validateConfig validates that all required config values are present
func validateConfig(config *Config) error {
	// Validate DB config
	if config.DB.Host == "" {
		return errors.New("database host is required")
	}
	if config.DB.User == "" {
		return errors.New("database user is required")
	}
	if config.DB.DBName == "" {
		return errors.New("database name is required")
	}

	// Validate JWT config
	if config.JWT.Secret == "" {
		return errors.New("JWT secret is required")
	}

	return nil
}
