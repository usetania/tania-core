package sqlite

import (
	"database/sql"

	"github.com/usetania/tania-core/src/growth/repository"
	"github.com/usetania/tania-core/src/growth/storage"
)

type CropReadRepositoryMysql struct {
	DB *sql.DB
}

func NewCropReadRepositoryMysql(db *sql.DB) repository.CropRead {
	return &CropReadRepositoryMysql{DB: db}
}

func (f *CropReadRepositoryMysql) Save(cropRead *storage.CropRead) <-chan error {
	result := make(chan error)

	go func() {
		count := 0

		err := f.DB.QueryRow(`SELECT COUNT(*) FROM CROP_READ WHERE UID = ?`, cropRead.UID.Bytes()).Scan(&count)
		if err != nil {
			result <- err
		}

		if count > 0 {
			_, err = f.DB.Exec(`UPDATE CROP_READ SET
				BATCH_ID = ?, STATUS = ?, TYPE = ?,
				CONTAINER_QUANTITY = ?, CONTAINER_TYPE = ?, CONTAINER_CELL = ?,
				INVENTORY_UID = ?, INVENTORY_TYPE =?, INVENTORY_PLANT_TYPE = ?, INVENTORY_NAME = ?,
				AREA_STATUS_SEEDING = ?, AREA_STATUS_GROWING = ?, AREA_STATUS_DUMPED = ?,
				FARM_UID = ?,
				INITIAL_AREA_UID = ?, INITIAL_AREA_NAME = ?,
				INITIAL_AREA_INITIAL_QUANTITY = ?, INITIAL_AREA_CURRENT_QUANTITY = ?,
				INITIAL_AREA_LAST_WATERED = ?, INITIAL_AREA_LAST_FERTILIZED = ?,
				INITIAL_AREA_LAST_PESTICIDED = ?, INITIAL_AREA_LAST_PRUNED = ?,
				INITIAL_AREA_CREATED_DATE = ?, INITIAL_AREA_LAST_UPDATED = ?
				WHERE UID = ?`,
				cropRead.BatchID,
				cropRead.Status,
				cropRead.Type,
				cropRead.Container.Quantity,
				cropRead.Container.Type,
				cropRead.Container.Cell,
				cropRead.Inventory.UID.Bytes(),
				cropRead.Inventory.Type,
				cropRead.Inventory.PlantType,
				cropRead.Inventory.Name,
				cropRead.AreaStatus.Seeding,
				cropRead.AreaStatus.Growing,
				cropRead.AreaStatus.Dumped,
				cropRead.FarmUID.Bytes(),
				cropRead.InitialArea.AreaUID.Bytes(),
				cropRead.InitialArea.Name,
				cropRead.InitialArea.InitialQuantity,
				cropRead.InitialArea.CurrentQuantity,
				cropRead.InitialArea.LastWatered,
				cropRead.InitialArea.LastFertilized,
				cropRead.InitialArea.LastPesticided,
				cropRead.InitialArea.LastPruned,
				cropRead.InitialArea.CreatedDate,
				cropRead.InitialArea.LastUpdated,
				cropRead.UID.Bytes())

			if err != nil {
				result <- err
			}

			if len(cropRead.Photos) > 0 {
				for _, v := range cropRead.Photos {
					res, err := f.DB.Exec(`UPDATE CROP_READ_PHOTO
						SET FILENAME = ?, MIMETYPE = ?, SIZE = ?,
						WIDTH = ?, HEIGHT = ?, DESCRIPTION = ?
						WHERE UID = ?`,
						v.Filename, v.MimeType, v.Size, v.Width, v.Height, v.Description, v.UID.Bytes())
					if err != nil {
						result <- err
					}

					rowsAffected, err := res.RowsAffected()
					if err != nil {
						result <- err
					}

					if rowsAffected == 0 {
						f.DB.Exec(`INSERT INTO CROP_READ_PHOTO (
							UID, CROP_UID, FILENAME, MIMETYPE, SIZE, WIDTH, HEIGHT, DESCRIPTION)
							VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
							v.UID.Bytes(), cropRead.UID.Bytes(), v.Filename, v.MimeType, v.Size, v.Width, v.Height, v.Description)

						if err != nil {
							result <- err
						}
					}
				}
			}

			if len(cropRead.MovedArea) > 0 {
				for _, v := range cropRead.MovedArea {
					res, err := f.DB.Exec(`UPDATE CROP_READ_MOVED_AREA
						SET NAME = ?, INITIAL_QUANTITY = ?, CURRENT_QUANTITY = ?,
						LAST_WATERED = ?, LAST_FERTILIZED = ?, LAST_PESTICIDED = ?, LAST_PRUNED = ?,
						CREATED_DATE = ?, LAST_UPDATED = ?
						WHERE CROP_UID = ? AND AREA_UID = ?`,
						v.Name, v.InitialQuantity, v.CurrentQuantity,
						v.LastWatered, v.LastFertilized, v.LastPesticided, v.LastPruned,
						v.CreatedDate, v.LastUpdated,
						cropRead.UID.Bytes(), v.AreaUID.Bytes())
					if err != nil {
						result <- err
					}

					rowsAffected, err := res.RowsAffected()
					if err != nil {
						result <- err
					}

					if rowsAffected == 0 {
						_, err = f.DB.Exec(`INSERT INTO CROP_READ_MOVED_AREA (
							CROP_UID, AREA_UID, NAME, INITIAL_QUANTITY, CURRENT_QUANTITY,
							LAST_WATERED, LAST_FERTILIZED, LAST_PESTICIDED, LAST_PRUNED,
							CREATED_DATE, LAST_UPDATED)
							VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
							cropRead.UID.Bytes(), v.AreaUID.Bytes(), v.Name, v.InitialQuantity, v.CurrentQuantity,
							v.LastWatered, v.LastFertilized, v.LastPesticided, v.LastPruned,
							v.CreatedDate, v.LastUpdated)

						if err != nil {
							result <- err
						}
					}
				}
			}

			if len(cropRead.HarvestedStorage) > 0 {
				for _, v := range cropRead.HarvestedStorage {
					res, err := f.DB.Exec(`UPDATE CROP_READ_HARVESTED_STORAGE
						SET QUANTITY = ?, PRODUCED_GRAM_QUANTITY = ?,
						SOURCE_AREA_NAME = ?,
						CREATED_DATE = ?, LAST_UPDATED = ?
						WHERE CROP_UID = ? AND SOURCE_AREA_UID = ?`,
						v.Quantity, v.ProducedGramQuantity,
						v.SourceAreaName,
						v.CreatedDate, v.LastUpdated,
						cropRead.UID.Bytes(), v.SourceAreaUID.Bytes())
					if err != nil {
						result <- err
					}

					rowsAffected, err := res.RowsAffected()
					if err != nil {
						result <- err
					}

					if rowsAffected == 0 {
						_, err = f.DB.Exec(`INSERT INTO CROP_READ_HARVESTED_STORAGE (
							CROP_UID, QUANTITY, PRODUCED_GRAM_QUANTITY,
							SOURCE_AREA_UID, SOURCE_AREA_NAME,
							CREATED_DATE, LAST_UPDATED)
							VALUES (?, ?, ?, ?, ?, ?, ?)`,
							cropRead.UID.Bytes(), v.Quantity, v.ProducedGramQuantity,
							v.SourceAreaUID.Bytes(), v.SourceAreaName, v.CreatedDate, v.LastUpdated)

						if err != nil {
							result <- err
						}
					}
				}
			}

			if len(cropRead.Trash) > 0 {
				for _, v := range cropRead.Trash {
					res, err := f.DB.Exec(`UPDATE CROP_READ_TRASH
						SET QUANTITY = ?, SOURCE_AREA_NAME = ?,
						CREATED_DATE = ?, LAST_UPDATED = ?
						WHERE CROP_UID = ? AND SOURCE_AREA_UID = ?`,
						v.Quantity, v.SourceAreaName, v.CreatedDate, v.LastUpdated,
						cropRead.UID.Bytes(), v.SourceAreaUID.Bytes())
					if err != nil {
						result <- err
					}

					rowsAffected, err := res.RowsAffected()
					if err != nil {
						result <- err
					}

					if rowsAffected == 0 {
						_, err = f.DB.Exec(`INSERT INTO CROP_READ_TRASH (
							CROP_UID, QUANTITY, SOURCE_AREA_UID, SOURCE_AREA_NAME,
							CREATED_DATE, LAST_UPDATED)
							VALUES (?, ?, ?, ?, ?, ?)`,
							cropRead.UID.Bytes(), v.Quantity,
							v.SourceAreaUID.Bytes(), v.SourceAreaName, v.CreatedDate, v.LastUpdated)

						if err != nil {
							result <- err
						}
					}
				}
			}

			if len(cropRead.Notes) > 0 {
				// Just delete them all then insert them all again.
				// We can refactor it later.
				_, err := f.DB.Exec(`DELETE FROM CROP_READ_NOTES WHERE CROP_UID = ?`, cropRead.UID.Bytes())
				if err != nil {
					result <- err
				}

				for _, v := range cropRead.Notes {
					_, err := f.DB.Exec(`INSERT INTO CROP_READ_NOTES (UID, CROP_UID, CONTENT, CREATED_DATE)
							VALUES (?, ?, ?, ?)`, v.UID.Bytes(), cropRead.UID.Bytes(), v.Content, v.CreatedDate)
					if err != nil {
						result <- err
					}
				}
			}
		} else {
			_, err = f.DB.Exec(`INSERT INTO CROP_READ
				(UID, BATCH_ID, STATUS, TYPE, CONTAINER_QUANTITY, CONTAINER_TYPE, CONTAINER_CELL,
				INVENTORY_UID, INVENTORY_TYPE, INVENTORY_PLANT_TYPE, INVENTORY_NAME,
				AREA_STATUS_SEEDING, AREA_STATUS_GROWING, AREA_STATUS_DUMPED,
				FARM_UID,
				INITIAL_AREA_UID, INITIAL_AREA_NAME,
				INITIAL_AREA_INITIAL_QUANTITY, INITIAL_AREA_CURRENT_QUANTITY,
				INITIAL_AREA_LAST_WATERED, INITIAL_AREA_LAST_FERTILIZED, INITIAL_AREA_LAST_PESTICIDED,
				INITIAL_AREA_LAST_PRUNED, INITIAL_AREA_CREATED_DATE, INITIAL_AREA_LAST_UPDATED)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				cropRead.UID.Bytes(),
				cropRead.BatchID,
				cropRead.Status,
				cropRead.Type,
				cropRead.Container.Quantity,
				cropRead.Container.Type,
				cropRead.Container.Cell,
				cropRead.Inventory.UID.Bytes(),
				cropRead.Inventory.Type,
				cropRead.Inventory.PlantType,
				cropRead.Inventory.Name,
				cropRead.AreaStatus.Seeding,
				cropRead.AreaStatus.Growing,
				cropRead.AreaStatus.Dumped,
				cropRead.FarmUID.Bytes(),
				cropRead.InitialArea.AreaUID.Bytes(),
				cropRead.InitialArea.Name,
				cropRead.InitialArea.InitialQuantity,
				cropRead.InitialArea.CurrentQuantity,
				cropRead.InitialArea.LastWatered,
				cropRead.InitialArea.LastFertilized,
				cropRead.InitialArea.LastPesticided,
				cropRead.InitialArea.LastPruned,
				cropRead.InitialArea.CreatedDate,
				cropRead.InitialArea.LastUpdated)

			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}
