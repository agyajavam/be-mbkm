package database

import (
	"context"
	"fmt"
	"log"
	"mbkm-api/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Pool *pgxpool.Pool
	GORM *gorm.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	ctx := context.Background()

	// Setup pgx pool for native SQL queries
	connString := cfg.DatabaseURL()
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = time.Minute * 30

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	// Setup GORM for auto-migration only
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("unable to connect GORM: %w", err)
	}

	log.Println("✅ Database connected successfully (pgx + GORM)")

	return &Database{
		Pool: pool,
		GORM: gormDB,
	}, nil
}

func (db *Database) Close() {
	db.Pool.Close()
}

func (db *Database) AutoMigrate(models ...interface{}) error {
	log.Println("🔄 Running GORM auto-migration...")

	if err := db.GORM.AutoMigrate(models...); err != nil {
		return fmt.Errorf("auto-migration failed: %w", err)
	}

	log.Println("✅ GORM auto-migration completed")
	return nil
}
