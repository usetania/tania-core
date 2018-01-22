package domain_test

import (
	"testing"

	. "github.com/Tanibox/tania-server/src/growth/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CropServiceMock struct {
	mock.Mock
}

func (m CropServiceMock) FindInventoryMaterialByID(uid uuid.UUID) ServiceResult {
	args := m.Called(uid)
	return args.Get(0).(ServiceResult)
}
func (m CropServiceMock) FindByBatchID(batchID string) ServiceResult {
	args := m.Called(batchID)
	return args.Get(0).(ServiceResult)
}
func (m CropServiceMock) FindAreaByID(uid uuid.UUID) ServiceResult {
	args := m.Called(uid)
	return args.Get(0).(ServiceResult)
}

func TestCreateCropBatch(t *testing.T) {
	// Given
	cropServiceMock := new(CropServiceMock)

	areaUID, _ := uuid.NewV4()
	areaServiceResult := ServiceResult{
		Result: CropArea{
			UID: areaUID,
		},
	}
	cropServiceMock.On("FindAreaByID", areaUID).Return(areaServiceResult)

	inventoryUID, _ := uuid.NewV4()
	inventoryServiceResult := ServiceResult{
		Result: CropInventory{
			UID: inventoryUID,
		},
	}
	cropServiceMock.On("FindInventoryMaterialByID", inventoryUID).Return(inventoryServiceResult)

	containerType := Tray{Cell: 15}

	// When
	crop, _ := CreateCropBatch(cropServiceMock, areaUID, "SEEDING", inventoryUID, 10, containerType)
	crop.Harvest(cropServiceMock, areaUID, 6)
	crop.Dump(cropServiceMock, areaUID, 1)

	// Then
	cropServiceMock.AssertExpectations(t)

	assert.Equal(t, CropActive, crop.Status.Code)
	assert.Equal(t, CropTypeSeeding, crop.Type.Code)
	assert.Equal(t, 3, crop.InitialArea.CurrentQuantity)

	// Harvest
	assert.Equal(t, areaUID, crop.HarvestedStorage[0].SourceAreaUID)
	assert.Equal(t, 6, crop.HarvestedStorage[0].Quantity)

	// Dump
	assert.Equal(t, areaUID, crop.Trash[0].SourceAreaUID)
	assert.Equal(t, 1, crop.Trash[0].Quantity)

}
