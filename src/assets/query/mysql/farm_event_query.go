package mysql

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Tanibox/tania-core/src/assets/decoder"
	"github.com/Tanibox/tania-core/src/assets/query"
	"github.com/Tanibox/tania-core/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type FarmEventQueryMysql struct {
	DB *sql.DB
}

func NewFarmEventQueryMysql(db *sql.DB) query.FarmEventQuery {
	return &FarmEventQueryMysql{DB: db}
}

func (f *FarmEventQueryMysql) FindAllByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		events := []storage.FarmEvent{}

		rows, err := f.DB.Query("SELECT * FROM FARM_EVENT WHERE FARM_UID = ? ORDER BY VERSION ASC", uid.Bytes())
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		rowsData := struct {
			ID          int
			FarmUID     []byte
			Version     int
			CreatedDate time.Time
			Event       []byte
		}{}

		for rows.Next() {
			rows.Scan(&rowsData.ID, &rowsData.FarmUID, &rowsData.Version, &rowsData.CreatedDate, &rowsData.Event)

			wrapper := decoder.FarmEventWrapper{}
			err := json.Unmarshal(rowsData.Event, &wrapper)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			farmUID, err := uuid.FromBytes(rowsData.FarmUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			events = append(events, storage.FarmEvent{
				FarmUID:     farmUID,
				Version:     rowsData.Version,
				CreatedDate: rowsData.CreatedDate,
				Event:       wrapper.EventData,
			})
		}

		result <- query.QueryResult{Result: events}
		close(result)
	}()

	return result
}
