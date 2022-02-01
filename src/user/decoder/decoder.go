package decoder

import (
	"encoding/base64"
	"reflect"
	"time"

	"github.com/gofrs/uuid"
	"github.com/mitchellh/mapstructure"
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

// PasswordHook is used because the password []byte data type is encoded to base64 string by json.Marshal
// so we need to decode it back to base64
// https://golang.org/pkg/encoding/json/#Marshal
func PasswordHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf([]byte{}) {
			return data, nil
		}

		return base64.StdEncoding.DecodeString(data.(string))
	}
}
