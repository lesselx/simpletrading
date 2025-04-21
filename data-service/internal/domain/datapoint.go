package domain

import "time"

type DataPoint struct {
	ID        uint      `gorm:"primaryKey"`
	Value     float64   `gorm:"not null"`
	Timestamp time.Time `gorm:"autoCreateTime"`
}
