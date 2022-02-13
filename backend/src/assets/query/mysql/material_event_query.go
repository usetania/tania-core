package mysql

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/decoder"
	"github.com/usetania/tania-core/src/assets/query"
	"github.com/usetania/tania-core/src/assets/storage"
)

type MaterialEventQueryMysql struct {
	DB *sql.DB
}

func NewMaterialEventQueryMysql(db *sql.DB) query.MaterialEvent {
	return &MaterialEventQueryMysql{DB: db}
}

func (f *MaterialEventQueryMysql) FindAllByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		events := []storage.MaterialEvent{}

		rows, err := f.DB.Query("SELECT * FROM MATERIAL_EVENT WHERE MATERIAL_UID = ? ORDER BY VERSION ASC", uid.Bytes())
		if err != nil {
			result <- query.Result{Error: err}
		}

		rowsData := struct {
			ID          int
			MaterialUID []byte
			Version     int
			CreatedDate time.Time
			Event       []byte
		}{}

		for rows.Next() {
			err := rows.Scan(&rowsData.ID, &rowsData.MaterialUID, &rowsData.Version, &rowsData.CreatedDate, &rowsData.Event)
			if err != nil {
				result <- query.Result{Error: err}
			}

			wrapper := decoder.MaterialEventWrapper{}

			err = json.Unmarshal(rowsData.Event, &wrapper)
			if err != nil {
				result <- query.Result{Error: err}
			}

			materialUID, err := uuid.FromBytes(rowsData.MaterialUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			events = append(events, storage.MaterialEvent{
				MaterialUID: materialUID,
				Version:     rowsData.Version,
				CreatedDate: rowsData.CreatedDate,
				Event:       wrapper.EventData,
			})
		}

		result <- query.Result{Result: events}
		close(result)
	}()

	return result
}
