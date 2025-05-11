package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"VDM2-BankBE/internal/config"
	"VDM2-BankBE/internal/handler"
	"VDM2-BankBE/internal/middleware"
	"VDM2-BankBE/internal/repository"
	"VDM2-BankBE/internal/router"
	"VDM2-BankBE/internal/service"
	"VDM2-BankBE/pkg/cache"
	"VDM2-BankBE/pkg/oauth"
)

// @title VDM2 Banking API
// @version 1.0
// @description A RESTful API for banking operations including authentication, account management, and transfers
// @contact.name API Support
// @contact.email support@vdmsquare.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	var logger *zap.Logger
	if cfg.Server.Debug {
		logger, _ = zap.NewDevelopment()
		gin.SetMode(gin.DebugMode)
	} else {
		logger, _ = zap.NewProduction()
		gin.SetMode(gin.ReleaseMode)
	}
	defer logger.Sync()

	// Connect to database
	db, err := connectToDatabase(cfg.DB, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Connect to Redis
	redisClient, err := cache.NewRedisClient(&cfg.Redis)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.Close()

	// Initialize repositories
	userRepo := repository.NewGormUserRepository(db)
	accountRepo := repository.NewGormAccountRepository(db)
	movementRepo := repository.NewGormMovementRepository(db)
	oauthTokenRepo := repository.NewGormOAuthTokenRepository(db)
	transferRepo := repository.NewGormTransferRepository(db)

	repos := repository.NewRepository(
		userRepo,
		accountRepo,
		movementRepo,
		oauthTokenRepo,
		transferRepo,
	)

	// Initialize OAuth client
	googleOAuth := oauth.NewGoogleOAuthClient(&cfg.OAuth.Google)

	// Initialize services
	authService := service.NewAuthService(
		repos.User,
		repos.Account,
		repos.OAuthToken,
		redisClient,
		googleOAuth,
		cfg,
	)

	accountService := service.NewAccountService(
		repos.Account,
		redisClient,
	)

	movementService := service.NewMovementService(
		repos.Movement,
		repos.Account,
		redisClient,
	)

	transferService := service.NewTransferService(
		repos.Transfer,
		repos.Account,
		repos.Movement,
		redisClient,
		db,
	)

	services := service.NewService(
		authService,
		accountService,
		movementService,
		transferService,
	)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(services.Auth)
	accountHandler := handler.NewAccountHandler(services.Account)
	movementHandler := handler.NewMovementHandler(services.Movement, services.Account)
	transferHandler := handler.NewTransferHandler(services.Transfer, services.Account)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(services.Auth, logger)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(redisClient, &cfg.Security.RateLimit, logger)

	// Setup router
	r := router.NewRouter(
		authHandler,
		accountHandler,
		movementHandler,
		transferHandler,
		authMiddleware,
		rateLimitMiddleware,
		logger,
	)

	// Initialize server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r.Setup(),
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.Timeout * 2,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server", zap.Int("port", cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Register metrics
	registerMetrics()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline to wait for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}

