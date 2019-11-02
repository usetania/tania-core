package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/Tanibox/tania-core/src/growth/decoder"
	"github.com/Tanibox/tania-core/src/growth/query"
	"github.com/Tanibox/tania-core/src/growth/storage"
	uuid "github.com/satori/go.uuid"
)

type CropActivityQuerySqlite struct {
	DB *sql.DB
}

func NewCropActivityQuerySqlite(db *sql.DB) query.CropActivityQuery {
	return CropActivityQuerySqlite{DB: db}
}

type cropActivityResult struct {
	ID               string
	CropUID          string
	BatchID          string
	ContainerType    string
	ActivityType     []byte
	ActivityTypeCode string
	CreatedDate      string
	Description      string
}

func (s CropActivityQuerySqlite) FindAllByCropID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		cropActivities := []storage.CropActivity{}
		rowsData := cropActivityResult{}

		rows, err := s.DB.Query(`SELECT * FROM CROP_ACTIVITY WHERE CROP_UID = ? ORDER BY CREATED_DATE DESC`, uid)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		for rows.Next() {
			err = rows.Scan(
				&rowsData.ID,
				&rowsData.CropUID,
				&rowsData.BatchID,
				&rowsData.ContainerType,
				&rowsData.ActivityType,
				&rowsData.ActivityTypeCode,
				&rowsData.CreatedDate,
				&rowsData.Description,
			)

			wrapper := decoder.CropActivityTypeWrapper{}
			json.Unmarshal(rowsData.ActivityType, &wrapper)

			activityType, ok := wrapper.Data.(storage.ActivityType)
			if !ok {
				result <- query.QueryResult{Error: errors.New("Error type assertion")}
			}

			cropUID, err := uuid.FromString(rowsData.CropUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			cropActivities = append(cropActivities, storage.CropActivity{
				UID:           cropUID,
				BatchID:       rowsData.BatchID,
				ContainerType: rowsData.ContainerType,
				ActivityType:  activityType,
				CreatedDate:   createdDate,
				Description:   rowsData.Description,
			})
		}

		result <- query.QueryResult{Result: cropActivities}
		close(result)
	}()

	return result
}

func (s CropActivityQuerySqlite) FindByCropIDAndActivityType(uid uuid.UUID, activityType interface{}) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		cropActivity := storage.CropActivity{}
		rowsData := cropActivityResult{}

		at, ok := activityType.(storage.ActivityType)
		if !ok {
			result <- query.QueryResult{Error: errors.New("Wrong activity type")}
		}

		rows, err := s.DB.Query(`SELECT * FROM CROP_ACTIVITY
			WHERE CROP_UID = ? AND ACTIVITY_TYPE_CODE = ?`, uid, at.Code())
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		for rows.Next() {
			err = rows.Scan(
				&rowsData.ID,
				&rowsData.CropUID,
				&rowsData.BatchID,
				&rowsData.ContainerType,
				&rowsData.ActivityType,
				&rowsData.ActivityTypeCode,
				&rowsData.CreatedDate,
				&rowsData.Description,
			)

			wrapper := decoder.CropActivityTypeWrapper{}
			json.Unmarshal(rowsData.ActivityType, &wrapper)

			activityType, ok := wrapper.Data.(storage.ActivityType)
			if !ok {
				result <- query.QueryResult{Error: errors.New("Error type assertion")}
			}

			cropUID, err := uuid.FromString(rowsData.CropUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			cropActivity = storage.CropActivity{
				UID:           cropUID,
				BatchID:       rowsData.BatchID,
				ContainerType: rowsData.ContainerType,
				ActivityType:  activityType,
				CreatedDate:   createdDate,
				Description:   rowsData.Description,
			}
		}

		result <- query.QueryResult{Result: cropActivity}
		close(result)
	}()

	return result
}
