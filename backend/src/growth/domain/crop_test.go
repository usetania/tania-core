package domain_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	. "github.com/usetania/tania-core/src/growth/domain"
	"github.com/usetania/tania-core/src/growth/query"
)

type CropServiceMock struct {
	mock.Mock
}

func (m *CropServiceMock) FindMaterialByID(uid uuid.UUID) ServiceResult {
	args := m.Called(uid)

	return args.Get(0).(ServiceResult)
}

func (m *CropServiceMock) FindByBatchID(batchID string) ServiceResult {
	args := m.Called(batchID)

	return args.Get(0).(ServiceResult)
}

func (m *CropServiceMock) FindAreaByID(uid uuid.UUID) ServiceResult {
	args := m.Called(uid)

	return args.Get(0).(ServiceResult)
}

func TestCreateCropBatch(t *testing.T) {
	t.Parallel()
	// Given
	cropServiceMock := new(CropServiceMock)

	areaAUID, _ := uuid.NewV4()
	areaBUID, _ := uuid.NewV4()
	areaAServiceResult := ServiceResult{
		Result: query.CropAreaQueryResult{UID: areaAUID, Type: "SEEDING"},
	}
	areaBServiceResult := ServiceResult{
		Result: query.CropAreaQueryResult{UID: areaBUID, Type: "GROWING"},
	}

	cropServiceMock.On("FindAreaByID", areaAUID).Return(areaAServiceResult)
	cropServiceMock.On("FindAreaByID", areaBUID).Return(areaBServiceResult)

	inventoryUID, _ := uuid.NewV4()
	inventoryServiceResult := ServiceResult{
		Result: query.CropMaterialQueryResult{
			UID:  inventoryUID,
			Name: "Tomato Super One",
		},
	}
	cropServiceMock.On("FindMaterialByID", inventoryUID).Return(inventoryServiceResult)

	date := strings.ToLower(time.Now().Format("2Jan"))
	batchID := fmt.Sprintf("%s%s", "tom-sup-one-", date)
	cropServiceMock.On("FindByBatchID", batchID).Return(ServiceResult{})

	containerType := Tray{Cell: 15}

	// When
	crop, _ := CreateCropBatch(cropServiceMock, areaAUID, CropTypeSeeding, inventoryUID, 20, containerType)
	crop.MoveToArea(cropServiceMock, areaAUID, areaBUID, 15)
	crop.Dump(cropServiceMock, areaBUID, 5, "Notes")
	crop.Fertilize()
	crop.Pesticide()
	crop.Prune()

	// Then
	cropServiceMock.AssertExpectations(t)

	assert.Equal(t, CropActive, crop.Status.Code)
	assert.Equal(t, CropTypeSeeding, crop.Type.Code)
	assert.Equal(t, 5, crop.InitialArea.CurrentQuantity)

	// Dump
	assert.Equal(t, areaBUID, crop.Trash[0].SourceAreaUID)
	assert.Equal(t, 5, crop.Trash[0].Quantity)
}

