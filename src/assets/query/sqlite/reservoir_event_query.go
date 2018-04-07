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

type ReservoirEventQuerySqlite struct {
	DB *sql.DB
}

func NewReservoirEventQuerySqlite(db *sql.DB) query.ReservoirEventQuery {
	return &ReservoirEventQuerySqlite{DB: db}
}

func (f *ReservoirEventQuerySqlite) FindAllByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		events := []storage.ReservoirEvent{}

		rows, err := f.DB.Query("SELECT * FROM RESERVOIR_EVENT WHERE RESERVOIR_UID = ? ORDER BY VERSION ASC", uid)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		rowsData := struct {
			ID           int
			ReservoirUID string
			Version      int
			CreatedDate  string
			Event        []byte
		}{}

		for rows.Next() {
			rows.Scan(&rowsData.ID, &rowsData.ReservoirUID, &rowsData.Version, &rowsData.CreatedDate, &rowsData.Event)

			wrapper := decoder.ReservoirEventWrapper{}
			err := json.Unmarshal(rowsData.Event, &wrapper)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			reservoirUID, err := uuid.FromString(rowsData.ReservoirUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			events = append(events, storage.ReservoirEvent{
				ReservoirUID: reservoirUID,
				Version:      rowsData.Version,
				CreatedDate:  createdDate,
				Event:        wrapper.EventData,
			})
		}

		result <- query.QueryResult{Result: events}
		close(result)
	}()

	return result
}
