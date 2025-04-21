package usecase

import (
	"database/sql"
	"simpletrading/dataservice/internal/domain"
	repository "simpletrading/dataservice/internal/repository/memory"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Open in-memory SQLite DB
	sqlDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	db, err := gorm.Open(gorm.Dialector(sqlite.Dialector{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Migrate the schema to create the DataPoint table
	err = db.AutoMigrate(&domain.DataPoint{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	// Return the GORM DB object
	return db

}

func TestGetLowestPriceInLast24Hours(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewDataRepo(db)
	uc := NewDataUsecase(*repo)
	now := time.Now().UTC()

	dummies := []domain.DataPoint{
		{Value: 1200, Timestamp: now.Add(-2 * time.Hour)},
		{Value: 950, Timestamp: now.Add(-4 * time.Hour)},
		{Value: 870, Timestamp: now.Add(-5 * time.Hour)},
		{Value: 600, Timestamp: now.Add(-25 * time.Hour)}, // too old
	}

	for _, dp := range dummies {
		err := db.Create(&dp).Error
		if err != nil {
			t.Fatalf("failed to insert dummy data: %v", err)
		}
	}

	lowest, err := uc.GetLowestPriceInLast24Hours()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lowest != 870 {
		t.Errorf("expected lowest to be 870, got %v", lowest)
	}
}
