package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/storage"
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

	switch wrapper.EventName {
	case "ReservoirCreated":
		e := domain.ReservoirCreated{}

		for key, v := range mapped {
			if key == "UID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.UID = uid
			}
			if key == "Name" {
				val := v.(string)
				e.Name = val
			}
			if key == "WaterSource" {
				ws, err := makeWaterSource(v)
				if err != nil {
					return nil, err
				}

				e.WaterSource = ws
			}
			if key == "FarmUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.FarmUID = uid
			}
			if key == "CreatedDate" {
				d, err := makeTime(v)
				if err != nil {
					return nil, err
				}

				e.CreatedDate = d
			}
		}

		return e, nil

	case "ReservoirWaterSourceChanged":
		e := domain.ReservoirWaterSourceChanged{}

		for key, v := range mapped {
			if key == "ReservoirUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.ReservoirUID = uid
			}

			if key == "WaterSource" {
				ws, err := makeWaterSource(v)
				if err != nil {
					return nil, err
				}

				e.WaterSource = ws
			}
		}

		return e, nil

	case "ReservoirNameChanged":
		e := domain.ReservoirNameChanged{}

		for key, v := range mapped {
			if key == "ReservoirUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.ReservoirUID = uid
			}
			if key == "Name" {
				name, ok := v.(string)
				if !ok {
					return nil, errors.New("Internal server error. Error type assertion")
				}

				e.Name = name
			}
		}

		return e, nil

	case "ReservoirNoteAdded":
		e := domain.ReservoirNoteAdded{}

		for key, v := range mapped {
			if key == "ReservoirUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.ReservoirUID = uid
			}
			if key == "UID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.UID = uid
			}
			if key == "Content" {
				val := v.(string)
				e.Content = val
			}
			if key == "CreatedDate" {
				d, err := makeTime(v)
				if err != nil {
					return nil, err
				}

				e.CreatedDate = d
			}
		}

		return e, nil

	case "ReservoirNoteRemoved":
		e := domain.ReservoirNoteRemoved{}

		for key, v := range mapped {
			if key == "ReservoirUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.ReservoirUID = uid
			}
			if key == "UID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.UID = uid
			}
		}

		return e, nil
	}

	return nil, nil
}
