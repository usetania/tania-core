package domain_test

import (
	"testing"
	"time"

	. "github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/service"
	"github.com/Tanibox/tania-server/src/assets/storage"
	uuid "github.com/satori/go.uuid"
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

	cropStorage := storage.CropStorage{CropMap: make(map[uuid.UUID]Crop)}
	cropQuery := query.NewCropQueryInMemory(&cropStorage)
	cropService := service.CropService{CropQuery: cropQuery}

	// When
	cropBatch1, errCrop1 := CreateCropBatch(areaNursery)
	cropBatch2, errCrop2 := CreateCropBatch(areaGrowing)

	errType1 := cropBatch1.ChangeCropType(Nursery{})
	errType2 := cropBatch2.ChangeCropType(Growing{})

	inventory1 := InventoryMaterial{PlantType: Vegetable{}, Variety: "Sawi"}
	errPlantType1 := cropService.ChangeInventory(&cropBatch1, inventory1)
	errPlantType2 := cropService.ChangeInventory(&cropBatch2, inventory1)

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

func TestChangeInventoryWithBatchID(t *testing.T) {
	// Given
	time, timeErr := time.Parse(time.RFC3339, "2018-01-25T22:08:41+07:00")

	farm, farmErr := CreateFarm("MyFarm1", "organic")
	area, areaErr := CreateArea(farm, "AreaNursery", "nursery")

	cropBatch, cropErr := CreateCropBatch(area)
	cropBatch.CreatedDate = time

	inventory := InventoryMaterial{PlantType: Vegetable{}, Variety: "Sawi Putih Super"}

	cropStorage := storage.CropStorage{CropMap: make(map[uuid.UUID]Crop)}
	cropQuery := query.NewCropQueryInMemory(&cropStorage)
	cropService := service.CropService{CropQuery: cropQuery}

	// When
	err := cropService.ChangeInventory(&cropBatch, inventory)

	// Then
	assert.Nil(t, timeErr)
	assert.Nil(t, farmErr)
	assert.Nil(t, areaErr)
	assert.Nil(t, cropErr)
	assert.Nil(t, err)

	assert.Equal(t, "Sawi Putih Super", cropBatch.Inventory.Variety)
	assert.Equal(t, "saw-put-sup-25jan", cropBatch.BatchID)
}

func TestCropCreateRemoveNote(t *testing.T) {
	// Given
	farm, farmErr := CreateFarm("MyFarm1", "organic")
	area, areaErr := CreateArea(farm, "Area1", "nursery")
	crop, cropErr := CreateCropBatch(area)

	// When
	crop.AddNewNote("This is my new note")

	// Then
	assert.Nil(t, farmErr)
	assert.Nil(t, areaErr)
	assert.Nil(t, cropErr)

	assert.Equal(t, 1, len(crop.Notes))

	uid := uuid.UUID{}
	for k, v := range crop.Notes {
		assert.Equal(t, "This is my new note", v.Content)
		assert.NotNil(t, v.CreatedDate)
		uid = k
	}

	// When
	crop.RemoveNote(uid.String())

	assert.Equal(t, 0, len(crop.Notes))
}
