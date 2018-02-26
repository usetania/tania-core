package sqlite

import (
	"database/sql"
	"encoding/json"

	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/storage"
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
			ID      int
			FarmUID string
			Version int
			Events  json.RawMessage
		}{}

		for rows.Next() {
			rows.Scan(&rowsData.ID, &rowsData.FarmUID, &rowsData.Version, &rowsData.Events)

			event := storage.FarmEvent{}
			err := json.Unmarshal([]byte(rowsData.Events), &event)
			farmUID, err := uuid.FromString(rowsData.FarmUID)

			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			events = append(events, storage.FarmEvent{
				FarmUID: farmUID,
				Version: rowsData.Version,
				Event:   event,
			})
		}

		result <- query.QueryResult{Result: events}
		close(result)
	}()

	return result
}
