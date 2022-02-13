package sqlite

import (
	"database/sql"
	"time"

	"github.com/usetania/tania-core/src/growth/repository"
	"github.com/usetania/tania-core/src/growth/storage"
)

type CropReadRepositorySqlite struct {
	DB *sql.DB
}

func NewCropReadRepositorySqlite(db *sql.DB) repository.CropRead {
	return &CropReadRepositorySqlite{DB: db}
}

func (f *CropReadRepositorySqlite) Save(cropRead *storage.CropRead) <-chan error {
	result := make(chan error)

	go func() {
		count := 0

		err := f.DB.QueryRow(`SELECT COUNT(*) FROM CROP_READ WHERE UID = ?`, cropRead.UID).Scan(&count)
		if err != nil {
			result <- err
		}

		var initialAreaLastWatered string
		if cropRead.InitialArea.LastWatered != nil && !cropRead.InitialArea.LastWatered.IsZero() {
			initialAreaLastWatered = cropRead.InitialArea.LastWatered.Format(time.RFC3339)
		}

		var initialAreaLastFertilized string
		if cropRead.InitialArea.LastFertilized != nil && !cropRead.InitialArea.LastFertilized.IsZero() {
			initialAreaLastFertilized = cropRead.InitialArea.LastFertilized.Format(time.RFC3339)
		}

		var initialAreaLastPesticided string
		if cropRead.InitialArea.LastPesticided != nil && !cropRead.InitialArea.LastPesticided.IsZero() {
			initialAreaLastPesticided = cropRead.InitialArea.LastPesticided.Format(time.RFC3339)
		}

		var initialAreaLastPruned string
		if cropRead.InitialArea.LastPruned != nil && !cropRead.InitialArea.LastPruned.IsZero() {
			initialAreaLastPruned = cropRead.InitialArea.LastPruned.Format(time.RFC3339)
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
				cropRead.Inventory.UID,
				cropRead.Inventory.Type,
				cropRead.Inventory.PlantType,
				cropRead.Inventory.Name,
				cropRead.AreaStatus.Seeding,
				cropRead.AreaStatus.Growing,
				cropRead.AreaStatus.Dumped,
				cropRead.FarmUID,
				cropRead.InitialArea.AreaUID,
				cropRead.InitialArea.Name,
				cropRead.InitialArea.InitialQuantity,
				cropRead.InitialArea.CurrentQuantity,
				initialAreaLastWatered,
				initialAreaLastFertilized,
				initialAreaLastPesticided,
				initialAreaLastPruned,
				cropRead.InitialArea.CreatedDate.Format(time.RFC3339),
				cropRead.InitialArea.LastUpdated.Format(time.RFC3339),
				cropRead.UID)

			if err != nil {
				result <- err
			}

			if len(cropRead.Photos) > 0 {
				for _, v := range cropRead.Photos {
					res, err := f.DB.Exec(`UPDATE CROP_READ_PHOTO
						SET FILENAME = ?, MIMETYPE = ?, SIZE = ?,
						WIDTH = ?, HEIGHT = ?, DESCRIPTION = ?
						WHERE UID = ?`,
						v.Filename, v.MimeType, v.Size, v.Width, v.Height, v.Description, v.UID)
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
							v.UID, cropRead.UID, v.Filename, v.MimeType, v.Size, v.Width, v.Height, v.Description)

						if err != nil {
							result <- err
						}
					}
				}
			}

			if len(cropRead.MovedArea) > 0 {
				for _, v := range cropRead.MovedArea {
					var movedLastWatered string
					if v.LastWatered != nil && !v.LastWatered.IsZero() {
						movedLastWatered = v.LastWatered.Format(time.RFC3339)
					}

					var movedLastFertilized string
					if v.LastFertilized != nil && !v.LastFertilized.IsZero() {
						movedLastFertilized = v.LastFertilized.Format(time.RFC3339)
					}

					var movedLastPesticided string
					if v.LastPesticided != nil && !v.LastPesticided.IsZero() {
						movedLastPesticided = v.LastPesticided.Format(time.RFC3339)
					}

					var movedLastPruned string
					if v.LastPruned != nil && !v.LastPruned.IsZero() {
						movedLastPruned = v.LastPruned.Format(time.RFC3339)
					}

					cd := v.CreatedDate.Format(time.RFC3339)
					lu := v.LastUpdated.Format(time.RFC3339)

					res, err := f.DB.Exec(`UPDATE CROP_READ_MOVED_AREA
						SET NAME = ?, INITIAL_QUANTITY = ?, CURRENT_QUANTITY = ?,
						LAST_WATERED = ?, LAST_FERTILIZED = ?, LAST_PESTICIDED = ?, LAST_PRUNED = ?,
						CREATED_DATE = ?, LAST_UPDATED = ?
						WHERE CROP_UID = ? AND AREA_UID = ?`,
						v.Name, v.InitialQuantity, v.CurrentQuantity,
						movedLastWatered, movedLastFertilized, movedLastPesticided, movedLastPruned,
						cd, lu,
						cropRead.UID, v.AreaUID)
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
							cropRead.UID, v.AreaUID, v.Name, v.InitialQuantity, v.CurrentQuantity,
							movedLastWatered, movedLastFertilized, movedLastPesticided, movedLastPruned,
							cd, lu)

						if err != nil {
							result <- err
						}
					}
				}
			}

			if len(cropRead.HarvestedStorage) > 0 {
				for _, v := range cropRead.HarvestedStorage {
					cd := v.CreatedDate.Format(time.RFC3339)
					lu := v.LastUpdated.Format(time.RFC3339)

					res, err := f.DB.Exec(`UPDATE CROP_READ_HARVESTED_STORAGE
						SET QUANTITY = ?, PRODUCED_GRAM_QUANTITY = ?,
						SOURCE_AREA_NAME = ?,
						CREATED_DATE = ?, LAST_UPDATED = ?
						WHERE CROP_UID = ? AND SOURCE_AREA_UID = ?`,
						v.Quantity, v.ProducedGramQuantity,
						v.SourceAreaName,
						cd, lu,
						cropRead.UID, v.SourceAreaUID)
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
							cropRead.UID, v.Quantity, v.ProducedGramQuantity,
							v.SourceAreaUID, v.SourceAreaName, cd, lu)

						if err != nil {
							result <- err
						}
					}
				}
			}

			if len(cropRead.Trash) > 0 {
				for _, v := range cropRead.Trash {
					cd := v.CreatedDate.Format(time.RFC3339)
					lu := v.LastUpdated.Format(time.RFC3339)

					res, err := f.DB.Exec(`UPDATE CROP_READ_TRASH
						SET QUANTITY = ?, SOURCE_AREA_NAME = ?,
						CREATED_DATE = ?, LAST_UPDATED = ?
						WHERE CROP_UID = ? AND SOURCE_AREA_UID = ?`,
						v.Quantity, v.SourceAreaName, cd, lu,
						cropRead.UID, v.SourceAreaUID)
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
							cropRead.UID, v.Quantity, v.SourceAreaUID, v.SourceAreaName, cd, lu)

						if err != nil {
							result <- err
						}
					}
				}
			}

			if len(cropRead.Notes) > 0 {
				// Just delete them all then insert them all again.
				// We can refactor it later.
				_, err := f.DB.Exec(`DELETE FROM CROP_READ_NOTES WHERE CROP_UID = ?`, cropRead.UID)
				if err != nil {
					result <- err
				}

				for _, v := range cropRead.Notes {
					_, err := f.DB.Exec(`INSERT INTO CROP_READ_NOTES (UID, CROP_UID, CONTENT, CREATED_DATE)
							VALUES (?, ?, ?, ?)`, v.UID, cropRead.UID, v.Content, v.CreatedDate.Format(time.RFC3339))
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
				cropRead.UID,
				cropRead.BatchID,
				cropRead.Status,
				cropRead.Type,
				cropRead.Container.Quantity,
				cropRead.Container.Type,
				cropRead.Container.Cell,
				cropRead.Inventory.UID,
				cropRead.Inventory.Type,
				cropRead.Inventory.PlantType,
				cropRead.Inventory.Name,
				cropRead.AreaStatus.Seeding,
				cropRead.AreaStatus.Growing,
				cropRead.AreaStatus.Dumped,
				cropRead.FarmUID,
				cropRead.InitialArea.AreaUID,
				cropRead.InitialArea.Name,
				cropRead.InitialArea.InitialQuantity,
				cropRead.InitialArea.CurrentQuantity,
				initialAreaLastWatered,
				initialAreaLastFertilized,
				initialAreaLastPesticided,
				initialAreaLastPruned,
				cropRead.InitialArea.CreatedDate.Format(time.RFC3339),
				cropRead.InitialArea.LastUpdated.Format(time.RFC3339))

			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}
