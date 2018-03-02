package query

import (
	"reflect"
	"time"

	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
)

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

func WaterSourceHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f != reflect.TypeOf(map[string]interface{}{}) {
			return data, nil
		}

		mapped := data.(map[string]interface{})
		cap := float32(0)

		// If Tap, it won't have Capacity, then it won't go inside loop.
		for key, val := range mapped {
			if key != "Capacity" {
				return data, nil
			}

			c := val.(float64)
			cap = float32(c)
		}

		// reflect.TypeOf((*domain.WaterSource)(nil)).Elem() is to find
		// the reflect.Type from interface variable.
		if t != reflect.TypeOf((*domain.WaterSource)(nil)).Elem() {
			return data, nil
		}

		if cap == 0 {
			return domain.Tap{}, nil
		}

		return domain.Bucket{Capacity: cap}, nil
	}
}
