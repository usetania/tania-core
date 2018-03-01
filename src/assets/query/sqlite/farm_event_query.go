package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/storage"
	"github.com/mitchellh/mapstructure"
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

			wrapper := query.EventWrapper{}
			json.Unmarshal(rowsData.Event, &wrapper)

			event, err := assertFarmEvent(wrapper)
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
				Event:       event,
			})
		}

		result <- query.QueryResult{Result: events}
		close(result)
	}()

	return result
}

func assertFarmEvent(wrapper query.EventWrapper) (interface{}, error) {
	mapped := wrapper.EventData.(map[string]interface{})

	f := mapstructure.ComposeDecodeHookFunc(
		query.UIDHook(),
		query.TimeRFC3339Hook(time.RFC3339),
	)

	switch wrapper.EventName {
	case "FarmCreated":
		e := domain.FarmCreated{}

		query.Decode(f, &mapped, &e)

		return e, nil

	case "FarmNameChanged":
		e := domain.FarmNameChanged{}

		query.Decode(f, &mapped, &e)

		return e, nil

	case "FarmTypeChanged":
		e := domain.FarmTypeChanged{}

		query.Decode(f, &mapped, &e)

		return e, nil

	case "FarmGeolocationChanged":
		e := domain.FarmGeolocationChanged{}

		query.Decode(f, &mapped, &e)

		return e, nil

	case "FarmRegionChanged":
		e := domain.FarmRegionChanged{}

		query.Decode(f, &mapped, &e)

		return e, nil
	}

	return nil, errors.New("Event not decoded succesfully")
}
