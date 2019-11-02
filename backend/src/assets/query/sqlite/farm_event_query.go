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

type FarmEventQuerySqlite struct {
	DB *sql.DB
}

func NewFarmEventQuerySqlite(db *sql.DB) query.FarmEventQuery {
	return &FarmEventQuerySqlite{DB: db}
}

func (f *FarmEventQuerySqlite) FindAllByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		events := []storage.FarmEvent{}

		rows, err := f.DB.Query("SELECT * FROM FARM_EVENT WHERE FARM_UID = ? ORDER BY VERSION ASC", uid)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		rowsData := struct {
			ID          int
			FarmUID     string
			Version     int
			CreatedDate string
			Event       []byte
		}{}

		for rows.Next() {
			rows.Scan(&rowsData.ID, &rowsData.FarmUID, &rowsData.Version, &rowsData.CreatedDate, &rowsData.Event)

			wrapper := decoder.FarmEventWrapper{}
			err := json.Unmarshal(rowsData.Event, &wrapper)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			farmUID, err := uuid.FromString(rowsData.FarmUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			events = append(events, storage.FarmEvent{
				FarmUID:     farmUID,
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
