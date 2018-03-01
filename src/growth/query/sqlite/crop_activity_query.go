package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/Tanibox/tania-server/src/growth/query"
	"github.com/Tanibox/tania-server/src/growth/storage"
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

		rows, err := s.DB.Query(`SELECT * FROM CROP_ACTIVITY WHERE CROP_UID = ?`, uid)
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

			cropUID, err := uuid.FromString(rowsData.CropUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			wrapper := query.ActivityTypeWrapper{}
			json.Unmarshal(rowsData.ActivityType, &wrapper)

			activityType, err := assertActivityType(wrapper)
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

			cropUID, err := uuid.FromString(rowsData.CropUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			wrapper := query.ActivityTypeWrapper{}
			json.Unmarshal(rowsData.ActivityType, &wrapper)

			rowsActivityType, err := assertActivityType(wrapper)
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
				ActivityType:  rowsActivityType,
				CreatedDate:   createdDate,
				Description:   rowsData.Description,
			}
		}

		result <- query.QueryResult{Result: cropActivity}
		close(result)
	}()

	return result
}

func assertActivityType(wrapper query.ActivityTypeWrapper) (storage.ActivityType, error) {
	mapped := wrapper.ActivityData.(map[string]interface{})

	switch wrapper.ActivityName {
	case storage.SeedActivityCode:
		a := storage.SeedActivity{}

		if v, ok := mapped["area_id"]; ok {
			val, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			a.AreaUID = val
		}
		if v, ok := mapped["area_name"]; ok {
			val := v.(string)
			a.AreaName = val
		}
		if v, ok := mapped["quantity"]; ok {
			val := v.(float64)
			a.Quantity = int(val)
		}
		if v, ok := mapped["seeding_date"]; ok {
			val, err := makeTime(v)
			if err != nil {
				return nil, err
			}

			a.SeedingDate = val
		}

		return a, nil

	case storage.MoveActivityCode:
		a := storage.MoveActivity{}

		if v, ok := mapped["source_area_id"]; ok {
			val, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			a.SrcAreaUID = val
		}
		if v, ok := mapped["source_area_name"]; ok {
			val := v.(string)
			a.SrcAreaName = val
		}
		if v, ok := mapped["destination_area_id"]; ok {
			val, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			a.DstAreaUID = val
		}
		if v, ok := mapped["destination_area_name"]; ok {
			val := v.(string)
			a.DstAreaName = val
		}
		if v, ok := mapped["quantity"]; ok {
			val := v.(float64)
			a.Quantity = int(val)
		}
		if v, ok := mapped["moved_date"]; ok {
			val, err := makeTime(v)
			if err != nil {
				return nil, err
			}

			a.MovedDate = val
		}

		return a, nil

	case storage.HarvestActivityCode:
		a := storage.HarvestActivity{}

		if v, ok := mapped["type"]; ok {
			val := v.(string)
			a.Type = val
		}
		if v, ok := mapped["source_area_id"]; ok {
			val, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			a.SrcAreaUID = val
		}
		if v, ok := mapped["source_area_name"]; ok {
			val := v.(string)
			a.SrcAreaName = val
		}
		if v, ok := mapped["quantity"]; ok {
			val := v.(float64)
			a.Quantity = int(val)
		}
		if v, ok := mapped["produced_gram_quantity"]; ok {
			val := v.(float64)
			a.ProducedGramQuantity = float32(val)
		}
		if v, ok := mapped["harvest_date"]; ok {
			val, err := makeTime(v)
			if err != nil {
				return nil, err
			}

			a.HarvestDate = val
		}

		return a, nil

	case storage.DumpActivityCode:
		a := storage.DumpActivity{}

		if v, ok := mapped["source_area_id"]; ok {
			val, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			a.SrcAreaUID = val
		}
		if v, ok := mapped["source_area_name"]; ok {
			val := v.(string)
			a.SrcAreaName = val
		}
		if v, ok := mapped["quantity"]; ok {
			val := v.(float64)
			a.Quantity = int(val)
		}
		if v, ok := mapped["dump_date"]; ok {
			val, err := makeTime(v)
			if err != nil {
				return nil, err
			}

			a.DumpDate = val
		}

		return a, nil

	case storage.WaterActivityCode:
		a := storage.WaterActivity{}

		if v, ok := mapped["area_id"]; ok {
			val, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			a.AreaUID = val
		}
		if v, ok := mapped["area_name"]; ok {
			val := v.(string)
			a.AreaName = val
		}
		if v, ok := mapped["watering_date"]; ok {
			val, err := makeTime(v)
			if err != nil {
				return nil, err
			}

			a.WateringDate = val
		}

		return a, nil

	case storage.PhotoActivityCode:
		a := storage.PhotoActivity{}

		if v, ok := mapped["uid"]; ok {
			val, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			a.UID = val
		}
		if v, ok := mapped["filename"]; ok {
			val := v.(string)
			a.Filename = val
		}
		if v, ok := mapped["mime_type"]; ok {
			val := v.(string)
			a.MimeType = val
		}
		if v, ok := mapped["size"]; ok {
			val := v.(float64)
			a.Size = int(val)
		}
		if v, ok := mapped["width"]; ok {
			val := v.(float64)
			a.Width = int(val)
		}
		if v, ok := mapped["height"]; ok {
			val := v.(float64)
			a.Height = int(val)
		}
		if v, ok := mapped["description"]; ok {
			val := v.(string)
			a.Description = val
		}

		return a, nil
	}

	return nil, nil
}
