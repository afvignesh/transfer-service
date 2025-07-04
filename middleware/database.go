package middleware

import (
	"database/sql"
	"fmt"
	"os"
	"go.uber.org/zap"
	_ "github.com/lib/pq"
)

// DatabaseMiddleware provides database connection to handlers
type DatabaseMiddleware struct {
	DB *sql.DB
}

// NewDatabaseMiddleware creates a new database middleware
func NewDatabaseMiddleware() (*DatabaseMiddleware, error) {
	log := GetLogger()
	
	// Get database connection string from environment or use default
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://user:password@localhost:5432/internal_transfer?sslmode=disable"
	}
	
	log.Info("Connecting to database", zap.String("dsn", maskPassword(dsn)))
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Error("Failed to open database connection", zap.Error(err))
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	// Test the connection
	if err := db.Ping(); err != nil {
		log.Error("Failed to ping database", zap.Error(err))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	log.Info("Database connection established successfully")
	
	return &DatabaseMiddleware{DB: db}, nil
}

// Close closes the database connection
func (dm *DatabaseMiddleware) Close() error {
	if dm.DB != nil {
		return dm.DB.Close()
	}
	return nil
}

// GetDB returns the database connection
func (dm *DatabaseMiddleware) GetDB() *sql.DB {
	return dm.DB
}

// maskPassword masks the password in the connection string for logging
func maskPassword(dsn string) string {
	// Simple masking - in production you might want more sophisticated masking
	if len(dsn) > 20 {
		return dsn[:20] + "***"
	}
	return "***"
} 