package usecase

import (
	"errors"
	"fmt"
	"simpletrading/authservice/internal/config"
	"simpletrading/authservice/internal/domain"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthUsecase struct {
	repo domain.UserRepository
	cfg  config.Config
}

func NewAuthUsecase(repo domain.UserRepository, cfg config.Config) *AuthUsecase {
	return &AuthUsecase{repo: repo, cfg: cfg}
}

func (uc *AuthUsecase) Register(user domain.User) error {
	return uc.repo.Create(user)
}

func (uc *AuthUsecase) Login(username, password string) (*domain.User, error) {
	user, err := uc.repo.GetByUsername(username)
	if err != nil || user.Password != password {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (uc *AuthUsecase) GetToken(clientID, clientSecret string) (string, error) {
	if clientID != uc.cfg.ClientId || clientSecret != uc.cfg.ClientSecret {
		return "", fmt.Errorf("invalid client credentials")
	}

	// Create a JWT token for machine authentication
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"client_id": clientID,
		"sub":       clientID, // Set sub claim to the client ID
		"exp":       time.Now().Add(1 * time.Hour).Unix(),
	})

	// Sign the token
	tokenString, err := token.SignedString([]byte(uc.cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("error signing token: %v", err)
	}

	return tokenString, nil
}