// Connect to the database
func connectToDatabase(cfg config.DBConfig, logger *zap.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.GetDBURL()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Force auto-migration to run regardless of config setting
	logger.Info("Starting database schema migration")
	if err := migrateDatabase(db, logger); err != nil {
		logger.Error("Database migration failed", zap.Error(err))
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	logger.Info("Database schema migration completed successfully")

	return db, nil
}

// migrateDatabase automatically creates or updates the database schema based on GORM models
func migrateDatabase(db *gorm.DB, logger *zap.Logger) error {
	// Enable PostgreSQL-specific extensions
	logger.Info("Setting up PostgreSQL extensions")
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		logger.Error("Failed to create uuid-ossp extension", zap.Error(err))
		return err
	}
	logger.Info("PostgreSQL extensions setup complete")

	// First check if tables exist
	var tableCount int
	if err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'").Scan(&tableCount).Error; err != nil {
		logger.Error("Failed to check existing tables", zap.Error(err))
		return fmt.Errorf("failed to check existing tables: %w", err)
	}

	logger.Info("Current database state", zap.Int("existing_tables", tableCount))

	// If tables already exist, verify they're the ones we need and skip migration
	if tableCount > 0 {
		tables := []string{"users", "accounts", "movements", "oauth_tokens", "transfers"}
		allTablesExist := true

		for _, table := range tables {
			var exists bool
			checkTableSQL := `
				SELECT EXISTS (
					SELECT FROM information_schema.tables 
					WHERE table_schema = 'public' 
					AND table_name = $1
				)
			`
			if err := db.Raw(checkTableSQL, table).Scan(&exists).Error; err != nil {
				logger.Error("Failed to check if table exists",
					zap.String("table", table),
					zap.Error(err))
				allTablesExist = false
				break
			} else if !exists {
				logger.Warn("Required table does not exist",
					zap.String("table", table))
				allTablesExist = false
				break
			} else {
				logger.Info("Table verified", zap.String("table", table))
			}
		}

		if allTablesExist {
			logger.Info("All required tables already exist, skipping migration")
			return nil
		} else {
			logger.Warn("Some required tables are missing, will attempt SQL migration")
		}
	}

	// If no tables exist or required tables are missing, try using SQL migrations
	logger.Info("Database requires setup, attempting to run SQL migrations")
	if err := runSQLMigrations(db, logger); err != nil {
		logger.Error("SQL migrations failed", zap.Error(err))
		return fmt.Errorf("failed to run SQL migrations: %w", err)
	}

	logger.Info("SQL migrations completed successfully")

	// Verify all tables were created
	tables := []string{"users", "accounts", "movements", "oauth_tokens", "transfers"}
	for _, table := range tables {
		var exists bool
		checkTableSQL := `
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = $1
			)
		`
		if err := db.Raw(checkTableSQL, table).Scan(&exists).Error; err != nil {
			logger.Error("Failed to check if table exists",
				zap.String("table", table),
				zap.Error(err))
			return fmt.Errorf("failed to verify table %s: %w", table, err)
		} else if !exists {
			logger.Error("Table does not exist after migration",
				zap.String("table", table))
			return fmt.Errorf("table %s was not created by migration", table)
		}

		logger.Info("Table verified after migration", zap.String("table", table))
	}

	logger.Info("Database migration complete")
	return nil
}

// runSQLMigrations runs the SQL migration files directly
func runSQLMigrations(db *gorm.DB, logger *zap.Logger) error {
	logger.Info("Running SQL migrations from files")

	// Run the initial schema migration script
	migrationPath := "migrations/000001_initial_schema.up.sql"
	logger.Info("Reading migration file", zap.String("path", migrationPath))

	// Read the migration file
	content, err := os.ReadFile(migrationPath)
	if err != nil {
		logger.Error("Failed to read migration file",
			zap.String("path", migrationPath),
			zap.Error(err))
		return err
	}

	// Execute the SQL script as a single batch
	logger.Info("Executing SQL migration script")
	if err := db.Exec(string(content)).Error; err != nil {
		logger.Error("Failed to execute SQL migration script", zap.Error(err))
		return err
	}

	logger.Info("SQL migrations completed successfully")
	return nil
}

// initGormLogger creates a GORM logger that integrates with our zap logger
func initGormLogger(logger *zap.Logger) gormlogger.Interface {
	zapLogger := logger.Named("gorm")
	return &GormToZapLogger{
		ZapLogger:                 zapLogger,
		LogLevel:                  gormlogger.Info,
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: false,
	}
}

// GormToZapLogger implements GORM's logger interface to use Zap for logging
type GormToZapLogger struct {
	ZapLogger                 *zap.Logger
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

// LogMode sets the log level
func (l *GormToZapLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info logs info messages
func (l *GormToZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.ZapLogger.Sugar().Infof(msg, data...)
	}
}

// Warn logs warn messages
func (l *GormToZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.ZapLogger.Sugar().Warnf(msg, data...)
	}
}

// Error logs error messages
func (l *GormToZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.ZapLogger.Sugar().Errorf(msg, data...)
	}
}

// Trace logs SQL messages with execution time
func (l *GormToZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
	}

	// Log query execution details
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		l.ZapLogger.Error("GORM error", append(fields, zap.Error(err))...)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		l.ZapLogger.Warn("GORM slow query", fields...)
	case l.LogLevel >= gormlogger.Info:
		l.ZapLogger.Debug("GORM query", fields...)
	}
}

// Register Prometheus metrics
func registerMetrics() {
	// Define metrics
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	dbOperationsDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_operation_duration_seconds",
			Help:    "Database operation duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// Register metrics
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(dbOperationsDuration)
}
