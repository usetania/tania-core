package mysql

import (
	"database/sql"

	"github.com/Tanibox/tania-core/src/tasks/query"
	uuid "github.com/satori/go.uuid"
)

type CropQueryMysql struct {
	DB *sql.DB
}

func NewCropQueryMysql(db *sql.DB) query.CropQuery {
	return CropQueryMysql{DB: db}
}

func (s CropQueryMysql) FindCropByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		rowsData := struct {
			UID     []byte
			BatchID string
		}{}
		crop := query.TaskCropQueryResult{}

		err := s.DB.QueryRow(`SELECT UID, BATCH_ID
			FROM CROP_READ WHERE UID = ?`, uid.Bytes()).Scan(&rowsData.UID, &rowsData.BatchID)

		cropUID, err := uuid.FromBytes(rowsData.UID)
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
