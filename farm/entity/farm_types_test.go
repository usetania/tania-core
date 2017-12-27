package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAllFarmTypes(t *testing.T) {
	types := []FarmType{
		FarmType{Code: "organic", Name: "Organic / Soil-Based"},
		FarmType{Code: "hydroponic", Name: "Hydroponic"},
		FarmType{Code: "aquaponic", Name: "Aquaponic"},
		FarmType{Code: "mushroom", Name: "Mushroom"},
		FarmType{Code: "livestock", Name: "Livestock"},
		FarmType{Code: "fisheries", Name: "Fisheries"},
		FarmType{Code: "permaculture", Name: "Permaculture"},
	}

	farmTypes := FindAllFarmTypes()

	assert.Equal(t, types, farmTypes)
}

func TestFindFarmTypeByCode(t *testing.T) {
	// Given
	farmType := FarmType{Code: "organic", Name: "Organic / Soil-Based"}

	result, err := FindFarmTypeByCode(farmType.Code)

	assert.Nil(t, err)
	assert.Equal(t, farmType, result)
}

func TestInvalidFindFarmTypeByCode(t *testing.T) {
	code := "er"

	_, err := FindFarmTypeByCode(code)

	assert.Equal(t, FarmError{FarmErrorInvalidFarmTypeCode}, err)
}
