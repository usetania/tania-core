package sqlite

import (
	"errors"
	"time"

	"github.com/Tanibox/tania-core/src/assets/domain"
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

func makeMaterialPricePerUnit(v interface{}) (domain.PricePerUnit, error) {
	mapped := v.(map[string]interface{})

	amount := ""
	currencyCode := ""
	for i, v2 := range mapped {
		if i == "amount" {
			val := v2.(string)
			amount = val
		}
		if i == "code" {
			val := v2.(string)
			currencyCode = val
		}
	}

	return domain.PricePerUnit{
		Amount:       amount,
		CurrencyCode: currencyCode,
	}, nil
}

func makeMaterialType(v interface{}) (domain.MaterialType, error) {
	mapped := v.(map[string]interface{})

	switch mapped["Type"] {
	case domain.MaterialTypePlantCode:
		mapped2 := mapped["Data"].(map[string]interface{})
		mapped3 := mapped2["PlantType"].(map[string]interface{})
		typeCode := mapped3["code"].(string)

		t, err := domain.CreateMaterialTypePlant(typeCode)
		if err != nil {
			return nil, err
		}

		return t, nil

	case domain.MaterialTypeSeedCode:
		mapped2 := mapped["Data"].(map[string]interface{})
		mapped3 := mapped2["PlantType"].(map[string]interface{})
		typeCode := mapped3["code"].(string)

		t, err := domain.CreateMaterialTypeSeed(typeCode)
		if err != nil {
			return nil, err
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
			return nil, err
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

	return nil, nil
}

func makeMaterialQuantity(v interface{}, materialTypeCode string) (domain.MaterialQuantity, error) {
	mapped := v.(map[string]interface{})

	value := float32(0)
	qtyUnit := domain.MaterialQuantityUnit{}
	for i, v2 := range mapped {
		if i == "value" {
			val := v2.(float64)
			value = float32(val)
		}
		if i == "unit" {
			mapped2 := v2.(map[string]interface{})
			unitCode := mapped2["code"].(string)
			u := domain.GetMaterialQuantityUnit(materialTypeCode, unitCode)

			qtyUnit = u
		}
	}

	return domain.MaterialQuantity{Unit: qtyUnit, Value: value}, nil
}
