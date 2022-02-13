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

type CropEventQueryMysql struct {
	DB *sql.DB
}

func NewCropEventQueryMysql(db *sql.DB) query.CropEventQuery {
	return &CropEventQueryMysql{DB: db}
}

func (f *CropEventQueryMysql) FindAllByCropID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		events := []storage.CropEvent{}

		rows, err := f.DB.Query("SELECT * FROM CROP_EVENT WHERE CROP_UID = ? ORDER BY VERSION ASC", uid.Bytes())
		if err != nil {
			result <- query.Result{Error: err}
		}

		rowsData := struct {
			ID          int
			CropUID     []byte
			Version     int
			CreatedDate time.Time
			Event       []byte
		}{}

		for rows.Next() {
			rows.Scan(&rowsData.ID, &rowsData.CropUID, &rowsData.Version, &rowsData.CreatedDate, &rowsData.Event)

			wrapper := decoder.CropEventWrapper{}

			err = json.Unmarshal(rowsData.Event, &wrapper)
			if err != nil {
				result <- query.Result{Error: err}
			}

			cropUID, err := uuid.FromBytes(rowsData.CropUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			events = append(events, storage.CropEvent{
				CropUID:     cropUID,
				Version:     rowsData.Version,
				CreatedDate: rowsData.CreatedDate,
				Event:       wrapper.Data,
			})
		}

		result <- query.Result{Result: events}
		close(result)
	}()

	return result
}
