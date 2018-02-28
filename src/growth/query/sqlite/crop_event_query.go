package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Tanibox/tania-server/src/growth/domain"
	"github.com/Tanibox/tania-server/src/growth/query"
	"github.com/Tanibox/tania-server/src/growth/storage"
	uuid "github.com/satori/go.uuid"
)

type CropEventQuerySqlite struct {
	DB *sql.DB
}

func NewCropEventQuerySqlite(db *sql.DB) query.CropEventQuery {
	return &CropEventQuerySqlite{DB: db}
}

func (f *CropEventQuerySqlite) FindAllByCropID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		events := []storage.CropEvent{}

		rows, err := f.DB.Query("SELECT * FROM CROP_EVENT WHERE CROP_UID = ? ORDER BY VERSION ASC", uid)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		rowsData := struct {
			ID          int
			CropUID     string
			Version     int
			CreatedDate string
			Event       []byte
		}{}

		for rows.Next() {
			rows.Scan(&rowsData.ID, &rowsData.CropUID, &rowsData.Version, &rowsData.CreatedDate, &rowsData.Event)

			wrapper := query.EventWrapper{}
			json.Unmarshal(rowsData.Event, &wrapper)

			event, err := assertCropEvent(wrapper)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			cropUID, err := uuid.FromString(rowsData.CropUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			events = append(events, storage.CropEvent{
				CropUID:     cropUID,
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

func assertCropEvent(wrapper query.EventWrapper) (interface{}, error) {
	mapped := wrapper.EventData.(map[string]interface{})

	switch wrapper.EventName {
	case "CropBatchCreated":
		e := domain.CropBatchCreated{}

		if v, ok := mapped["UID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.UID = uid
		}
		if v, ok := mapped["BatchID"]; ok {
			val := v.(string)
			e.BatchID = val
		}
		if v, ok := mapped["Status"]; ok {
			mapped2 := v.(map[string]interface{})

			if v2, ok2 := mapped2["code"]; ok2 {
				st := v2.(string)
				e.Status = domain.GetCropStatus(st)
			}
		}
		if v, ok := mapped["Type"]; ok {
			mapped2 := v.(map[string]interface{})

			if v2, ok2 := mapped2["Code"]; ok2 {
				c := v2.(string)
				e.Type = domain.GetCropType(c)
			}
		}
		if v, ok := mapped["Container"]; ok {
			mapped2 := v.(map[string]interface{})

			if v2, ok2 := mapped2["Quantity"]; ok2 {
				qty := v2.(float64)
				e.Container.Quantity = int(qty)
			}
			if v2, ok2 := mapped2["Type"]; ok2 {
				mapped3 := v2.(map[string]interface{})

				if v3, ok3 := mapped3["Cell"]; ok3 {
					cell := v3.(float64)
					cellInt := int(cell)

					if cellInt == 0 {
						e.Container.Type = domain.Pot{}
					} else {
						e.Container.Type = domain.Tray{Cell: cellInt}
					}
				}
			}

			fmt.Println("CODE", e.Container.Type.Code())
		}
		if v, ok := mapped["InventoryUID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.InventoryUID = uid
		}
		if v, ok := mapped["FarmUID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.FarmUID = uid
		}
		if v, ok := mapped["CreatedDate"]; ok {
			d, err := makeTime(v)
			if err != nil {
				return nil, err
			}

			e.CreatedDate = d
		}
		if v, ok := mapped["InitialAreaUID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.InitialAreaUID = uid
		}
		if v, ok := mapped["Quantity"]; ok {
			val := v.(float64)
			e.Quantity = int(val)
		}

		return e, nil
	}

	return nil, nil
}
