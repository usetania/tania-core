package decoder

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/usetania/tania-core/src/growth/domain"
)

type CropEventWrapper InterfaceWrapper

func (w *CropEventWrapper) UnmarshalJSON(b []byte) error {
	wrapper := InterfaceWrapper{}

	err := json.Unmarshal(b, &wrapper)
	if err != nil {
		return err
	}

	mapped, ok := wrapper.Data.(map[string]interface{})
	if !ok {
		return errors.New("error type assertion")
	}

	f := mapstructure.ComposeDecodeHookFunc(
		UIDHook(),
		TimeHook(time.RFC3339),
		CropContainerHook(),
	)

	switch wrapper.Name {
	case "CropBatchCreated":
		e := domain.CropBatchCreated{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case "CropBatchTypeChanged":
		e := domain.CropBatchTypeChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case "CropBatchInventoryChanged":
		e := domain.CropBatchInventoryChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case "CropBatchContainerChanged":
		e := domain.CropBatchContainerChanged{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case "CropBatchMoved":
		e := domain.CropBatchMoved{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		// This decoding is too complex so we do this here instead in DecodeHookFunc
		if v, ok := mapped["UpdatedSrcArea"]; ok {
			code, ok2 := mapped["UpdatedSrcAreaCode"].(string)
			if !ok2 {
				return errors.New("error type assertion")
			}

			if code == "INITIAL_AREA" {
				initialArea, err := makeCropInitialArea(v)
				if err != nil {
					return err
				}

				e.UpdatedSrcArea = initialArea
			}

			if code == "MOVED_AREA" {
				movedArea, err := makeCropMovedArea(v)
				if err != nil {
					return err
				}

				e.UpdatedSrcArea = movedArea
			}
		}

		if v, ok := mapped["UpdatedDstArea"]; ok {
			code, ok2 := mapped["UpdatedDstAreaCode"].(string)
			if !ok2 {
				return errors.New("error type assertion")
			}

			if code == "INITIAL_AREA" {
				initialArea, err := makeCropInitialArea(v)
				if err != nil {
					return err
				}

				e.UpdatedDstArea = initialArea
			}

			if code == "MOVED_AREA" {
				movedArea, err := makeCropMovedArea(v)
				if err != nil {
					return err
				}

				e.UpdatedDstArea = movedArea
			}
		}

		w.Data = e

	case "CropBatchHarvested":
		e := domain.CropBatchHarvested{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		// This decoding is too complex so we do this here instead in DecodeHookFunc
		if v, ok := mapped["UpdatedHarvestedStorage"]; ok {
			mapped2, ok2 := v.(map[string]interface{})
			if !ok2 {
				return errors.New("error type assertion")
			}

			harvestedStorage := domain.HarvestedStorage{}

			if v2, ok2 := mapped2["quantity"]; ok2 {
				val, ok3 := v2.(float64)
				if !ok3 {
					return errors.New("error type assertion")
				}

				harvestedStorage.Quantity = int(val)
			}

			if v2, ok2 := mapped2["produced_gram_quantity"]; ok2 {
				val, ok3 := v2.(float64)
				if !ok3 {
					return errors.New("error type assertion")
				}

				harvestedStorage.ProducedGramQuantity = float32(val)
			}

			if v2, ok2 := mapped2["source_area_id"]; ok2 {
				uid, err := makeUUID(v2)
				if err != nil {
					return err
				}

				harvestedStorage.SourceAreaUID = uid
			}

			if v2, ok2 := mapped2["created_date"]; ok2 {
				val, err := makeTime(v2)
				if err != nil {
					return err
				}

				harvestedStorage.CreatedDate = val
			}

			if v2, ok2 := mapped2["last_updated"]; ok2 {
				val, err := makeTime(v2)
				if err != nil {
					return err
				}

				harvestedStorage.LastUpdated = val
			}

			e.UpdatedHarvestedStorage = harvestedStorage
		}

		if v, ok := mapped["HarvestedArea"]; ok {
			code, ok2 := mapped["HarvestedAreaCode"].(string)
			if !ok2 {
				return errors.New("error type assertion")
			}

			if code == "INITIAL_AREA" {
				initialArea, err := makeCropInitialArea(v)
				if err != nil {
					return err
				}

				e.HarvestedArea = initialArea
			}

			if code == "MOVED_AREA" {
				movedArea, err := makeCropMovedArea(v)
				if err != nil {
					return err
				}

				e.HarvestedArea = movedArea
			}
		}

		w.Data = e

	case "CropBatchDumped":
		e := domain.CropBatchDumped{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		if v, ok := mapped["UpdatedTrash"]; ok {
			mapped2, ok2 := v.(map[string]interface{})
			if !ok2 {
				return errors.New("error type assertion")
			}

			trash := domain.Trash{}

			if v2, ok2 := mapped2["quantity"]; ok2 {
				val, ok3 := v2.(float64)
				if !ok3 {
					return errors.New("error type assertion")
				}

				trash.Quantity = int(val)
			}

			if v2, ok2 := mapped2["source_area_id"]; ok2 {
				uid, err := makeUUID(v2)
				if err != nil {
					return err
				}

				trash.SourceAreaUID = uid
			}

			if v2, ok2 := mapped2["created_date"]; ok2 {
				val, err := makeTime(v2)
				if err != nil {
					return err
				}

				trash.CreatedDate = val
			}

			if v2, ok2 := mapped2["last_updated"]; ok2 {
				val, err := makeTime(v2)
				if err != nil {
					return err
				}

				trash.LastUpdated = val
			}

			e.UpdatedTrash = trash
		}

		if v, ok := mapped["DumpedArea"]; ok {
			code, ok2 := mapped["DumpedAreaCode"].(string)
			if !ok2 {
				return errors.New("error type assertion")
			}

			if code == "INITIAL_AREA" {
				initialArea, err := makeCropInitialArea(v)
				if err != nil {
					return err
				}

				e.DumpedArea = initialArea
			}

			if code == "MOVED_AREA" {
				movedArea, err := makeCropMovedArea(v)
				if err != nil {
					return err
				}

				e.DumpedArea = movedArea
			}
		}

		w.Data = e

	case "CropBatchWatered":
		e := domain.CropBatchWatered{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case "CropBatchPhotoCreated":
		e := domain.CropBatchPhotoCreated{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case "CropBatchNoteCreated":
		e := domain.CropBatchNoteCreated{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e

	case "CropBatchNoteRemoved":
		e := domain.CropBatchNoteRemoved{}

		_, err := Decode(f, &mapped, &e)
		if err != nil {
			return err
		}

		w.Data = e
	}

	return nil
}

func makeUUID(v interface{}) (uuid.UUID, error) {
	val, ok := v.(string)
	if !ok {
		return uuid.UUID{}, errors.New("error type assertion")
	}

	uid, err := uuid.FromString(val)
	if err != nil {
		return uuid.UUID{}, err
	}

	return uid, nil
}

func makeTime(v interface{}) (time.Time, error) {
	val, ok := v.(string)
	if !ok {
		return time.Time{}, errors.New("error type assertion")
	}

	date, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

func makeCropInitialArea(v interface{}) (domain.InitialArea, error) {
	initialArea := domain.InitialArea{}

	mapped, ok := v.(map[string]interface{})
	if !ok {
		return domain.InitialArea{}, errors.New("error type assertion")
	}

	if v, ok := mapped["area_id"]; ok {
		uid, err := makeUUID(v)
		if err != nil {
			return domain.InitialArea{}, err
		}

		initialArea.AreaUID = uid
	}

	if v, ok := mapped["initial_quantity"]; ok {
		qty, ok2 := v.(float64)
		if !ok2 {
			return domain.InitialArea{}, errors.New("error type assertion")
		}

		initialArea.InitialQuantity = int(qty)
	}

	if v, ok := mapped["current_quantity"]; ok {
		qty, ok2 := v.(float64)
		if !ok2 {
			return domain.InitialArea{}, errors.New("error type assertion")
		}

		initialArea.CurrentQuantity = int(qty)
	}

	if v, ok := mapped["created_date"]; ok {
		val, err := makeTime(v)
		if err != nil {
			return domain.InitialArea{}, err
		}

		initialArea.CreatedDate = val
	}

	if v, ok := mapped["last_updated"]; ok {
		val, err := makeTime(v)
		if err != nil {
			return domain.InitialArea{}, err
		}

		initialArea.LastUpdated = val
	}

	if v, ok := mapped["last_watered"]; ok {
		val, err := makeTime(v)
		if err != nil {
			return domain.InitialArea{}, err
		}

		initialArea.LastWatered = val
	}

	if v, ok := mapped["last_fertilized"]; ok {
		val, err := makeTime(v)
		if err != nil {
			return domain.InitialArea{}, err
		}

		initialArea.LastFertilized = val
	}

	if v, ok := mapped["last_pesticided"]; ok {
		val, err := makeTime(v)
		if err != nil {
			return domain.InitialArea{}, err
		}

		initialArea.LastPesticided = val
	}

	if v, ok := mapped["last_pruned"]; ok {
		val, err := makeTime(v)
		if err != nil {
			return domain.InitialArea{}, err
		}

		initialArea.LastPruned = val
	}

	return initialArea, nil
}

func makeCropMovedArea(v interface{}) (domain.MovedArea, error) {
	movedArea := domain.MovedArea{}

	mapped, ok := v.(map[string]interface{})
	if !ok {
		return domain.MovedArea{}, errors.New("error type assertion")
	}

	if v, ok := mapped["area_id"]; ok {
		uid, err := makeUUID(v)
		if err != nil {
			return domain.MovedArea{}, err
		}

		movedArea.AreaUID = uid
	}

	if v, ok := mapped["source_area_id"]; ok {
		uid, err := makeUUID(v)
		if err != nil {
			return domain.MovedArea{}, err
		}

		movedArea.SourceAreaUID = uid
	}

	if v, ok := mapped["initial_quantity"]; ok {
		qty, ok2 := v.(float64)
		if !ok2 {
			return domain.MovedArea{}, errors.New("error type assertion")
		}

		movedArea.InitialQuantity = int(qty)
	}

	if v, ok := mapped["current_quantity"]; ok {
		qty, ok2 := v.(float64)
		if !ok2 {
			if !ok2 {
				return domain.MovedArea{}, errors.New("error type assertion")
			}
		}

		movedArea.CurrentQuantity = int(qty)
	}

	if v, ok := mapped["created_date"]; ok {
		val, err := makeTime(v)
		if err != nil {
			return domain.MovedArea{}, err
		}

		movedArea.CreatedDate = val
	}

	if v, ok := mapped["last_updated"]; ok {
		val, err := makeTime(v)
		if err != nil {
			return domain.MovedArea{}, err
		}

		movedArea.LastUpdated = val
	}

	if v, ok := mapped["last_watered"]; ok {
		val, err := makeTime(v)
		if err != nil {
			return domain.MovedArea{}, err
		}

		movedArea.LastWatered = val
	}

	if v, ok := mapped["last_fertilized"]; ok {
		val, err := makeTime(v)
		if err != nil {
			return domain.MovedArea{}, err
		}

		movedArea.LastFertilized = val
	}

	if v, ok := mapped["last_pesticided"]; ok {
		val, err := makeTime(v)
		if err != nil {
			return domain.MovedArea{}, err
		}

		movedArea.LastPesticided = val
	}

	if v, ok := mapped["last_pruned"]; ok {
		val, err := makeTime(v)
		if err != nil {
			return domain.MovedArea{}, err
		}

		movedArea.LastPruned = val
	}

	return movedArea, nil
}
