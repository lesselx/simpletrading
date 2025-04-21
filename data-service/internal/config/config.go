package config

import (
	"database/sql"
	"log"
	"os"
	"simpletrading/dataservice/internal/domain"

	"github.com/joho/godotenv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "modernc.org/sqlite"
)

type Config struct {
	DBPath    string
	Port      string
	JWTSecret string
}

func Init() (*gorm.DB, *Config) {
	// Load environment variables and set up DB and OAuth
	cfg := Load()
	db := InitDatabase(cfg.DBPath)
	// SeedDummyData(db) // Seed dummy data

	return db, cfg
}

// func SeedDummyData(db *gorm.DB) {
// 	now := time.Now().UTC()
// 	dummies := []domain.DataPoint{
// 		{Value: 1200, Timestamp: now.Add(-2 * time.Hour)},
// 		{Value: 950, Timestamp: now.Add(-3 * time.Hour)},
// 		{Value: 870, Timestamp: now.Add(-6 * time.Hour)},
// 		{Value: 600, Timestamp: now.Add(-25 * time.Hour)}, // older than 24h
// 	}

// 	for _, dp := range dummies {
// 		if err := db.Create(&dp).Error; err != nil {
// 			log.Println("Failed to insert dummy datapoint:", dp, err)
// 		}
// 	}
// }

// Load reads the configuration values from environment variables (or defaults).
func Load() *Config {
	// Set default values

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Remaining logic stays the same

	cfg := &Config{
		DBPath:    "auth.db",  // Default SQLite DB path
		Port:      ":8081",    // Default server port
		JWTSecret: "mysecret", // Default JWT secret key
	}

	// Override with environment variables if they exist
	if dbPath := os.Getenv("DB_PATH"); dbPath != "" {
		cfg.DBPath = dbPath
	}
	if port := os.Getenv("PORT"); port != "" {
		cfg.Port = port
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		cfg.JWTSecret = jwtSecret
	}

	return cfg
}

func InitDatabase(path string) *gorm.DB {
	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		log.Fatalf("Failed to open SQL database: %v", err)
	}

	db, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate User schema
	err = db.AutoMigrate(&domain.DataPoint{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	return db
}
