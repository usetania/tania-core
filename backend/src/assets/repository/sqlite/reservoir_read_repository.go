package sqlite

import (
	"database/sql"
	"time"

	"github.com/usetania/tania-core/src/assets/repository"
	"github.com/usetania/tania-core/src/assets/storage"
)

type ReservoirReadRepositorySqlite struct {
	DB *sql.DB
}

func NewReservoirReadRepositorySqlite(db *sql.DB) repository.ReservoirRead {
	return &ReservoirReadRepositorySqlite{DB: db}
}

func (f *ReservoirReadRepositorySqlite) Save(reservoirRead *storage.ReservoirRead) <-chan error {
	result := make(chan error)

	go func() {
		count := 0

		err := f.DB.QueryRow(`SELECT COUNT(*) FROM RESERVOIR_READ WHERE UID = ?`, reservoirRead.UID).Scan(&count)
		if err != nil {
			result <- err
		}

		if count > 0 {
			_, err = f.DB.Exec(`UPDATE RESERVOIR_READ SET
				NAME = ?, WATERSOURCE_TYPE = ?, WATERSOURCE_CAPACITY = ?, FARM_UID = ?,
				FARM_NAME = ?, CREATED_DATE = ?
				WHERE UID = ?`,
				reservoirRead.Name,
				reservoirRead.WaterSource.Type,
				reservoirRead.WaterSource.Capacity,
				reservoirRead.Farm.UID,
				reservoirRead.Farm.Name,
				reservoirRead.CreatedDate.Format(time.RFC3339),
				reservoirRead.UID)

			if err != nil {
				result <- err
			}

			if len(reservoirRead.Notes) > 0 {
				// Just delete them all then insert them all again.
				// We can refactor it later.
				_, err := f.DB.Exec(`DELETE FROM RESERVOIR_READ_NOTES WHERE RESERVOIR_UID = ?`, reservoirRead.UID)
				if err != nil {
					result <- err
				}

				for _, v := range reservoirRead.Notes {
					_, err := f.DB.Exec(`INSERT INTO RESERVOIR_READ_NOTES (UID, RESERVOIR_UID, CONTENT, CREATED_DATE)
							VALUES (?, ?, ?, ?)`, v.UID, reservoirRead.UID, v.Content, v.CreatedDate.Format(time.RFC3339))
					if err != nil {
						result <- err
					}
				}
			}
		} else {
			_, err = f.DB.Exec(`INSERT INTO RESERVOIR_READ
				(UID, NAME, WATERSOURCE_TYPE, WATERSOURCE_CAPACITY, FARM_UID, FARM_NAME, CREATED_DATE)
				VALUES (?, ?, ?, ?, ?, ?, ?)`,
				reservoirRead.UID,
				reservoirRead.Name,
				reservoirRead.WaterSource.Type,
				reservoirRead.WaterSource.Capacity,
				reservoirRead.Farm.UID,
				reservoirRead.Farm.Name,
				reservoirRead.CreatedDate.Format(time.RFC3339))

			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}
