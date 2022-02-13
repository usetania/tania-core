package sqlite

import (
	"database/sql"
	"time"

	"github.com/usetania/tania-core/src/assets/repository"
	"github.com/usetania/tania-core/src/assets/storage"
)

type FarmReadRepositorySqlite struct {
	DB *sql.DB
}

func NewFarmReadRepositorySqlite(db *sql.DB) repository.FarmRead {
	return &FarmReadRepositorySqlite{DB: db}
}

func (f *FarmReadRepositorySqlite) Save(farmRead *storage.FarmRead) <-chan error {
	result := make(chan error)

	go func() {
		count := 0

		err := f.DB.QueryRow(`SELECT COUNT(*) FROM FARM_READ WHERE UID = ?`, farmRead.UID).Scan(&count)
		if err != nil {
			result <- err
		}

		if count > 0 {
			_, err := f.DB.Exec(`UPDATE FARM_READ SET
				NAME = ?, LATITUDE = ?, LONGITUDE = ?, TYPE = ?, COUNTRY = ?, CITY = ?,
				IS_ACTIVE = ?, CREATED_DATE = ?
				WHERE UID = ?`,
				farmRead.Name, farmRead.Latitude, farmRead.Longitude, farmRead.Type,
				farmRead.Country, farmRead.City, farmRead.IsActive, farmRead.CreatedDate.Format(time.RFC3339),
				farmRead.UID)
			if err != nil {
				result <- err
			}
		} else {
			_, err := f.DB.Exec(`INSERT INTO FARM_READ
				(UID, NAME, LATITUDE, LONGITUDE, TYPE, COUNTRY, CITY, IS_ACTIVE, CREATED_DATE)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				farmRead.UID, farmRead.Name, farmRead.Latitude, farmRead.Longitude, farmRead.Type,
				farmRead.Country, farmRead.City, farmRead.IsActive, farmRead.CreatedDate.Format(time.RFC3339))
			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}
