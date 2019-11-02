package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateInventorySeed(t *testing.T) {
	// Given

	// When
	mts, err1 := CreateMaterialTypeSeed(PlantTypeVegetable)
	material1, err2 := CreateMaterial("Bayam Lu Hsieh", "12", MoneyEUR, mts, 20, MaterialUnitPackets, nil, nil, nil)
	tp, ok := material1.Type.(MaterialTypeSeed)

	// Then
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Equal(t, "Bayam Lu Hsieh", material1.Name)
	assert.Equal(t, "12", material1.PricePerUnit.Amount)
	assert.Equal(t, true, ok)
	assert.Equal(t, PlantTypeVegetable, tp.PlantType.Code)

	// When
	mta, err1 := CreateMaterialTypeAgrochemical(ChemicalTypeDisinfectant)
	material2, err2 := CreateMaterial("Green Disinfectant", "5", MoneyEUR, mta, 5, MaterialUnitPackets, nil, nil, nil)
	ta, ok := material2.Type.(MaterialTypeAgrochemical)

	// Then
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Equal(t, "Green Disinfectant", material2.Name)
	assert.Equal(t, "5", material2.PricePerUnit.Amount)
	assert.Equal(t, true, ok)
	assert.Equal(t, ChemicalTypeDisinfectant, ta.ChemicalType.Code)

	// When
	mtsc, err1 := CreateMaterialTypeSeedingContainer(ContainerTypeTray)
	material3, err2 := CreateMaterial("Soft Indoor Tray Pack", "10", MoneyEUR, mtsc, 10, MaterialUnitPieces, nil, nil, nil)
	tsc, ok := material3.Type.(MaterialTypeSeedingContainer)

	// Then
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Equal(t, "Soft Indoor Tray Pack", material3.Name)
	assert.Equal(t, "10", material3.PricePerUnit.Amount)
	assert.Equal(t, true, ok)
	assert.Equal(t, ContainerTypeTray, tsc.ContainerType.Code)

	// When
	mtgm := MaterialTypeGrowingMedium{}
	material4, err1 := CreateMaterial("Organic Super Soil", "2", MoneyEUR, mtgm, 5, MaterialUnitBags, nil, nil, nil)
	tgm, ok := material4.Type.(MaterialTypeGrowingMedium)

	// Then
	assert.Nil(t, err1)
	assert.Equal(t, "Organic Super Soil", material4.Name)
	assert.Equal(t, "2", material4.PricePerUnit.Amount)
	assert.Equal(t, true, ok)
	assert.Equal(t, MaterialTypeGrowingMediumCode, tgm.Code())

	// When
	mtl := MaterialTypeLabelAndCropSupport{}
	material5, err1 := CreateMaterial("Clean Label", "5", MoneyEUR, mtl, 5, MaterialUnitPieces, nil, nil, nil)
	tl, ok := material5.Type.(MaterialTypeLabelAndCropSupport)

	// Then
	assert.Nil(t, err1)
	assert.Equal(t, "Clean Label", material5.Name)
	assert.Equal(t, "5", material5.PricePerUnit.Amount)
	assert.Equal(t, true, ok)
	assert.Equal(t, MaterialTypeLabelAndCropSupportCode, tl.Code())

	// When
	mtph := MaterialTypePostHarvestSupply{}
	material6, err1 := CreateMaterial("Warm Solid Plastic", "5", MoneyEUR, mtph, 5, MaterialUnitPieces, nil, nil, nil)
	tph, ok := material6.Type.(MaterialTypePostHarvestSupply)

	// Then
	assert.Nil(t, err1)
	assert.Equal(t, "Warm Solid Plastic", material6.Name)
	assert.Equal(t, "5", material6.PricePerUnit.Amount)
	assert.Equal(t, true, ok)
	assert.Equal(t, MaterialTypePostHarvestSupplyCode, tph.Code())

	// When
	mto := MaterialTypeOther{}
	material7, err1 := CreateMaterial("Night Lamp Bright", "3", MoneyEUR, mto, 3, MaterialUnitPieces, nil, nil, nil)
	mo, ok := material7.Type.(MaterialTypeOther)

	// Then
	assert.Nil(t, err1)
	assert.Equal(t, "Night Lamp Bright", material7.Name)
	assert.Equal(t, "3", material7.PricePerUnit.Amount)
	assert.Equal(t, true, ok)
	assert.Equal(t, MaterialTypeOtherCode, mo.Code())
}
