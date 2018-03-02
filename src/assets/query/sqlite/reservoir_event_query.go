package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/storage"
	"github.com/Tanibox/tania-server/src/assets/util/decoder"
	"github.com/mitchellh/mapstructure"
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

			wrapper := query.EventWrapper{}
			json.Unmarshal(rowsData.Event, &wrapper)

			event, err := assertReservoirEvent(wrapper)
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
				Event:        event,
			})
		}

		result <- query.QueryResult{Result: events}
		close(result)
	}()

	return result
}

func assertReservoirEvent(wrapper query.EventWrapper) (interface{}, error) {
	mapped := wrapper.EventData.(map[string]interface{})

	f := mapstructure.ComposeDecodeHookFunc(
		decoder.UIDHook(),
		decoder.TimeHook(time.RFC3339),
		decoder.WaterSourceHook(),
	)

	switch wrapper.EventName {
	case "ReservoirCreated":
		e := domain.ReservoirCreated{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "ReservoirWaterSourceChanged":
		e := domain.ReservoirWaterSourceChanged{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "ReservoirNameChanged":
		e := domain.ReservoirNameChanged{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "ReservoirNoteAdded":
		e := domain.ReservoirNoteAdded{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "ReservoirNoteRemoved":
		e := domain.ReservoirNoteRemoved{}

		decoder.Decode(f, &mapped, &e)

		return e, nil
	}

	return nil, nil
}
