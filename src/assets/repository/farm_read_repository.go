package repository

import (
	"database/sql"
	"time"

	"github.com/Tanibox/tania-server/src/assets/storage"
)

type FarmReadRepository interface {
	Save(farmRead *storage.FarmRead) <-chan error
}

type FarmReadRepositoryInMemory struct {
	Storage *storage.FarmReadStorage
}

func NewFarmReadRepositoryInMemory(s *storage.FarmReadStorage) FarmReadRepository {
	return &FarmReadRepositoryInMemory{Storage: s}
}

type FarmReadRepositorySqlite struct {
	DB *sql.DB
}

func NewFarmReadRepositorySqlite(db *sql.DB) FarmReadRepository {
	return &FarmReadRepositorySqlite{DB: db}
}

// Save is to save
func (f *FarmReadRepositoryInMemory) Save(farmRead *storage.FarmRead) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		f.Storage.FarmReadMap[farmRead.UID] = *farmRead

		result <- nil

		close(result)
	}()

	return result
}

func (f *FarmReadRepositorySqlite) Save(farmRead *storage.FarmRead) <-chan error {
	result := make(chan error)

	go func() {
		stmt, err := f.DB.Prepare(`INSERT INTO FARM_READ
			(UID, NAME, LATITUDE, LONGITUDE, TYPE, COUNTRY_CODE, CITY_CODE, IS_ACTIVE, CREATED_DATE)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)

		if err != nil {
			result <- err
			close(result)
		}

		_, err = stmt.Exec(farmRead.UID, farmRead.Name, farmRead.Latitude, farmRead.Longitude, farmRead.Type,
			farmRead.CountryCode, farmRead.CityCode, farmRead.IsActive, farmRead.CreatedDate.Format(time.RFC3339))

		if err != nil {
			result <- err
			close(result)
		}

		result <- err
		close(result)
	}()

	return result
}
