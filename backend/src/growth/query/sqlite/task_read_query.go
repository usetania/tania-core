package sqlite

import (
	"database/sql"

	"github.com/Tanibox/tania-core/src/growth/query"
	uuid "github.com/satori/go.uuid"
)

type TaskReadQuerySqlite struct {
	DB *sql.DB
}

func NewTaskReadQuerySqlite(db *sql.DB) query.TaskReadQuery {
	return TaskReadQuerySqlite{DB: db}
}

type taskReadResult struct {
	UID         string
	Title       string
	Description string
	Category    string
	Status      string
	Domain      string
	AssetID     string
	AreaID      sql.NullString
	MaterialID  sql.NullString
}

func (s TaskReadQuerySqlite) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		taskQueryResult := query.CropTaskQueryResult{}
		rowsData := taskReadResult{}

		err := s.DB.QueryRow(`SELECT UID, TITLE, DESCRIPTION, CATEGORY, STATUS, DOMAIN_CODE,
			ASSET_ID, DOMAIN_DATA_AREA_ID, DOMAIN_DATA_MATERIAL_ID
			FROM TASK_READ WHERE UID = ?`, uid).Scan(
			&rowsData.UID,
			&rowsData.Title,
			&rowsData.Description,
			&rowsData.Category,
			&rowsData.Status,
			&rowsData.Domain,
			&rowsData.AssetID,
			&rowsData.AreaID,
			&rowsData.MaterialID,
		)

		if err != nil && err != sql.ErrNoRows {
			result <- query.QueryResult{Error: err}
		}

		if err == sql.ErrNoRows {
			result <- query.QueryResult{Result: taskQueryResult}
		}

		taskUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		assetUID, err := uuid.FromString(rowsData.AssetID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		areaUID := uuid.UUID{}
		if rowsData.AreaID.Valid {
			areaUID, err = uuid.FromString(rowsData.AreaID.String)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}
		}

		materialUID := uuid.UUID{}
		if rowsData.MaterialID.Valid {
			materialUID, err = uuid.FromString(rowsData.MaterialID.String)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}
		}

		taskQueryResult.UID = taskUID
		taskQueryResult.Title = rowsData.Title
		taskQueryResult.Description = rowsData.Description
		taskQueryResult.Status = rowsData.Status
		taskQueryResult.Category = rowsData.Category
		taskQueryResult.Domain = rowsData.Domain
		taskQueryResult.AssetUID = assetUID
		taskQueryResult.AreaUID = areaUID
		taskQueryResult.MaterialUID = materialUID

		result <- query.QueryResult{Result: taskQueryResult}
		close(result)
	}()

	return result
}
