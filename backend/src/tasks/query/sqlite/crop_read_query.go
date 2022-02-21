package sqlite

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/tasks/query"
)

type CropQuerySqlite struct {
	DB *sql.DB
}

func NewCropQuerySqlite(db *sql.DB) query.Crop {
	return CropQuerySqlite{DB: db}
}

func (s CropQuerySqlite) FindCropByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		rowsData := struct {
			UID     string
			BatchID string
		}{}
		crop := query.TaskCropResult{}

		s.DB.QueryRow(`SELECT UID, BATCH_ID
			FROM CROP_READ WHERE UID = ?`, uid).Scan(&rowsData.UID, &rowsData.BatchID)

		cropUID, err := uuid.FromString(rowsData.UID)
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
