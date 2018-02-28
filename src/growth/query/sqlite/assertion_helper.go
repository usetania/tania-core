package sqlite

import (
	"errors"
	"time"

	"github.com/Tanibox/tania-server/src/growth/domain"
	uuid "github.com/satori/go.uuid"
)

func makeUUID(v interface{}) (uuid.UUID, error) {
	val := v.(string)
	uid, err := uuid.FromString(val)
	if err != nil {
		return uuid.UUID{}, err
	}

	return uid, nil
}

func makeTime(v interface{}) (time.Time, error) {
	val := v.(string)

	date, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

func makeCropType(v interface{}) (domain.CropType, error) {
	mapped := v.(map[string]interface{})

	if v2, ok2 := mapped["Code"]; ok2 {
		c := v2.(string)
		return domain.GetCropType(c), nil
	}

	return domain.CropType{}, errors.New("Crop type not found")
}

func makeCropContainer(v interface{}) (domain.CropContainer, error) {
	mapped := v.(map[string]interface{})

	var containerType domain.CropContainerType
	quantity := 0

	if v, ok := mapped["Quantity"]; ok {
		qty := v.(float64)
		quantity = int(qty)
	}
	if v, ok := mapped["Type"]; ok {
		mapped2 := v.(map[string]interface{})

		if v2, ok2 := mapped2["Cell"]; ok2 {
			cell := v2.(float64)
			cellInt := int(cell)

			if cellInt == 0 {
				containerType = domain.Pot{}
			} else {
				containerType = domain.Tray{Cell: cellInt}
			}
		}
	}

	return domain.CropContainer{
		Quantity: quantity,
		Type:     containerType,
	}, nil
}

func makeCropInitialArea(v interface{}) (domain.InitialArea, error) {
	initialArea := domain.InitialArea{}
	mapped := v.(map[string]interface{})

	if v, ok := mapped["area_id"]; ok {
		uid, err := makeUUID(v)
		if err != nil {
			return domain.InitialArea{}, err
		}
		initialArea.AreaUID = uid
	}
	if v, ok := mapped["initial_quantity"]; ok {
		qty := v.(float64)
		initialArea.InitialQuantity = int(qty)
	}
	if v, ok := mapped["current_quantity"]; ok {
		qty := v.(float64)
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
	mapped := v.(map[string]interface{})

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
		qty := v.(float64)
		movedArea.InitialQuantity = int(qty)
	}
	if v, ok := mapped["current_quantity"]; ok {
		qty := v.(float64)
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
