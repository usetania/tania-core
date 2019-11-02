package sqlite

import (
	"database/sql"

	"github.com/Tanibox/tania-core/src/tasks/query"
	uuid "github.com/satori/go.uuid"
)

type CropQuerySqlite struct {
	DB *sql.DB
}

func NewCropQuerySqlite(db *sql.DB) query.CropQuery {
	return CropQuerySqlite{DB: db}
}

func (s CropQuerySqlite) FindCropByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		rowsData := struct {
			UID     string
			BatchID string
		}{}
		crop := query.TaskCropQueryResult{}

		err := s.DB.QueryRow(`SELECT UID, BATCH_ID
			FROM CROP_READ WHERE UID = ?`, uid).Scan(&rowsData.UID, &rowsData.BatchID)

		cropUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		crop.UID = cropUID
		crop.BatchID = rowsData.BatchID

		result <- query.QueryResult{Result: crop}

		close(result)
	}()

	return result
}
