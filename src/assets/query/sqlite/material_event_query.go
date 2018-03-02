package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/storage"
	"github.com/mitchellh/mapstructure"
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

	f := mapstructure.ComposeDecodeHookFunc(
		query.UIDHook(),
		query.TimeHook(time.RFC3339),
	)

	switch wrapper.EventName {
	case "MaterialCreated":
		e := domain.MaterialCreated{}

		query.Decode(f, &mapped, &e)

		return e, nil

	case "MaterialNameChanged":
		e := domain.MaterialNameChanged{}

		query.Decode(f, &mapped, &e)

		return e, nil

	case "MaterialPriceChanged":
		e := domain.MaterialPriceChanged{}

		query.Decode(f, &mapped, &e)

		return e, nil

	case "MaterialQuantityChanged":
		e := domain.MaterialQuantityChanged{}

		query.Decode(f, &mapped, &e)

		return e, nil

	case "MaterialTypeChanged":
		e := domain.MaterialTypeChanged{}

		query.Decode(f, &mapped, &e)

		return e, nil

	case "MaterialExpirationDateChanged":
		e := domain.MaterialExpirationDateChanged{}

		query.Decode(f, &mapped, &e)

		return e, nil

	case "MaterialNotesChanged":
		e := domain.MaterialNotesChanged{}

		query.Decode(f, &mapped, &e)

		return e, nil

	case "MaterialProducedByChanged":
		e := domain.MaterialProducedByChanged{}

		query.Decode(f, &mapped, &e)

		return e, nil
	}

	return nil, nil
}
