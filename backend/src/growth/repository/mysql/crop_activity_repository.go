package sqlite

import (
	"database/sql"
	"encoding/json"

	"github.com/usetania/tania-core/src/growth/decoder"
	"github.com/usetania/tania-core/src/growth/repository"
	"github.com/usetania/tania-core/src/growth/storage"
)

type CropActivityRepositoryMysql struct {
	DB *sql.DB
}

func NewCropActivityRepositoryMysql(db *sql.DB) repository.CropActivity {
	return &CropActivityRepositoryMysql{DB: db}
}

func (f *CropActivityRepositoryMysql) Save(cropActivity *storage.CropActivity, isUpdate bool) <-chan error {
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
				cropActivity.CreatedDate,
				cropActivity.Description,
				cropActivity.UID.Bytes())

			if err != nil {
				result <- err
			}
		} else {
			_, err = f.DB.Exec(`INSERT INTO CROP_ACTIVITY
				(CROP_UID, BATCH_ID, CONTAINER_TYPE, ACTIVITY_TYPE, ACTIVITY_TYPE_CODE, CREATED_DATE, DESCRIPTION)
				VALUES (?, ?, ?, ?, ?, ?, ?)`,
				cropActivity.UID.Bytes(),
				cropActivity.BatchID,
				cropActivity.ContainerType,
				at,
				cropActivity.ActivityType.Code(),
				cropActivity.CreatedDate,
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
