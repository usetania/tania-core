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

type MaterialEventQuerySqlite struct {
	DB *sql.DB
}

func NewMaterialEventQuerySqlite(db *sql.DB) query.MaterialEventQuery {
	return &MaterialEventQuerySqlite{DB: db}
}

func (f *MaterialEventQuerySqlite) FindAllByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		events := []storage.MaterialEvent{}

		rows, err := f.DB.Query("SELECT * FROM MATERIAL_EVENT WHERE MATERIAL_UID = ? ORDER BY VERSION ASC", uid)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		rowsData := struct {
			ID          int
			MaterialUID string
			Version     int
			CreatedDate string
			Event       []byte
		}{}

		for rows.Next() {
			rows.Scan(&rowsData.ID, &rowsData.MaterialUID, &rowsData.Version, &rowsData.CreatedDate, &rowsData.Event)
			wrapper := query.EventWrapper{}
			json.Unmarshal(rowsData.Event, &wrapper)

			event, err := assertMaterialEvent(wrapper)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			materialUID, err := uuid.FromString(rowsData.MaterialUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			events = append(events, storage.MaterialEvent{
				MaterialUID: materialUID,
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

func assertMaterialEvent(wrapper query.EventWrapper) (interface{}, error) {
	mapped := wrapper.EventData.(map[string]interface{})

	switch wrapper.EventName {
	case "MaterialCreated":
		e := domain.MaterialCreated{}

		if v, ok := mapped["UID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.UID = uid
		}

		if v, ok := mapped["Name"]; ok {
			val := v.(string)
			e.Name = val
		}

		if v, ok := mapped["PricePerUnit"]; ok {
			val, err := makeMaterialPricePerUnit(v)
			if err != nil {
				return nil, err
			}

			e.PricePerUnit = val
		}

		if v, ok := mapped["Type"]; ok {
			val, err := makeMaterialType(v)
			if err != nil {
				return nil, err
			}

			e.Type = val
		}

		if v, ok := mapped["Quantity"]; ok {
			val, err := makeMaterialQuantity(v, e.Type.Code())
			if err != nil {
				return nil, err
			}

			e.Quantity = val
		}

		if v, ok := mapped["ExpirationDate"]; ok {
			if v != nil {
				d, err := makeTime(v)
				if err != nil {
					return nil, err
				}

				e.ExpirationDate = &d
			}
		}

		if v, ok := mapped["Notes"]; ok {
			if v != nil {
				v := v.(string)
				e.Notes = &v
			}
		}

		if v, ok := mapped["ProducedBy"]; ok {
			if v != nil {
				v := v.(string)
				e.ProducedBy = &v
			}
		}

		if v, ok := mapped["CreatedDate"]; ok {
			d, err := makeTime(v)
			if err != nil {
				return nil, err
			}

			e.CreatedDate = d
		}

		return e, nil

	case "MaterialNameChanged":
		e := domain.MaterialNameChanged{}

		if v, ok := mapped["MaterialUID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.MaterialUID = uid
		}
		if v, ok := mapped["Name"]; ok {
			val := v.(string)
			e.Name = val
		}

		return e, nil

	case "MaterialPriceChanged":
		e := domain.MaterialPriceChanged{}

		if v, ok := mapped["MaterialUID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.MaterialUID = uid
		}
		if v, ok := mapped["Price"]; ok {
			val, err := makeMaterialPricePerUnit(v)
			if err != nil {
				return nil, err
			}

			e.Price = val
		}

		return e, nil

	case "MaterialQuantityChanged":
		e := domain.MaterialQuantityChanged{}

		if v, ok := mapped["MaterialUID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.MaterialUID = uid
		}
		if v, ok := mapped["Quantity"]; ok {
			typeCode := mapped["MaterialTypeCode"].(string)
			val, err := makeMaterialQuantity(v, typeCode)
			if err != nil {
				return nil, err
			}

			e.Quantity = val
		}

		return e, nil

	case "MaterialTypeChanged":
		e := domain.MaterialTypeChanged{}

		if v, ok := mapped["MaterialUID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.MaterialUID = uid
		}
		if v, ok := mapped["MaterialType"]; ok {
			val, err := makeMaterialType(v)
			if err != nil {
				return nil, err
			}

			e.MaterialType = val
		}

		return e, nil

	case "MaterialExpirationDateChanged":
		e := domain.MaterialExpirationDateChanged{}

		if v, ok := mapped["MaterialUID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.MaterialUID = uid
		}
		if v, ok := mapped["ExpirationDate"]; ok {
			if v != nil {
				d, err := makeTime(v)
				if err != nil {
					return nil, err
				}

				e.ExpirationDate = d
			}
		}

	case "MaterialNotesChanged":
		e := domain.MaterialNotesChanged{}

		if v, ok := mapped["MaterialUID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.MaterialUID = uid
		}
		if v, ok := mapped["Notes"]; ok {
			if v != nil {
				val := v.(string)
				e.Notes = val
			}
		}

		return e, nil

	case "MaterialProducedByChanged":
		e := domain.MaterialProducedByChanged{}

		if v, ok := mapped["MaterialUID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.MaterialUID = uid
		}
		if v, ok := mapped["ProducedBy"]; ok {
			if v != nil {
				val := v.(string)
				e.ProducedBy = val
			}
		}

		return e, nil
	}

	return nil, nil
}
