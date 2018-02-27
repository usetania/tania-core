package sqlite

import (
	"errors"
	"time"

	"github.com/Tanibox/tania-server/src/assets/domain"
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

func makeWaterSource(v interface{}) (domain.WaterSource, error) {
	ws, ok := v.(map[string]interface{})
	if !ok {
		return nil, errors.New("Internal server error. Error type assertion")
	}

	convertedMap := map[string]float64{}
	for i, v := range ws {
		convertedMap[i] = v.(float64)
	}

	if convertedMap["Capacity"] == 0 {
		return domain.Tap{}, nil
	}

	return domain.Bucket{Capacity: float32(convertedMap["Capacity"])}, nil
}
