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

type FarmEventQueryMysql struct {
	DB *sql.DB
}

func NewFarmEventQueryMysql(db *sql.DB) query.FarmEvent {
	return &FarmEventQueryMysql{DB: db}
}

func (f *FarmEventQueryMysql) FindAllByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		events := []storage.FarmEvent{}

		rows, err := f.DB.Query("SELECT * FROM FARM_EVENT WHERE FARM_UID = ? ORDER BY VERSION ASC", uid.Bytes())
		if err != nil {
			result <- query.Result{Error: err}
		}

		rowsData := struct {
			ID          int
			FarmUID     []byte
			Version     int
			CreatedDate time.Time
			Event       []byte
		}{}

		for rows.Next() {
			err := rows.Scan(&rowsData.ID, &rowsData.FarmUID, &rowsData.Version, &rowsData.CreatedDate, &rowsData.Event)
			if err != nil {
				result <- query.Result{Error: err}
			}

			wrapper := decoder.FarmEventWrapper{}

			err = json.Unmarshal(rowsData.Event, &wrapper)
			if err != nil {
				result <- query.Result{Error: err}
			}

			farmUID, err := uuid.FromBytes(rowsData.FarmUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			events = append(events, storage.FarmEvent{
				FarmUID:     farmUID,
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
