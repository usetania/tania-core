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

	switch wrapper.EventName {
	case "AreaCreated":
		e := domain.AreaCreated{}

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
				at, err := makeAreaType(v)
				if err != nil {
					return nil, err
				}

				e.Type = at
			}
			if key == "Location" {
				al, err := makeAreaLocation(v)
				if err != nil {
					return nil, err
				}

				e.Location = al
			}
			if key == "Size" {
				areaSize, err := makeAreaSize(v)
				if err != nil {
					return nil, err
				}

				e.Size = areaSize
			}
			if key == "FarmUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.FarmUID = uid
			}
			if key == "ReservoirUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.ReservoirUID = uid
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

	case "AreaNameChanged":
		e := domain.AreaNameChanged{}

		for key, v := range mapped {
			if key == "AreaUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.AreaUID = uid
			}
			if key == "Name" {
				val := v.(string)
				e.Name = val
			}
		}

		return e, nil

	case "AreaSizeChanged":
		e := domain.AreaSizeChanged{}

		for key, v := range mapped {
			if key == "AreaUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.AreaUID = uid
			}
			if key == "Size" {
				areaSize, err := makeAreaSize(v)
				if err != nil {
					return nil, err
				}

				e.Size = areaSize
			}
		}

		return e, nil

	case "AreaTypeChanged":
		e := domain.AreaTypeChanged{}

		for key, v := range mapped {
			if key == "AreaUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.AreaUID = uid
			}
			if key == "Type" {
				at, err := makeAreaType(v)
				if err != nil {
					return nil, err
				}

				e.Type = at
			}
		}

		return e, nil

	case "AreaLocationChanged":
		e := domain.AreaLocationChanged{}

		for key, v := range mapped {
			if key == "AreaUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.AreaUID = uid
			}
			if key == "Location" {
				al, err := makeAreaLocation(v)
				if err != nil {
					return nil, err
				}

				e.Location = al
			}
		}

		return e, nil

	case "AreaReservoirChanged":
		e := domain.AreaReservoirChanged{}

		for key, v := range mapped {
			if key == "AreaUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.AreaUID = uid
			}
			if key == "ReservoirUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.ReservoirUID = uid
			}
		}

		return e, nil

	case "AreaPhotoAdded":
		e := domain.AreaPhotoAdded{}

		for key, v := range mapped {
			if key == "AreaUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.AreaUID = uid
			}
			if key == "Filename" {
				val := v.(string)
				e.Filename = val
			}
			if key == "MimeType" {
				val := v.(string)
				e.MimeType = val
			}
			if key == "Size" {
				val := v.(float64)
				e.Size = int(val)
			}
			if key == "Width" {
				val := v.(float64)
				e.Width = int(val)
			}
			if key == "Height" {
				val := v.(float64)
				e.Height = int(val)
			}
		}

		return e, nil

	case "AreaNoteAdded":
		e := domain.AreaNoteAdded{}

		for key, v := range mapped {
			if key == "AreaUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.AreaUID = uid
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

	case "AreaNoteRemoved":
		e := domain.AreaNoteRemoved{}

		for key, v := range mapped {
			if key == "AreaUID" {
				uid, err := makeUUID(v)
				if err != nil {
					return nil, err
				}

				e.AreaUID = uid
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