func TestHarvestCropBatch(t *testing.T) {
	t.Parallel()
	// Given
	cropServiceMock := new(CropServiceMock)

	areaAUID, _ := uuid.NewV4()
	areaBUID, _ := uuid.NewV4()
	areaAServiceResult := ServiceResult{
		Result: query.CropAreaQueryResult{UID: areaAUID, Type: "SEEDING"},
	}
	areaBServiceResult := ServiceResult{
		Result: query.CropAreaQueryResult{UID: areaBUID, Type: "GROWING"},
	}

	cropServiceMock.On("FindAreaByID", areaAUID).Return(areaAServiceResult)
	cropServiceMock.On("FindAreaByID", areaBUID).Return(areaBServiceResult)

	inventoryUID, _ := uuid.NewV4()
	inventoryServiceResult := ServiceResult{
		Result: query.CropMaterialQueryResult{
			UID:  inventoryUID,
			Name: "Tomato Super One",
		},
	}
	cropServiceMock.On("FindMaterialByID", inventoryUID).Return(inventoryServiceResult)

	date := strings.ToLower(time.Now().Format("2Jan"))
	batchID := fmt.Sprintf("%s%s", "tom-sup-one-", date)
	cropServiceMock.On("FindByBatchID", batchID).Return(ServiceResult{})

	containerType := Tray{Cell: 15}

	// When
	crop, _ := CreateCropBatch(cropServiceMock, areaAUID, CropTypeSeeding, inventoryUID, 20, containerType)
	crop.MoveToArea(cropServiceMock, areaAUID, areaBUID, 15)
	err1 := crop.Harvest(cropServiceMock, areaBUID, HarvestTypePartial, 10, GetProducedUnit(Kg), "Notes")
	err2 := crop.Harvest(cropServiceMock, areaAUID, HarvestTypePartial, 10, GetProducedUnit(Kg), "Notes")

	// Then
	cropServiceMock.AssertExpectations(t)

	assert.Nil(t, err1)
	assert.Equal(t, areaBUID, crop.HarvestedStorage[0].SourceAreaUID)
	assert.Equal(t, 0, crop.HarvestedStorage[0].Quantity)
	assert.Equal(t, float32(10000), crop.HarvestedStorage[0].ProducedGramQuantity)

	assert.NotNil(t, err2)

	// When
	crop.Harvest(cropServiceMock, areaBUID, HarvestTypeAll, 2000, GetProducedUnit(Gr), "Notes")

	// Then
	assert.Equal(t, 0, crop.MovedArea[0].CurrentQuantity)
	assert.Equal(t, 15, crop.HarvestedStorage[0].Quantity)
	assert.Equal(t, float32(12000), crop.HarvestedStorage[0].ProducedGramQuantity)
}

func TestWaterCrop(t *testing.T) {
	t.Parallel()
	// Given
	cropServiceMock := new(CropServiceMock)

	areaAUID, _ := uuid.NewV4()
	areaBUID, _ := uuid.NewV4()
	areaAServiceResult := ServiceResult{
		Result: query.CropAreaQueryResult{UID: areaAUID, Type: "SEEDING"},
	}
	areaBServiceResult := ServiceResult{
		Result: query.CropAreaQueryResult{UID: areaBUID, Type: "GROWING"},
	}

	cropServiceMock.On("FindAreaByID", areaAUID).Return(areaAServiceResult)
	cropServiceMock.On("FindAreaByID", areaBUID).Return(areaBServiceResult)

	inventoryUID, _ := uuid.NewV4()
	inventoryServiceResult := ServiceResult{
		Result: query.CropMaterialQueryResult{
			UID:  inventoryUID,
			Name: "Tomato Super One",
		},
	}
	cropServiceMock.On("FindMaterialByID", inventoryUID).Return(inventoryServiceResult)

	date := strings.ToLower(time.Now().Format("2Jan"))
	batchID := fmt.Sprintf("%s%s", "tom-sup-one-", date)
	cropServiceMock.On("FindByBatchID", batchID).Return(ServiceResult{})

	containerType := Tray{Cell: 15}

	wateringDate := "2018-Jan-15"
	wDate, _ := time.Parse("2006-Jan-02", wateringDate)

	// When
	crop, errCrop := CreateCropBatch(cropServiceMock, areaAUID, CropTypeSeeding, inventoryUID, 20, containerType)
	errMove := crop.MoveToArea(cropServiceMock, areaAUID, areaBUID, 15)
	errWater1 := crop.Water(cropServiceMock, areaAUID, wDate)
	errWater2 := crop.Water(cropServiceMock, areaBUID, wDate)

	// Then
	cropServiceMock.AssertExpectations(t)

	assert.Nil(t, errCrop)
	assert.Nil(t, errMove)
	assert.Nil(t, errWater1)
	assert.Nil(t, errWater2)
	assert.Equal(t, wDate, crop.InitialArea.LastWatered)
	assert.Equal(t, wDate, crop.MovedArea[0].LastWatered)
}

