package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/usetania/tania-core/src/assets/domain"
)

func TestFindAllFarmTypes(t *testing.T) {
	t.Parallel()

	types := []FarmType{
		{Code: "organic", Name: "Organic / Soil-Based"},
		{Code: "hydroponic", Name: "Hydroponic"},
		{Code: "aquaponic", Name: "Aquaponic"},
		{Code: "mushroom", Name: "Mushroom"},
		{Code: "livestock", Name: "Livestock"},
		{Code: "fisheries", Name: "Fisheries"},
		{Code: "permaculture", Name: "Permaculture"},
	}

	farmTypes := FindAllFarmTypes()

	assert.Equal(t, types, farmTypes)
}

func TestFindFarmTypeByCode(t *testing.T) {
	t.Parallel()
	// Given
	farmType := FarmType{Code: "organic", Name: "Organic / Soil-Based"}

	result, err := FindFarmTypeByCode(farmType.Code)

	assert.Nil(t, err)
	assert.Equal(t, farmType, result)
}

func TestInvalidFindFarmTypeByCode(t *testing.T) {
	t.Parallel()

	code := "er"

	_, err := FindFarmTypeByCode(code)

	assert.Equal(t, FarmError{FarmErrorInvalidFarmTypeCode}, err)
}
