package domain

import "time"

type Trade struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    string    `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
