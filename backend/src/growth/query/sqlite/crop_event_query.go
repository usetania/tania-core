package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/growth/decoder"
	"github.com/usetania/tania-core/src/growth/query"
	"github.com/usetania/tania-core/src/growth/storage"
)

type CropEventQuerySqlite struct {
	DB *sql.DB
}

func NewCropEventQuerySqlite(db *sql.DB) query.CropEventQuery {
	return &CropEventQuerySqlite{DB: db}
}

func (f *CropEventQuerySqlite) FindAllByCropID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		events := []storage.CropEvent{}

		rows, err := f.DB.Query("SELECT * FROM CROP_EVENT WHERE CROP_UID = ? ORDER BY VERSION ASC", uid)
		if err != nil {
			result <- query.Result{Error: err}
		}

		rowsData := struct {
			ID          int
			CropUID     string
			Version     int
			CreatedDate string
			Event       []byte
		}{}

		for rows.Next() {
			rows.Scan(&rowsData.ID, &rowsData.CropUID, &rowsData.Version, &rowsData.CreatedDate, &rowsData.Event)

			wrapper := decoder.CropEventWrapper{}

			err = json.Unmarshal(rowsData.Event, &wrapper)
			if err != nil {
				result <- query.Result{Error: err}
			}

			cropUID, err := uuid.FromString(rowsData.CropUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.Result{Error: err}
			}

			events = append(events, storage.CropEvent{
				CropUID:     cropUID,
				Version:     rowsData.Version,
				CreatedDate: createdDate,
				Event:       wrapper.Data,
			})
		}

		result <- query.Result{Result: events}
		close(result)
	}()

	return result
}
