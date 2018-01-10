package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCropBatch(t *testing.T) {
	// Given
	farm, err := CreateFarm("MyFarm1", "organic")
	if err != nil {
		assert.Nil(t, err)
	}

	areaNursery, err1 := CreateArea(farm, "AreaNursery", "nursery")
	areaGrowing, err2 := CreateArea(farm, "AreaGrowing", "growing")

	// When
	cropBatch1, errCrop1 := CreateCropBatch(areaNursery)
	cropBatch2, errCrop2 := CreateCropBatch(areaGrowing)

	errType1 := cropBatch1.ChangeCropType(Nursery{})
	errType2 := cropBatch2.ChangeCropType(Growing{})

	inventory1 := InventoryMaterial{PlantType: Vegetable{}, Variety: "Sawi"}
	errPlantType1 := cropBatch1.ChangeInventory(inventory1)
	errPlantType2 := cropBatch2.ChangeInventory(inventory1)

	tray := CropContainer{Quantity: 10, Type: Tray{Cell: 20}}
	pot := CropContainer{Quantity: 50, Type: Pot{}}
	errContainer1 := cropBatch1.ChangeContainer(tray)
	errContainer2 := cropBatch2.ChangeContainer(pot)

	// Then
	assert.Nil(t, err1)
	assert.Nil(t, err2)

	assert.NotNil(t, cropBatch1)
	assert.Nil(t, errCrop1)
	assert.NotNil(t, cropBatch2)
	assert.Nil(t, errCrop2)

	assert.Nil(t, errType1)
	assert.Nil(t, errType2)

	assert.Nil(t, errPlantType1)
	assert.Nil(t, errPlantType2)

	assert.Nil(t, errContainer1)
	assert.Nil(t, errContainer2)
}
