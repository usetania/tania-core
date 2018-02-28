package sqlite

import (
	"database/sql"
	"encoding/json"
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
			val, err := makeCropType(v)
			if err != nil {
				return nil, err
			}

			e.Type = val
		}
		if v, ok := mapped["Container"]; ok {
			val, err := makeCropContainer(v)
			if err != nil {
				return nil, err
			}

			e.Container = val
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

	case "CropBatchTypeChanged":
		e := domain.CropBatchTypeChanged{}

		if v, ok := mapped["UID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.UID = uid
		}
		if v, ok := mapped["Type"]; ok {
			val, err := makeCropType(v)
			if err != nil {
				return nil, err
			}

			e.Type = val
		}

		return e, nil

	case "CropBatchInventoryChanged":
		e := domain.CropBatchInventoryChanged{}

		if v, ok := mapped["UID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.UID = uid
		}
		if v, ok := mapped["InventoryUID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.InventoryUID = uid
		}
		if v, ok := mapped["BatchID"]; ok {
			val := v.(string)
			e.BatchID = val
		}

		return e, nil

	case "CropBatchContainerChanged":
		e := domain.CropBatchContainerChanged{}

		if v, ok := mapped["UID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.UID = uid
		}
		if v, ok := mapped["Container"]; ok {
			val, err := makeCropContainer(v)
			if err != nil {
				return nil, err
			}

			e.Container = val
		}

		return e, nil

	case "CropBatchMoved":
		e := domain.CropBatchMoved{}

		if v, ok := mapped["UID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.UID = uid
		}
		if v, ok := mapped["Quantity"]; ok {
			val := v.(float64)
			e.Quantity = int(val)
		}
		if v, ok := mapped["SrcAreaUID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.SrcAreaUID = uid
		}
		if v, ok := mapped["DstAreaUID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.DstAreaUID = uid
		}
		if v, ok := mapped["MovedDate"]; ok {
			val, err := makeTime(v)
			if err != nil {
				return nil, err
			}

			e.MovedDate = val
		}
		if v, ok := mapped["UpdatedSrcArea"]; ok {
			code := mapped["UpdatedSrcAreaCode"].(string)

			if code == "INITIAL_AREA" {
				initialArea, err := makeCropInitialArea(v)
				if err != nil {
					return nil, err
				}

				e.UpdatedSrcArea = initialArea
			}
			if code == "MOVED_AREA" {
				movedArea, err := makeCropMovedArea(v)
				if err != nil {
					return nil, err
				}

				e.UpdatedSrcArea = movedArea
			}
		}
		if v, ok := mapped["UpdatedDstArea"]; ok {
			code := mapped["UpdatedDstAreaCode"].(string)

			if code == "INITIAL_AREA" {
				initialArea, err := makeCropInitialArea(v)
				if err != nil {
					return nil, err
				}

				e.UpdatedDstArea = initialArea
			}
			if code == "MOVED_AREA" {
				movedArea, err := makeCropMovedArea(v)
				if err != nil {
					return nil, err
				}

				e.UpdatedDstArea = movedArea
			}
		}

		return e, nil

	case "CropBatchHarvested":
		e := domain.CropBatchHarvested{}

		if v, ok := mapped["UID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.UID = uid
		}
		if v, ok := mapped["HarvestType"]; ok {
			val := v.(string)
			e.HarvestType = val
		}
		if v, ok := mapped["HarvestedQuantity"]; ok {
			val := v.(float64)
			e.HarvestedQuantity = int(val)
		}
		if v, ok := mapped["ProducedGramQuantity"]; ok {
			val := v.(float64)
			e.ProducedGramQuantity = float32(val)
		}
		if v, ok := mapped["UpdatedHarvestedStorage"]; ok {
			mapped2 := v.(map[string]interface{})
			harvestedStorage := domain.HarvestedStorage{}

			if v2, ok2 := mapped2["quantity"]; ok2 {
				val := v2.(float64)
				harvestedStorage.Quantity = int(val)
			}
			if v2, ok2 := mapped2["produced_gram_quantity"]; ok2 {
				val := v2.(float64)
				harvestedStorage.ProducedGramQuantity = float32(val)
			}
			if v2, ok2 := mapped2["source_area_id"]; ok2 {
				uid, err := makeUUID(v2)
				if err != nil {
					return nil, err
				}

				harvestedStorage.SourceAreaUID = uid
			}
			if v2, ok2 := mapped["created_date"]; ok2 {
				val, err := makeTime(v2)
				if err != nil {
					return nil, err
				}

				harvestedStorage.CreatedDate = val
			}
			if v2, ok2 := mapped["last_updated"]; ok2 {
				val, err := makeTime(v2)
				if err != nil {
					return nil, err
				}

				harvestedStorage.LastUpdated = val
			}

			e.UpdatedHarvestedStorage = harvestedStorage
		}
		if v, ok := mapped["HarvestedArea"]; ok {
			code := mapped["HarvestedAreaCode"].(string)

			if code == "INITIAL_AREA" {
				initialArea, err := makeCropInitialArea(v)
				if err != nil {
					return nil, err
				}

				e.HarvestedArea = initialArea
			}
			if code == "MOVED_AREA" {
				movedArea, err := makeCropMovedArea(v)
				if err != nil {
					return nil, err
				}

				e.HarvestedArea = movedArea
			}
		}
		if v, ok := mapped["HarvestDate"]; ok {
			val, err := makeTime(v)
			if err != nil {
				return nil, err
			}

			e.HarvestDate = val
		}
		if v, ok := mapped["Notes"]; ok {
			val := v.(string)
			e.Notes = val
		}

		return e, nil

	case "CropBatchDumped":
		e := domain.CropBatchDumped{}

		if v, ok := mapped["UID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.UID = uid
		}
		if v, ok := mapped["Quantity"]; ok {
			val := v.(float64)
			e.Quantity = int(val)
		}
		if v, ok := mapped["UpdatedTrash"]; ok {
			mapped2 := v.(map[string]interface{})
			trash := domain.Trash{}

			if v2, ok2 := mapped2["quantity"]; ok2 {
				val := v2.(float64)
				trash.Quantity = int(val)
			}
			if v2, ok2 := mapped2["source_area_id"]; ok2 {
				uid, err := makeUUID(v2)
				if err != nil {
					return nil, err
				}

				trash.SourceAreaUID = uid
			}
			if v2, ok2 := mapped["created_date"]; ok2 {
				val, err := makeTime(v2)
				if err != nil {
					return nil, err
				}

				trash.CreatedDate = val
			}
			if v2, ok2 := mapped["last_updated"]; ok2 {
				val, err := makeTime(v2)
				if err != nil {
					return nil, err
				}

				trash.LastUpdated = val
			}

			e.UpdatedTrash = trash
		}
		if v, ok := mapped["DumpedArea"]; ok {
			code := mapped["DumpedAreaCode"].(string)

			if code == "INITIAL_AREA" {
				initialArea, err := makeCropInitialArea(v)
				if err != nil {
					return nil, err
				}

				e.DumpedArea = initialArea
			}
			if code == "MOVED_AREA" {
				movedArea, err := makeCropMovedArea(v)
				if err != nil {
					return nil, err
				}

				e.DumpedArea = movedArea
			}
		}
		if v, ok := mapped["DumpDate"]; ok {
			val, err := makeTime(v)
			if err != nil {
				return nil, err
			}

			e.DumpDate = val
		}
		if v, ok := mapped["Notes"]; ok {
			val := v.(string)
			e.Notes = val
		}

		return e, nil

	case "CropBatchPhotoCreated":
		e := domain.CropBatchPhotoCreated{}

		if v, ok := mapped["UID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.UID = uid
		}
		if v, ok := mapped["CropUID"]; ok {
			uid, err := makeUUID(v)
			if err != nil {
				return nil, err
			}

			e.CropUID = uid
		}
		if v, ok := mapped["Filename"]; ok {
			val := v.(string)
			e.Filename = val
		}
		if v, ok := mapped["MimeType"]; ok {
			val := v.(string)
			e.MimeType = val
		}
		if v, ok := mapped["Size"]; ok {
			val := v.(float64)
			e.Size = int(val)
		}
		if v, ok := mapped["Width"]; ok {
			val := v.(float64)
			e.Width = int(val)
		}
		if v, ok := mapped["Height"]; ok {
			val := v.(float64)
			e.Height = int(val)
		}
		if v, ok := mapped["Description"]; ok {
			val := v.(string)
			e.Description = val
		}
	}

	return nil, nil
}
