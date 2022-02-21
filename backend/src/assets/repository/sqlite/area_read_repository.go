package sqlite

import (
	"database/sql"
	"time"

	"github.com/usetania/tania-core/src/assets/repository"
	"github.com/usetania/tania-core/src/assets/storage"
)

type AreaReadRepositorySqlite struct {
	DB *sql.DB
}

func NewAreaReadRepositorySqlite(db *sql.DB) repository.AreaRead {
	return &AreaReadRepositorySqlite{DB: db}
}

func (f *AreaReadRepositorySqlite) Save(areaRead *storage.AreaRead) <-chan error {
	result := make(chan error)

	go func() {
		count := 0

		err := f.DB.QueryRow(`SELECT COUNT(*) FROM AREA_READ WHERE UID = ?`, areaRead.UID).Scan(&count)
		if err != nil {
			result <- err
		}

		if count > 0 {
			_, err := f.DB.Exec(`UPDATE AREA_READ SET
				NAME = ?, SIZE_UNIT = ?, SIZE = ?, TYPE = ?, LOCATION = ?,
				PHOTO_FILENAME = ?, PHOTO_MIMETYPE = ?, PHOTO_SIZE = ?, PHOTO_WIDTH = ?, PHOTO_HEIGHT = ?,
				CREATED_DATE = ?, FARM_UID = ?, FARM_NAME = ?, RESERVOIR_UID = ?, RESERVOIR_NAME = ?
				WHERE UID = ?`,
				areaRead.Name, areaRead.Size.Unit.Symbol, areaRead.Size.Value, areaRead.Type,
				areaRead.Location.Code, areaRead.Photo.Filename, areaRead.Photo.MimeType,
				areaRead.Photo.Size, areaRead.Photo.Width, areaRead.Photo.Height, areaRead.CreatedDate.Format(time.RFC3339),
				areaRead.Farm.UID, areaRead.Farm.Name, areaRead.Reservoir.UID, areaRead.Reservoir.Name, areaRead.UID)
			if err != nil {
				result <- err
			}

			if len(areaRead.Notes) > 0 {
				// Just delete them all then insert them all again.
				// We can refactor it later.
				_, err := f.DB.Exec(`DELETE FROM AREA_READ_NOTES WHERE AREA_UID = ?`, areaRead.UID)
				if err != nil {
					result <- err
				}

				for _, v := range areaRead.Notes {
					_, err := f.DB.Exec(`INSERT INTO AREA_READ_NOTES (UID, AREA_UID, CONTENT, CREATED_DATE)
							VALUES (?, ?, ?, ?)`, v.UID, areaRead.UID, v.Content, v.CreatedDate.Format(time.RFC3339))
					if err != nil {
						result <- err
					}
				}
			}
		} else {
			_, err := f.DB.Exec(`INSERT INTO AREA_READ
				(UID, NAME, SIZE_UNIT, SIZE, TYPE, LOCATION, PHOTO_FILENAME, PHOTO_MIMETYPE,
				PHOTO_SIZE, PHOTO_WIDTH, PHOTO_HEIGHT, CREATED_DATE, FARM_UID, FARM_NAME, RESERVOIR_UID, RESERVOIR_NAME)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				areaRead.UID, areaRead.Name, areaRead.Size.Unit.Symbol, areaRead.Size.Value, areaRead.Type,
				areaRead.Location.Code, areaRead.Photo.Filename, areaRead.Photo.MimeType,
				areaRead.Photo.Size, areaRead.Photo.Width, areaRead.Photo.Height, areaRead.CreatedDate.Format(time.RFC3339),
				areaRead.Farm.UID, areaRead.Farm.Name, areaRead.Reservoir.UID, areaRead.Reservoir.Name)
			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}
