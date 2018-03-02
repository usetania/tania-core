package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Tanibox/tania-server/src/growth/domain"
	"github.com/Tanibox/tania-server/src/growth/query"
	"github.com/Tanibox/tania-server/src/growth/storage"
	"github.com/Tanibox/tania-server/src/growth/util/decoder"
	"github.com/mitchellh/mapstructure"
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

	f := mapstructure.ComposeDecodeHookFunc(
		decoder.UIDHook(),
		decoder.TimeHook(time.RFC3339),
		decoder.CropContainerHook(),
	)

	switch wrapper.EventName {
	case "CropBatchCreated":
		e := domain.CropBatchCreated{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "CropBatchTypeChanged":
		e := domain.CropBatchTypeChanged{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "CropBatchInventoryChanged":
		e := domain.CropBatchInventoryChanged{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "CropBatchContainerChanged":
		e := domain.CropBatchContainerChanged{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "CropBatchMoved":
		e := domain.CropBatchMoved{}

		decoder.Decode(f, &mapped, &e)

		// This decoding is too complex so we do this here instead in DecodeHookFunc
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

		decoder.Decode(f, &mapped, &e)

		// This decoding is too complex so we do this here instead in DecodeHookFunc
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
			if v2, ok2 := mapped2["created_date"]; ok2 {
				fmt.Println("MASUK SINI GA SIH")
				val, err := makeTime(v2)
				if err != nil {
					return nil, err
				}

				harvestedStorage.CreatedDate = val
			}
			if v2, ok2 := mapped2["last_updated"]; ok2 {
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

		return e, nil

	case "CropBatchDumped":
		e := domain.CropBatchDumped{}

		decoder.Decode(f, &mapped, &e)

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
			if v2, ok2 := mapped2["created_date"]; ok2 {
				val, err := makeTime(v2)
				if err != nil {
					return nil, err
				}

				trash.CreatedDate = val
			}
			if v2, ok2 := mapped2["last_updated"]; ok2 {
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

		fmt.Println("DUMPED E", e)

		return e, nil

	case "CropBatchWatered":
		e := domain.CropBatchWatered{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "CropBatchPhotoCreated":
		e := domain.CropBatchPhotoCreated{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "CropBatchNoteCreated":
		e := domain.CropBatchNoteCreated{}

		decoder.Decode(f, &mapped, &e)

		return e, nil

	case "CropBatchNoteRemoved":
		e := domain.CropBatchNoteRemoved{}

		decoder.Decode(f, &mapped, &e)

		return e, nil
	}

	return nil, nil
}
