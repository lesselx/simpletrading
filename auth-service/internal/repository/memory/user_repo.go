package memory

import (
	"simpletrading/authservice/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user domain.User) error {
	result := r.db.Create(&user)
	return result.Error
}

func (r *UserRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	result := r.db.First(&user, "username = ?", username)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
