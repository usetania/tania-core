package sqlite

import (
	"database/sql"

	"github.com/Tanibox/tania-server/src/growth/query"
	uuid "github.com/satori/go.uuid"
)

type TaskReadQueryMysql struct {
	DB *sql.DB
}

func NewTaskReadQueryMysql(db *sql.DB) query.TaskReadQuery {
	return TaskReadQueryMysql{DB: db}
}

type taskReadResult struct {
	UID         []byte
	Title       string
	Description string
	Category    string
	Status      string
	Domain      string
	CropID      []byte
	AreaID      []byte
	MaterialID  []byte
}

func (s TaskReadQueryMysql) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		taskQueryResult := query.CropTaskQueryResult{}
		rowsData := taskReadResult{}

		err := s.DB.QueryRow(`SELECT UID, TITLE, DESCRIPTION, CATEGORY, STATUS, DOMAIN_CODE,
			DOMAIN_DATA_CROP_ID, DOMAIN_DATA_AREA_ID, DOMAIN_DATA_MATERIAL_ID
			FROM TASK_READ WHERE UID = ?`, uid.Bytes()).Scan(
			&rowsData.UID,
			&rowsData.Title,
			&rowsData.Description,
			&rowsData.Category,
			&rowsData.Status,
			&rowsData.Domain,
			&rowsData.CropID,
			&rowsData.AreaID,
			&rowsData.MaterialID,
		)

		if err != nil && err != sql.ErrNoRows {
			result <- query.QueryResult{Error: err}
		}

		if err == sql.ErrNoRows {
			result <- query.QueryResult{Result: taskQueryResult}
		}

		taskUID, err := uuid.FromBytes(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		cropUID, err := uuid.FromBytes(rowsData.CropID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		areaUID, err := uuid.FromBytes(rowsData.AreaID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		materialUID, err := uuid.FromBytes(rowsData.MaterialID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		taskQueryResult.UID = taskUID
		taskQueryResult.Title = rowsData.Title
		taskQueryResult.Description = rowsData.Description
		taskQueryResult.Status = rowsData.Status
		taskQueryResult.Category = rowsData.Category
		taskQueryResult.Domain = rowsData.Domain
		taskQueryResult.CropUID = cropUID
		taskQueryResult.AreaUID = areaUID
		taskQueryResult.MaterialUID = materialUID

		result <- query.QueryResult{Result: taskQueryResult}
		close(result)
	}()

	return result
}
