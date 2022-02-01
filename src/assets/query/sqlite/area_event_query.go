package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Tanibox/tania-core/src/assets/decoder"
	"github.com/Tanibox/tania-core/src/assets/query"
	"github.com/Tanibox/tania-core/src/assets/storage"
	"github.com/gofrs/uuid"
)

type AreaEventQuerySqlite struct {
	DB *sql.DB
}

func NewAreaEventQuerySqlite(db *sql.DB) query.AreaEvent {
	return &AreaEventQuerySqlite{DB: db}
}

func (f *AreaEventQuerySqlite) FindAllByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		events := []storage.AreaEvent{}

		rows, err := f.DB.Query("SELECT * FROM AREA_EVENT WHERE AREA_UID = ? ORDER BY VERSION ASC", uid)
		if err != nil {
			result <- query.Result{Error: err}
		}

		rowsData := struct {
			ID          int
			AreaUID     string
			Version     int
			CreatedDate string
			Event       []byte
		}{}

		for rows.Next() {
			rows.Scan(&rowsData.ID, &rowsData.AreaUID, &rowsData.Version, &rowsData.CreatedDate, &rowsData.Event)

			wrapper := decoder.AreaEventWrapper{}

			err := json.Unmarshal(rowsData.Event, &wrapper)
			if err != nil {
				result <- query.Result{Error: err}
			}

			areaUID, err := uuid.FromString(rowsData.AreaUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.Result{Error: err}
			}

			events = append(events, storage.AreaEvent{
				AreaUID:     areaUID,
				Version:     rowsData.Version,
				CreatedDate: createdDate,
				Event:       wrapper.EventData,
			})
		}

		result <- query.Result{Result: events}
		close(result)
	}()

	return result
}
