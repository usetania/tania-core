package decoder

import (
	"reflect"
	"time"

	"github.com/gofrs/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/usetania/tania-core/src/growth/domain"
)

// InterfaceWrapper is used to wrap an interface with its struct name,
// so it will be easier to unmarshal later.
type InterfaceWrapper struct {
	Name string
	Data interface{}
}

func Decode(f mapstructure.DecodeHookFunc, data *map[string]interface{}, e interface{}) (interface{}, error) {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       f,
		TagName:          "json",
		Result:           e,
		WeaklyTypedInput: true,
	})
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(data)
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

func CropContainerHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f != reflect.TypeOf(map[string]interface{}{}) {
			return data, nil
		}

		if t != reflect.TypeOf(domain.CropContainer{}) {
			return data, nil
		}

		mapped, ok := data.(map[string]interface{})
		if !ok {
			return data, nil
		}

		var containerType domain.CropContainerType

		quantity := 0

		if v, ok := mapped["Quantity"]; ok {
			qty, ok2 := v.(float64)
			if !ok2 {
				return data, nil
			}

			quantity = int(qty)
		}

		if v, ok := mapped["Type"]; ok {
			mapped2, ok2 := v.(map[string]interface{})
			if !ok2 {
				return data, nil
			}

			if v2, ok2 := mapped2["Cell"]; ok2 {
				cell, ok3 := v2.(float64)
				if !ok3 {
					return data, nil
				}

				cellInt := int(cell)

				if cellInt != 0 {
					containerType = domain.Tray{Cell: cellInt}
				}
			} else {
				containerType = domain.Pot{}
			}
		}

		return domain.CropContainer{
			Quantity: quantity,
			Type:     containerType,
		}, nil
	}
}