func TestCropHarvestArchiveStatus(t *testing.T) {
	t.Parallel()
	// Given
	cropServiceMock := new(CropServiceMock)

	areaAUID, _ := uuid.NewV4()
	areaBUID, _ := uuid.NewV4()
	areaAServiceResult := ServiceResult{
		Result: query.CropAreaQueryResult{UID: areaAUID, Type: "SEEDING"},
	}
	areaBServiceResult := ServiceResult{
		Result: query.CropAreaQueryResult{UID: areaBUID, Type: "GROWING"},
	}

	cropServiceMock.On("FindAreaByID", areaAUID).Return(areaAServiceResult)
	cropServiceMock.On("FindAreaByID", areaBUID).Return(areaBServiceResult)

	inventoryUID, _ := uuid.NewV4()
	inventoryServiceResult := ServiceResult{
		Result: query.CropMaterialQueryResult{
			UID:  inventoryUID,
			Name: "Tomato Super One",
		},
	}
	cropServiceMock.On("FindMaterialByID", inventoryUID).Return(inventoryServiceResult)

	date := strings.ToLower(time.Now().Format("2Jan"))
	batchID := fmt.Sprintf("%s%s", "tom-sup-one-", date)
	cropServiceMock.On("FindByBatchID", batchID).Return(ServiceResult{})

	containerType := Tray{Cell: 15}

	// When
	crop, _ := CreateCropBatch(cropServiceMock, areaAUID, CropTypeSeeding, inventoryUID, 20, containerType)
	crop.MoveToArea(cropServiceMock, areaAUID, areaBUID, 15)
	crop.Harvest(cropServiceMock, areaBUID, HarvestTypeAll, 2000, GetProducedUnit(Gr), "Notes")

	// Then
	assert.Equal(t, crop.Status.Code, CropActive)

	// When
	crop.MoveToArea(cropServiceMock, areaAUID, areaBUID, 5)
	crop.Harvest(cropServiceMock, areaBUID, HarvestTypeAll, 3000, GetProducedUnit(Gr), "Notes")

	// Then
	assert.Equal(t, crop.Status.Code, CropArchived)
}

func TestCropDumpArchiveStatus(t *testing.T) {
	t.Parallel()
	// Given
	cropServiceMock := new(CropServiceMock)

	areaAUID, _ := uuid.NewV4()
	areaBUID, _ := uuid.NewV4()
	areaAServiceResult := ServiceResult{
		Result: query.CropAreaQueryResult{UID: areaAUID, Type: "SEEDING"},
	}
	areaBServiceResult := ServiceResult{
		Result: query.CropAreaQueryResult{UID: areaBUID, Type: "GROWING"},
	}

	cropServiceMock.On("FindAreaByID", areaAUID).Return(areaAServiceResult)
	cropServiceMock.On("FindAreaByID", areaBUID).Return(areaBServiceResult)

	inventoryUID, _ := uuid.NewV4()
	inventoryServiceResult := ServiceResult{
		Result: query.CropMaterialQueryResult{
			UID:  inventoryUID,
			Name: "Tomato Super One",
		},
	}
	cropServiceMock.On("FindMaterialByID", inventoryUID).Return(inventoryServiceResult)

	date := strings.ToLower(time.Now().Format("2Jan"))
	batchID := fmt.Sprintf("%s%s", "tom-sup-one-", date)
	cropServiceMock.On("FindByBatchID", batchID).Return(ServiceResult{})

	containerType := Tray{Cell: 15}

	// When
	crop, _ := CreateCropBatch(cropServiceMock, areaAUID, CropTypeSeeding, inventoryUID, 20, containerType)
	crop.MoveToArea(cropServiceMock, areaAUID, areaBUID, 15)
	crop.Dump(cropServiceMock, areaBUID, 15, "Notes")

	// Then
	assert.Equal(t, crop.Status.Code, CropActive)

	// When
	crop.Dump(cropServiceMock, areaAUID, 5, "Notes")

	// Then
	assert.Equal(t, crop.Status.Code, CropArchived)
}
