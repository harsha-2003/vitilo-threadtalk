package config

import (
	"database/sql"
	"log"
	"time"

	"github.com/harsha-2003/vitilo-threadtalk/backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

func InitDB() *gorm.DB {
	// Open raw SQL connection first with pure Go driver
	sqlDB, err := sql.Open("sqlite", "threadtalk.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Wrap with GORM using existing connection
	db, err := gorm.Open(sqlite.Dialector{
		Conn: sqlDB,
	}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		log.Fatal("Failed to initialize GORM:", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	db.AutoMigrate(
		&models.User{},
		&models.Community{},
		&models.Post{},
		&models.Comment{},
		&models.Vote{},
	)

	log.Println("✓ Database connected successfully (Pure Go SQLite)")
	return db
}
