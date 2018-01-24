package domain_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	. "github.com/Tanibox/tania-server/src/growth/domain"
	"github.com/Tanibox/tania-server/src/growth/query"
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
		Result: query.CropInventoryQueryResult{
			UID:     inventoryUID,
			Variety: "Tomato Super One",
		},
	}
	cropServiceMock.On("FindInventoryMaterialByID", inventoryUID).Return(inventoryServiceResult)

	date := strings.ToLower(time.Now().Format("2Jan"))
	batchID := fmt.Sprintf("%s%s", "tom-sup-one-", date)
	cropServiceMock.On("FindByBatchID", batchID).Return(ServiceResult{})

	containerType := Tray{Cell: 15}

	// When
	crop, _ := CreateCropBatch(cropServiceMock, areaAUID, CropTypeSeeding, inventoryUID, 20, containerType)
	crop.MoveToArea(cropServiceMock, areaAUID, areaBUID, 15)
	crop.Harvest(cropServiceMock, areaBUID, 6)
	crop.Dump(cropServiceMock, areaBUID, 5)
	crop.Fertilize()
	crop.Pesticide()
	crop.Prune()

	// Then
	cropServiceMock.AssertExpectations(t)

	assert.Equal(t, CropActive, crop.Status.Code)
	assert.Equal(t, CropTypeSeeding, crop.Type.Code)
	assert.Equal(t, 5, crop.InitialArea.CurrentQuantity)

	// Harvest
	assert.Equal(t, areaBUID, crop.HarvestedStorage[0].SourceAreaUID)
	assert.Equal(t, 6, crop.HarvestedStorage[0].Quantity)

	// Dump
	assert.Equal(t, areaBUID, crop.Trash[0].SourceAreaUID)
	assert.Equal(t, 5, crop.Trash[0].Quantity)

	// Care
	assert.Equal(t, false, crop.LastFertilized.IsZero())
	assert.Equal(t, false, crop.LastPesticided.IsZero())
	assert.Equal(t, false, crop.LastPruned.IsZero())
}
