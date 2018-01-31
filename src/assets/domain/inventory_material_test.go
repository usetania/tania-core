package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateInventorySeed(t *testing.T) {
	// Given

	// When
	mts, err1 := CreateMaterialTypeSeed(PlantTypeVegetable)
	material, err2 := CreateMaterial("Bayam Lu Hsieh", "12", MoneyEUR, mts, 20, MaterialUnitPackets)
	tp, ok := material.Type.(MaterialTypeSeed)

	// Then
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Equal(t, "Bayam Lu Hsieh", material.Name)
	assert.Equal(t, "12", material.PricePerUnit.Amount())
	assert.Equal(t, true, ok)
	assert.Equal(t, PlantTypeVegetable, tp.PlantType.Code)
}
