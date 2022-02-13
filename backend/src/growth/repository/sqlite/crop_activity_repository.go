package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/usetania/tania-core/src/growth/decoder"
	"github.com/usetania/tania-core/src/growth/repository"
	"github.com/usetania/tania-core/src/growth/storage"
)

type CropActivityRepositorySqlite struct {
	DB *sql.DB
}

func NewCropActivityRepositorySqlite(db *sql.DB) repository.CropActivity {
	return &CropActivityRepositorySqlite{DB: db}
}

func (f *CropActivityRepositorySqlite) Save(cropActivity *storage.CropActivity, isUpdate bool) <-chan error {
	result := make(chan error)

	go func() {
		at, err := json.Marshal(decoder.InterfaceWrapper{
			Name: cropActivity.ActivityType.Code(),
			Data: cropActivity.ActivityType,
		})
		if err != nil {
			result <- err
		}

		if isUpdate {
			_, err = f.DB.Exec(`UPDATE CROP_ACTIVITY
				SET BATCH_ID = ?, CONTAINER_TYPE = ?, ACTIVITY_TYPE = ?, ACTIVITY_TYPE_CODE = ?,
				CREATED_DATE = ?, DESCRIPTION = ?
				WHERE CROP_UID = ?`,
				cropActivity.BatchID,
				cropActivity.ContainerType,
				at,
				cropActivity.ActivityType.Code(),
				cropActivity.CreatedDate.Format(time.RFC3339),
				cropActivity.Description,
				cropActivity.UID)

			if err != nil {
				result <- err
			}
		} else {
			_, err = f.DB.Exec(`INSERT INTO CROP_ACTIVITY
				(CROP_UID, BATCH_ID, CONTAINER_TYPE, ACTIVITY_TYPE, ACTIVITY_TYPE_CODE, CREATED_DATE, DESCRIPTION)
				VALUES (?, ?, ?, ?, ?, ?, ?)`,
				cropActivity.UID,
				cropActivity.BatchID,
				cropActivity.ContainerType,
				at,
				cropActivity.ActivityType.Code(),
				cropActivity.CreatedDate.Format(time.RFC3339),
				cropActivity.Description)

			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}
