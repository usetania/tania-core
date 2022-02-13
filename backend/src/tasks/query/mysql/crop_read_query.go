package mysql

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/tasks/query"
)

type CropQueryMysql struct {
	DB *sql.DB
}

func NewCropQueryMysql(db *sql.DB) query.Crop {
	return CropQueryMysql{DB: db}
}

func (s CropQueryMysql) FindCropByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		rowsData := struct {
			UID     []byte
			BatchID string
		}{}
		crop := query.TaskCropResult{}

		s.DB.QueryRow(`SELECT UID, BATCH_ID
			FROM CROP_READ WHERE UID = ?`, uid.Bytes()).Scan(&rowsData.UID, &rowsData.BatchID)

		cropUID, err := uuid.FromBytes(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		crop.UID = cropUID
		crop.BatchID = rowsData.BatchID

		result <- query.Result{Result: crop}

		close(result)
	}()

	return result
}
