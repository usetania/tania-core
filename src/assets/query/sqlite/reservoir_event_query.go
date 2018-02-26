package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

			fmt.Println("EVENT NAME", wrapper.EventName)
			fmt.Println("EVENT DATA", wrapper.EventData)
			event, err := assertEvent(wrapper)

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

func assertEvent(wrapper query.EventWrapper) (interface{}, error) {
	mapped := wrapper.EventData.(map[string]interface{})

	switch wrapper.EventName {
	case "ReservoirCreated":
		e := domain.ReservoirCreated{}

		for key, v := range mapped {
			if key == "UID" {
				val := v.(string)
				uid, _ := uuid.FromString(val)
				e.UID = uid
			}
			if key == "Name" {
				val := v.(string)
				e.Name = val
			}
			if key == "WaterSource" {
				ws, _ := v.(map[string]interface{})

				convertedMap := map[string]float64{}
				for i2, v2 := range ws {
					convertedMap[i2] = v2.(float64)
				}

				if convertedMap["Capacity"] == 0 {
					e.WaterSource = domain.Tap{}
				} else {
					e.WaterSource = domain.Bucket{Capacity: float32(convertedMap["Capacity"])}
				}
			}
			if key == "FarmUID" {
				val := v.(string)
				uid, _ := uuid.FromString(val)
				e.FarmUID = uid
			}
			if key == "CreatedDate" {
				val := v.(string)

				createdDate, err := time.Parse(time.RFC3339, val)
				if err != nil {
					return nil, err
				}

				e.CreatedDate = createdDate
			}
		}

		return e, nil

	case "ReservoirNoteAdded":
		e := domain.ReservoirNoteAdded{}

		for key, v := range mapped {
			if key == "ReservoirUID" {
				val := v.(string)
				uid, err := uuid.FromString(val)
				if err != nil {
					return nil, err
				}

				e.ReservoirUID = uid
			}
			if key == "UID" {
				val := v.(string)
				uid, err := uuid.FromString(val)
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
				val := v.(string)

				createdDate, err := time.Parse(time.RFC3339, val)
				if err != nil {
					return nil, err
				}

				e.CreatedDate = createdDate
			}
		}

		return e, nil
	}

	return nil, nil
}
