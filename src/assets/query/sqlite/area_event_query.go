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

type AreaEventQuerySqlite struct {
	DB *sql.DB
}

func NewAreaEventQuerySqlite(db *sql.DB) query.AreaEventQuery {
	return &AreaEventQuerySqlite{DB: db}
}

func (f *AreaEventQuerySqlite) FindAllByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		events := []storage.AreaEvent{}

		rows, err := f.DB.Query("SELECT * FROM AREA_EVENT WHERE AREA_UID = ? ORDER BY VERSION ASC", uid)
		if err != nil {
			result <- query.QueryResult{Error: err}
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

			wrapper := query.EventWrapper{}
			json.Unmarshal(rowsData.Event, &wrapper)

			event, err := assertAreaEvent(wrapper)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			areaUID, err := uuid.FromString(rowsData.AreaUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			events = append(events, storage.AreaEvent{
				AreaUID:     areaUID,
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

func assertAreaEvent(wrapper query.EventWrapper) (interface{}, error) {
	mapped := wrapper.EventData.(map[string]interface{})

	f := mapstructure.ComposeDecodeHookFunc(
		decoder.UIDHook(),
		decoder.TimeHook(time.RFC3339),
	)

	switch wrapper.EventName {
	case "AreaCreated":
		e := domain.AreaCreated{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "AreaNameChanged":
		e := domain.AreaNameChanged{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "AreaSizeChanged":
		e := domain.AreaSizeChanged{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "AreaTypeChanged":
		e := domain.AreaTypeChanged{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "AreaLocationChanged":
		e := domain.AreaLocationChanged{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "AreaReservoirChanged":
		e := domain.AreaReservoirChanged{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "AreaPhotoAdded":
		e := domain.AreaPhotoAdded{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "AreaNoteAdded":
		e := domain.AreaNoteAdded{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "AreaNoteRemoved":
		e := domain.AreaNoteRemoved{}

		decoder.Decode(f, &mapped, &e)

		return e, nil
	}

	return nil, nil
}
