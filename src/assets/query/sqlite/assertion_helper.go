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

func makeAreaSize(v interface{}) (domain.AreaSize, error) {
	s, ok := v.(map[string]interface{})
	if !ok {
		return domain.AreaSize{}, errors.New("Internal server error. Error type assertion")
	}

	size := float32(0)
	sizeUnit := domain.AreaUnit{}
	for i, v := range s {
		if i == "unit" {
			sizeUnitMap := v.(map[string]interface{})

			for i2, v2 := range sizeUnitMap {
				if i2 == "symbol" {
					sizeUnit.Symbol = v2.(string)
				}
				if i2 == "label" {
					sizeUnit.Label = v2.(string)
				}
			}
		}
		if i == "value" {
			size64, ok := v.(float64)
			if !ok {
				return domain.AreaSize{}, errors.New("Internal server error. Error type assertion")
			}

			size = float32(size64)
		}
	}

	return domain.AreaSize{
		Value: size,
		Unit:  sizeUnit,
	}, nil
}

func makeAreaType(v interface{}) (domain.AreaType, error) {
	t, ok := v.(map[string]interface{})
	if !ok {
		return domain.AreaType{}, errors.New("Internal server error. Error type assertion")
	}

	convertedMap := map[string]string{}
	for i, v := range t {
		convertedMap[i] = v.(string)
	}

	return domain.GetAreaType(convertedMap["Code"]), nil
}

func makeAreaLocation(v interface{}) (domain.AreaLocation, error) {
	l, ok := v.(map[string]interface{})
	if !ok {
		return domain.AreaLocation{}, errors.New("Internal server error. Error type assertion")
	}

	convertedMap := map[string]string{}
	for i, v := range l {
		convertedMap[i] = v.(string)
	}

	return domain.GetAreaLocation(convertedMap["Code"]), nil
}
