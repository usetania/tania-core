package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/growth/decoder"
	"github.com/usetania/tania-core/src/growth/query"
	"github.com/usetania/tania-core/src/growth/storage"
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

func (s CropActivityQuerySqlite) FindAllByCropID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		cropActivities := []storage.CropActivity{}
		rowsData := cropActivityResult{}

		rows, err := s.DB.Query(`SELECT * FROM CROP_ACTIVITY WHERE CROP_UID = ? ORDER BY CREATED_DATE DESC`, uid)
		if err != nil {
			result <- query.Result{Error: err}
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
			if err != nil {
				result <- query.Result{Error: err}
			}

			wrapper := decoder.CropActivityTypeWrapper{}
			json.Unmarshal(rowsData.ActivityType, &wrapper)

			activityType, ok := wrapper.Data.(storage.ActivityType)
			if !ok {
				result <- query.Result{Error: errors.New("error type assertion")}
			}

			cropUID, err := uuid.FromString(rowsData.CropUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.Result{Error: err}
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

		result <- query.Result{Result: cropActivities}
		close(result)
	}()

	return result
}

func (s CropActivityQuerySqlite) FindByCropIDAndActivityType(
	uid uuid.UUID,
	activityType interface{},
) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		cropActivity := storage.CropActivity{}
		rowsData := cropActivityResult{}

		at, ok := activityType.(storage.ActivityType)
		if !ok {
			result <- query.Result{Error: errors.New("wrong activity type")}
		}

		rows, err := s.DB.Query(`SELECT * FROM CROP_ACTIVITY
			WHERE CROP_UID = ? AND ACTIVITY_TYPE_CODE = ?`, uid, at.Code())
		if err != nil {
			result <- query.Result{Error: err}
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
			if err != nil {
				result <- query.Result{Error: err}
			}

			wrapper := decoder.CropActivityTypeWrapper{}
			json.Unmarshal(rowsData.ActivityType, &wrapper)

			activityType, ok := wrapper.Data.(storage.ActivityType)
			if !ok {
				result <- query.Result{Error: errors.New("error type assertion")}
			}

			cropUID, err := uuid.FromString(rowsData.CropUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.Result{Error: err}
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

		result <- query.Result{Result: cropActivity}
		close(result)
	}()

	return result
}
