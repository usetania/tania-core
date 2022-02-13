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

type ReservoirEventQueryMysql struct {
	DB *sql.DB
}

func NewReservoirEventQueryMysql(db *sql.DB) query.ReservoirEvent {
	return &ReservoirEventQueryMysql{DB: db}
}

func (f *ReservoirEventQueryMysql) FindAllByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		events := []storage.ReservoirEvent{}

		rows, err := f.DB.Query("SELECT * FROM RESERVOIR_EVENT WHERE RESERVOIR_UID = ? ORDER BY VERSION ASC", uid.Bytes())
		if err != nil {
			result <- query.Result{Error: err}
		}

		rowsData := struct {
			ID           int
			ReservoirUID []byte
			Version      int
			CreatedDate  time.Time
			Event        []byte
		}{}

		for rows.Next() {
			err := rows.Scan(&rowsData.ID, &rowsData.ReservoirUID, &rowsData.Version, &rowsData.CreatedDate, &rowsData.Event)
			if err != nil {
				result <- query.Result{Error: err}
			}

			wrapper := decoder.ReservoirEventWrapper{}

			err = json.Unmarshal(rowsData.Event, &wrapper)
			if err != nil {
				result <- query.Result{Error: err}
			}

			reservoirUID, err := uuid.FromBytes(rowsData.ReservoirUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			events = append(events, storage.ReservoirEvent{
				ReservoirUID: reservoirUID,
				Version:      rowsData.Version,
				CreatedDate:  rowsData.CreatedDate,
				Event:        wrapper.EventData,
			})
		}

		result <- query.Result{Result: events}
		close(result)
	}()

	return result
}
