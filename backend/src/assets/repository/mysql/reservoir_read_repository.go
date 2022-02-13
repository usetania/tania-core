package mysql

import (
	"database/sql"

	"github.com/usetania/tania-core/src/assets/repository"
	"github.com/usetania/tania-core/src/assets/storage"
)

type ReservoirReadRepositoryMysql struct {
	DB *sql.DB
}

func NewReservoirReadRepositoryMysql(db *sql.DB) repository.ReservoirRead {
	return &ReservoirReadRepositoryMysql{DB: db}
}

func (f *ReservoirReadRepositoryMysql) Save(reservoirRead *storage.ReservoirRead) <-chan error {
	result := make(chan error)

	go func() {
		count := 0

		err := f.DB.QueryRow(`SELECT COUNT(*) FROM RESERVOIR_READ WHERE UID = ?`, reservoirRead.UID.Bytes()).Scan(&count)
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
				reservoirRead.Farm.UID.Bytes(),
				reservoirRead.Farm.Name,
				reservoirRead.CreatedDate,
				reservoirRead.UID.Bytes())

			if err != nil {
				result <- err
			}

			if len(reservoirRead.Notes) > 0 {
				// Just delete them all then insert them all again.
				// We can refactor it later.
				_, err := f.DB.Exec(`DELETE FROM RESERVOIR_READ_NOTES WHERE RESERVOIR_UID = ?`, reservoirRead.UID.Bytes())
				if err != nil {
					result <- err
				}

				for _, v := range reservoirRead.Notes {
					_, err := f.DB.Exec(`INSERT INTO RESERVOIR_READ_NOTES (UID, RESERVOIR_UID, CONTENT, CREATED_DATE)
							VALUES (?, ?, ?, ?)`, v.UID.Bytes(), reservoirRead.UID.Bytes(), v.Content, v.CreatedDate)
					if err != nil {
						result <- err
					}
				}
			}
		} else {
			_, err = f.DB.Exec(`INSERT INTO RESERVOIR_READ
				(UID, NAME, WATERSOURCE_TYPE, WATERSOURCE_CAPACITY, FARM_UID, FARM_NAME, CREATED_DATE)
				VALUES (?, ?, ?, ?, ?, ?, ?)`,
				reservoirRead.UID.Bytes(),
				reservoirRead.Name,
				reservoirRead.WaterSource.Type,
				reservoirRead.WaterSource.Capacity,
				reservoirRead.Farm.UID.Bytes(),
				reservoirRead.Farm.Name,
				reservoirRead.CreatedDate)

			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}
