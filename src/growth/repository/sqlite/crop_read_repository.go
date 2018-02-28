package sqlite

import (
	"database/sql"
	"time"

	"github.com/Tanibox/tania-server/src/growth/repository"
	"github.com/Tanibox/tania-server/src/growth/storage"
)

type CropReadRepositorySqlite struct {
	DB *sql.DB
}

func NewCropReadRepositorySqlite(db *sql.DB) repository.CropReadRepository {
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
				INVENTORY_UID = ?, INVENTORY_PLANT_TYPE = ?, INVENTORY_NAME = ?,
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

		} else {
			_, err = f.DB.Exec(`INSERT INTO CROP_READ
				(UID, BATCH_ID, STATUS, TYPE, CONTAINER_QUANTITY, CONTAINER_TYPE, CONTAINER_CELL,
				INVENTORY_UID, INVENTORY_PLANT_TYPE, INVENTORY_NAME,
				AREA_STATUS_SEEDING, AREA_STATUS_GROWING, AREA_STATUS_DUMPED,
				FARM_UID,
				INITIAL_AREA_UID, INITIAL_AREA_NAME,
				INITIAL_AREA_INITIAL_QUANTITY, INITIAL_AREA_CURRENT_QUANTITY,
				INITIAL_AREA_LAST_WATERED, INITIAL_AREA_LAST_FERTILIZED, INITIAL_AREA_LAST_PESTICIDED,
				INITIAL_AREA_LAST_PRUNED, INITIAL_AREA_CREATED_DATE, INITIAL_AREA_LAST_UPDATED)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				cropRead.UID,
				cropRead.BatchID,
				cropRead.Status,
				cropRead.Type,
				cropRead.Container.Quantity,
				cropRead.Container.Type,
				cropRead.Container.Cell,
				cropRead.Inventory.UID,
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
