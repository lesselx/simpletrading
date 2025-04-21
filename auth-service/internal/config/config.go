package config

import (
	"database/sql"
	"log"
	"os"
	"simpletrading/authservice/internal/domain"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "modernc.org/sqlite"
)

var GoogleOAuthConfig *oauth2.Config

type Config struct {
	DBPath             string
	Port               string
	JWTSecret          string
	ClientId           string // Machine ID for authentication
	ClientSecret       string // Machine secret for authentication
	GoogleClientId     string // Google OAuth2 Client ID
	GoogleClientSecret string // Google OAuth2 Client Secret
	GoogleRedirectURI  string // Google OAuth2 Redirect URI
}

func Init() (*gorm.DB, *Config) {
	// Load environment variables and set up DB and OAuth
	cfg := Load()
	db := InitDatabase(*cfg)
	initGoogleOAuth(*cfg)

	return db, cfg
}

// Load reads the configuration values from environment variables (or defaults).
func Load() *Config {
	// Set default values

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	cfg := &Config{
		DBPath:       "auth.db",  // Default SQLite DB path
		Port:         ":8080",    // Default server port
		JWTSecret:    "mysecret", // Default JWT secret key
		ClientId:     "myclientid",
		ClientSecret: "myclientsecret",
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

	if googleClientId := os.Getenv("GOOGLE_CLIENT_ID"); googleClientId != "" {
		cfg.GoogleClientId = googleClientId
	}

	if googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET"); googleClientSecret != "" {
		cfg.GoogleClientSecret = googleClientSecret
	}

	if googleRedirectURI := os.Getenv("GOOGLE_REDIRECT_URI"); googleRedirectURI != "" {
		cfg.GoogleRedirectURI = googleRedirectURI
	}

	return cfg
}

func GenerateJWT(email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "mysecret"
	}

	claims := jwt.MapClaims{
		"sub": email,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func InitDatabase(cfg Config) *gorm.DB {
	sqlDB, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to open SQL database: %v", err)
	}

	db, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate User schema
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	return db
}

func initGoogleOAuth(cfg Config) {

	// Check if essential environment variables are missing
	if cfg.GoogleClientId == "" || cfg.GoogleClientSecret == "" || cfg.GoogleClientSecret == "" {
		log.Fatalf("Error: Missing Google OAuth2 environment variables (CLIENT_ID, CLIENT_SECRET, or REDIRECT_URI).")
	}

	//https://developers.google.com/identity/protocols/oauth2/scopes#oauth2
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     cfg.GoogleClientId,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURI,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	log.Println("Google OAuth2 config initialized successfully")
}
