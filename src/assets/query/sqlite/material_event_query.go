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
			fmt.Println("EVENT", event)
			fmt.Println("ERR", err)
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
			mapped := v.(map[string]interface{})

			for i2, v2 := range mapped {
				if i2 == "amount" {
					val := v2.(string)
					e.PricePerUnit.Amount = val
				}
				if i2 == "code" {
					val := v2.(string)
					e.PricePerUnit.CurrencyCode = val
				}
			}
		}

		if v, ok := mapped["Type"]; ok {
			mapped := v.(map[string]interface{})

			switch mapped["Type"] {
			case domain.MaterialTypePlantCode:
				mapped2 := mapped["Data"].(map[string]interface{})
				mapped3 := mapped2["PlantType"].(map[string]interface{})
				typeCode := mapped3["code"].(string)

				t, err := domain.CreateMaterialTypePlant(typeCode)
				if err != nil {
					return nil, err
				}

				e.Type = t

			case domain.MaterialTypeSeedCode:
				mapped2 := mapped["Data"].(map[string]interface{})
				mapped3 := mapped2["PlantType"].(map[string]interface{})
				typeCode := mapped3["code"].(string)

				t, err := domain.CreateMaterialTypeSeed(typeCode)
				if err != nil {
					return nil, err
				}

				e.Type = t

			case domain.MaterialTypeGrowingMediumCode:
				e.Type = domain.MaterialTypeGrowingMedium{}

			case domain.MaterialTypeAgrochemicalCode:
				mapped2 := mapped["Data"].(map[string]interface{})
				mapped3 := mapped2["ChemicalType"].(map[string]interface{})
				typeCode := mapped3["code"].(string)

				t, err := domain.CreateMaterialTypeAgrochemical(typeCode)
				if err != nil {
					return nil, err
				}

				e.Type = t

			case domain.MaterialTypeLabelAndCropSupportCode:
				e.Type = domain.MaterialTypeLabelAndCropSupport{}

			case domain.MaterialTypeSeedingContainerCode:
				mapped2 := mapped["Data"].(map[string]interface{})
				mapped3 := mapped2["ContainerType"].(map[string]interface{})
				typeCode := mapped3["code"].(string)

				t, err := domain.CreateMaterialTypeSeedingContainer(typeCode)
				if err != nil {
					return nil, err
				}

				e.Type = t

			case domain.MaterialTypePostHarvestSupplyCode:
				e.Type = domain.MaterialTypePostHarvestSupply{}

			case domain.MaterialTypeOtherCode:
				e.Type = domain.MaterialTypeOther{}

			}

			fmt.Println("E TYPE", e.Type.Code())
		}

		if v, ok := mapped["Quantity"]; ok {
			mapped := v.(map[string]interface{})

			for i2, v2 := range mapped {
				if i2 == "value" {
					val := v2.(float64)
					e.Quantity.Value = float32(val)
				}
				if i2 == "unit" {
					mapped2 := v2.(map[string]interface{})
					fmt.Println("MAPPED 2", mapped2)
					unitCode := mapped2["code"].(string)
					fmt.Println("UNIT CODE", unitCode)
					unit := domain.GetMaterialQuantityUnit(e.Type.Code(), unitCode)

					e.Quantity.Unit = unit
				}
			}
		}

		if v, ok := mapped["ExpirationDate"]; ok {
			if v != nil {
				d, err := makeTime(v)
				if err != nil {
					return nil, err
				}

				e.CreatedDate = d
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

	case "ReservoirWaterSourceChanged":
	}

	return nil, nil
}
