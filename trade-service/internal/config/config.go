package config

import (
	"database/sql"
	"log"
	"os"
	"simpletrading/tradeservice/internal/domain"

	"github.com/joho/godotenv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "modernc.org/sqlite"
)

type Config struct {
	DBPath       string
	Port         string
	JWTSecret    string
	ClientId     string // Machine ID for authentication
	ClientSecret string // Machine secret for authentication
	AuthUrl      string // URL for authentication
	DataUrl      string // URL for data service
}

func Init() (*gorm.DB, *Config) {
	// Load environment variables and set up DB and OAuth
	cfg := Load()
	db := InitDatabase(cfg.DBPath)

	return db, cfg
}

// Load reads the configuration values from environment variables (or defaults).
func Load() *Config {
	// Set default values

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Remaining logic stays the same

	cfg := &Config{
		DBPath:       "auth.db",  // Default SQLite DB path
		Port:         ":8081",    // Default server port
		JWTSecret:    "mysecret", // Default JWT secret key
		ClientId:     "myclientid",
		ClientSecret: "myclientsecret",
		AuthUrl:      "http://localhost:8080/auth/token",  // Default authentication URL
		DataUrl:      "http://localhost:8081/data/lowest", // Default data URL
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

	if clientId := os.Getenv("CLIENT_ID"); clientId != "" {
		cfg.ClientId = clientId
	}

	if clientSecret := os.Getenv("CLIENT_SECRET"); clientSecret != "" {

		cfg.ClientSecret = clientSecret

	}

	if authUrl := os.Getenv("AUTH_URL"); authUrl != "" {

		cfg.AuthUrl = authUrl

	}

	if dataUrl := os.Getenv("DATA_URL"); dataUrl != "" {

		cfg.DataUrl = dataUrl

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
	err = db.AutoMigrate(&domain.Trade{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	return db
}
