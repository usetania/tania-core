package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Tanibox/tania-core/src/assets/decoder"
	"github.com/Tanibox/tania-core/src/assets/query"
	"github.com/Tanibox/tania-core/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type MaterialEventQuerySqlite struct {
	DB *sql.DB
}

func NewMaterialEventQuerySqlite(db *sql.DB) query.MaterialEventQuery {
	return &MaterialEventQuerySqlite{DB: db}
}

func (f *MaterialEventQuerySqlite) FindAllByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		events := []storage.MaterialEvent{}

		rows, err := f.DB.Query("SELECT * FROM MATERIAL_EVENT WHERE MATERIAL_UID = ? ORDER BY VERSION ASC", uid)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		rowsData := struct {
			ID          int
			MaterialUID string
			Version     int
			CreatedDate string
			Event       []byte
		}{}

		for rows.Next() {
			rows.Scan(&rowsData.ID, &rowsData.MaterialUID, &rowsData.Version, &rowsData.CreatedDate, &rowsData.Event)

			wrapper := decoder.MaterialEventWrapper{}
			err := json.Unmarshal(rowsData.Event, &wrapper)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			materialUID, err := uuid.FromString(rowsData.MaterialUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			events = append(events, storage.MaterialEvent{
				MaterialUID: materialUID,
				Version:     rowsData.Version,
				CreatedDate: createdDate,
				Event:       wrapper.EventData,
			})
		}

		result <- query.QueryResult{Result: events}
		close(result)
	}()

	return result
}
