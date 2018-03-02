package decoder

import (
	"reflect"
	"time"

	"github.com/Tanibox/tania-server/src/tasks/domain"
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

func TaskDomainDetailHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f != reflect.TypeOf(map[string]interface{}{}) {
			return data, nil
		}

		// reflect.TypeOf((*domain.TaskDomain)(nil)).Elem() is to find
		// the reflect.Type from interface variable.
		if t != reflect.TypeOf((*domain.TaskDomain)(nil)).Elem() {
			return data, nil
		}

		var domainDetails domain.TaskDomain
		if v, ok := data.(map[string]interface{}); ok {
			if v2, ok2 := v["InventoryUID"]; ok2 {
				u, ok3 := v2.(string)
				if !ok3 {
					return data, nil
				}

				uid, err := uuid.FromString(u)
				if err != nil {
					return data, err
				}

				domainDetails = domain.TaskDomainCrop{InventoryUID: &uid}
			}
		}

		return domainDetails, nil
	}
}
