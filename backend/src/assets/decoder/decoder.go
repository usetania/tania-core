package decoder

import (
	"reflect"
	"time"

	"github.com/gofrs/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/usetania/tania-core/src/assets/domain"
)

// EventWrapper is used to wrap the event interface with its struct name,
// so it will be easier to unmarshal later.
type EventWrapper struct {
	EventName string
	EventData interface{}
}

func Decode(f mapstructure.DecodeHookFunc, data *map[string]interface{}, e interface{}) (interface{}, error) {
	dc, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       f,
		TagName:          "json",
		Result:           e,
		WeaklyTypedInput: true,
	})
	if err != nil {
		return nil, err
	}

	err = dc.Decode(data)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func UIDHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(uuid.UUID{}) {
			return data, nil
		}

		return uuid.FromString(data.(string))
	}
}

func TimeHook(layout string) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		// Convert it by parsing
		return time.Parse(layout, data.(string))
	}
}

func WaterSourceHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f != reflect.TypeOf(map[string]interface{}{}) {
			return data, nil
		}

		// reflect.TypeOf((*domain.WaterSource)(nil)).Elem() is to find
		// the reflect.Type from interface variable.
		if t != reflect.TypeOf((*domain.WaterSource)(nil)).Elem() {
			return data, nil
		}

		mapped := data.(map[string]interface{})
		capacity := float32(0)

		// If Tap, it won't have Capacity, then it won't go inside loop.
		for key, val := range mapped {
			if key != "Capacity" {
				return data, nil
			}

			c := val.(float64)
			capacity = float32(c)
		}

		if capacity == 0 {
			return domain.Tap{}, nil
		}

		return domain.Bucket{Capacity: capacity}, nil
	}
}

func MaterialTypeHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f != reflect.TypeOf(map[string]interface{}{}) {
			return data, nil
		}

		// reflect.TypeOf((*domain.MaterialType)(nil)).Elem() is to find
		// the reflect.Type from interface variable.
		if t != reflect.TypeOf((*domain.MaterialType)(nil)).Elem() {
			return data, nil
		}

		mapped := data.(map[string]interface{})

		if val, ok := mapped["Type"]; ok {
			switch val {
			case domain.MaterialTypePlantCode:
				mapped2 := mapped["Data"].(map[string]interface{})
				mapped3 := mapped2["PlantType"].(map[string]interface{})
				typeCode := mapped3["code"].(string)

				t, err := domain.CreateMaterialTypePlant(typeCode)
				if err != nil {
					return data, err
				}

				return t, nil

			case domain.MaterialTypeSeedCode:
				mapped2 := mapped["Data"].(map[string]interface{})
				mapped3 := mapped2["PlantType"].(map[string]interface{})
				typeCode := mapped3["code"].(string)

				t, err := domain.CreateMaterialTypeSeed(typeCode)
				if err != nil {
					return data, err
				}

				return t, nil

			case domain.MaterialTypeGrowingMediumCode:
				return domain.MaterialTypeGrowingMedium{}, nil

			case domain.MaterialTypeAgrochemicalCode:
				mapped2 := mapped["Data"].(map[string]interface{})
				mapped3 := mapped2["ChemicalType"].(map[string]interface{})
				typeCode := mapped3["code"].(string)

				t, err := domain.CreateMaterialTypeAgrochemical(typeCode)
				if err != nil {
					return data, err
				}

				return t, nil

			case domain.MaterialTypeLabelAndCropSupportCode:
				return domain.MaterialTypeLabelAndCropSupport{}, nil

			case domain.MaterialTypeSeedingContainerCode:
				mapped2 := mapped["Data"].(map[string]interface{})
				mapped3 := mapped2["ContainerType"].(map[string]interface{})
				typeCode := mapped3["code"].(string)

				t, err := domain.CreateMaterialTypeSeedingContainer(typeCode)
				if err != nil {
					return nil, err
				}

				return t, nil

			case domain.MaterialTypePostHarvestSupplyCode:
				return domain.MaterialTypePostHarvestSupply{}, nil

			case domain.MaterialTypeOtherCode:
				return domain.MaterialTypeOther{}, nil
			}
		}

		return data, nil
	}
}
