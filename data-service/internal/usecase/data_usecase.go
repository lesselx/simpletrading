package usecase

import (
	"log"
	"simpletrading/dataservice/internal/domain"
	repository "simpletrading/dataservice/internal/repository/memory"
	"time"
)

// DataUsecase provides methods for working with data
type DataUsecase struct {
	repo repository.DataRepository
}

// NewDataUsecase initializes a new instance of DataUsecase
func NewDataUsecase(repo repository.DataRepository) *DataUsecase {
	return &DataUsecase{repo: repo}
}

// GenerateData generates a new data point and stores it
func (uc *DataUsecase) GenerateData(value float64) error {
	err := uc.repo.Insert(value)
	if err != nil {
		log.Println("Error saving data:", err)
		return err
	}
	log.Println("Generated and saved new data point:", value)
	return nil
}

// Get retrieves the most recent data points

func (uc *DataUsecase) GetRecentData(limit int) ([]domain.DataPoint, error) {
	data, err := uc.repo.GetRecent(limit)
	if err != nil {
		log.Println("Error retrieving recent data:", err)
		return nil, err
	}
	return data, nil
}

func (uc *DataUsecase) GetLowestPriceInLast24Hours() (float64, error) {
	// Filter data by timestamp for the last 24 hours
	oneDayAgo := time.Now().Add(-24 * time.Hour)
	data, err := uc.repo.GetDataSince(oneDayAgo)
	if err != nil {
		log.Println("Error fetching data:", err)
		return 0, err
	}

	// Find the lowest price
	var lowest float64
	for _, dp := range data {
		if lowest == 0 || dp.Value < lowest {
			lowest = dp.Value
		}
	}

	return lowest, nil
}
