package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Tanibox/tania-server/src/assets/domain"
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

	switch wrapper.EventName {
	case "FarmCreated":
		e := domain.FarmCreated{}

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
			if key == "Type" {
				val := v.(string)
				e.Type = val
			}
			if key == "Latitude" {
				val := v.(string)
				e.Latitude = val
			}
			if key == "Longitude" {
				val := v.(string)
				e.Longitude = val
			}
			if key == "CountryCode" {
				val := v.(string)
				e.CountryCode = val
			}
			if key == "CityCode" {
				val := v.(string)
				e.CityCode = val
			}
			if key == "IsActive" {
				val := v.(bool)
				e.IsActive = val
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

	case "FarmNameChanged":
		e := domain.FarmNameChanged{}

		for key, v := range mapped {
			if key == "FarmUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.FarmUID = uid
			}
			if key == "Name" {
				val := v.(string)
				e.Name = val
			}
		}

		return e, nil

	case "FarmTypeChanged":
		e := domain.FarmTypeChanged{}

		for key, v := range mapped {
			if key == "FarmUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.FarmUID = uid
			}
			if key == "Type" {
				val := v.(string)
				e.Type = val
			}
		}

		return e, nil

	case "FarmGeolocationChanged":
		e := domain.FarmGeolocationChanged{}

		for key, v := range mapped {
			if key == "FarmUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.FarmUID = uid
			}
			if key == "Latitude" {
				val := v.(string)
				e.Latitude = val
			}
			if key == "Longitude" {
				val := v.(string)
				e.Longitude = val
			}
		}

		return e, nil

	case "FarmRegionChanged":
		e := domain.FarmRegionChanged{}

		for key, v := range mapped {
			if key == "FarmUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.FarmUID = uid
			}
			if key == "CountryCode" {
				val := v.(string)
				e.CountryCode = val
			}
			if key == "CityCode" {
				val := v.(string)
				e.CityCode = val
			}
		}

		return e, nil
	}

	return nil, nil
}
