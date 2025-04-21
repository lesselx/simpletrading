package memory

import (
	"simpletrading/dataservice/internal/domain"
	"time"

	"gorm.io/gorm"
)

type DataRepository struct {
	db *gorm.DB
}

func NewDataRepo(db *gorm.DB) *DataRepository {
	return &DataRepository{db: db}
}

func (r *DataRepository) Insert(value float64) error {
	dp := domain.DataPoint{Value: value}
	return r.db.Create(&dp).Error
}

func (r *DataRepository) GetRecent(limit int) ([]domain.DataPoint, error) {
	var data []domain.DataPoint
	err := r.db.Order("timestamp desc").Limit(limit).Find(&data).Error
	return data, err
}

func (r *DataRepository) GetDataSince(startTime time.Time) ([]domain.DataPoint, error) {
	var data []domain.DataPoint
	err := r.db.Where("timestamp >= ?", startTime).Order("timestamp desc").Find(&data).Error
	return data, err
}
