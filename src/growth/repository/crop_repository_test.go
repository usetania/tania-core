package repository

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Tanibox/tania-server/src/growth/domain"
	"github.com/Tanibox/tania-server/src/growth/storage"
	deadlock "github.com/sasha-s/go-deadlock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CropServiceMock struct {
	mock.Mock
}

func (m CropServiceMock) FindInventoryMaterialByID(uid uuid.UUID) domain.ServiceResult {
	args := m.Called(uid)
	return args.Get(0).(domain.ServiceResult)
}
func (m CropServiceMock) FindByBatchID(batchID string) domain.ServiceResult {
	args := m.Called(batchID)
	return args.Get(0).(domain.ServiceResult)
}
func (m CropServiceMock) FindAreaByID(uid uuid.UUID) domain.ServiceResult {
	args := m.Called(uid)
	return args.Get(0).(domain.ServiceResult)
}

func TestCropInMemorySave(t *testing.T) {
	// Given
	done := make(chan bool)

	rwMutex := deadlock.RWMutex{}
	cropStorage := storage.CropStorage{CropMap: make(map[uuid.UUID]domain.Crop), Lock: &rwMutex}
	repo := NewCropRepositoryInMemory(&cropStorage)

	// Prepare mock
	cropServiceMock := new(CropServiceMock)
	areaAUID, _ := uuid.NewV4()
	areaBUID, _ := uuid.NewV4()
	areaAServiceResult := domain.ServiceResult{
		Result: domain.CropArea{UID: areaAUID, Type: "seeding"},
	}
	areaBServiceResult := domain.ServiceResult{
		Result: domain.CropArea{UID: areaBUID, Type: "growing"},
	}
	cropServiceMock.On("FindAreaByID", areaAUID).Return(areaAServiceResult)
	cropServiceMock.On("FindAreaByID", areaBUID).Return(areaBServiceResult)

	inventoryUID, _ := uuid.NewV4()
	inventoryServiceResult := domain.ServiceResult{
		Result: domain.CropInventory{
			UID:     inventoryUID,
			Variety: "Tomato Super One",
		},
	}
	cropServiceMock.On("FindInventoryMaterialByID", inventoryUID).Return(inventoryServiceResult)

	date := strings.ToLower(time.Now().Format("2Jan"))
	batchID := fmt.Sprintf("%s%s", "tom-sup-one-", date)
	cropServiceMock.On("FindByBatchID", batchID).Return(domain.ServiceResult{})

	containerType := domain.Tray{Cell: 15}

	crop1, cropErr1 := domain.CreateCropBatch(cropServiceMock, areaAUID, "SEEDING", inventoryUID, 20, containerType)
	crop2, cropErr2 := domain.CreateCropBatch(cropServiceMock, areaBUID, "GROWING", inventoryUID, 20, containerType)

	// When
	var err1, err2 error
	go func() {
		err1 = <-repo.Save(&crop1)
		err2 = <-repo.Save(&crop2)

		done <- true
	}()

	// Then
	<-done
	cropServiceMock.AssertExpectations(t)

	assert.Nil(t, cropErr1)
	assert.Nil(t, cropErr2)

	assert.Nil(t, err1)
	assert.Nil(t, err2)
}

func TestCropInMemoryFindAll(t *testing.T) {
	// Given
	done := make(chan bool)

	rwMutex := deadlock.RWMutex{}
	cropStorage := storage.CropStorage{CropMap: make(map[uuid.UUID]domain.Crop), Lock: &rwMutex}
	repo := NewCropRepositoryInMemory(&cropStorage)

	// Prepare mock
	cropServiceMock := new(CropServiceMock)
	areaAUID, _ := uuid.NewV4()
	areaBUID, _ := uuid.NewV4()
	areaAServiceResult := domain.ServiceResult{
		Result: domain.CropArea{UID: areaAUID, Type: "seeding"},
	}
	areaBServiceResult := domain.ServiceResult{
		Result: domain.CropArea{UID: areaBUID, Type: "growing"},
	}
	cropServiceMock.On("FindAreaByID", areaAUID).Return(areaAServiceResult)
	cropServiceMock.On("FindAreaByID", areaBUID).Return(areaBServiceResult)

	inventoryUID, _ := uuid.NewV4()
	inventoryServiceResult := domain.ServiceResult{
		Result: domain.CropInventory{
			UID:     inventoryUID,
			Variety: "Tomato Super One",
		},
	}
	cropServiceMock.On("FindInventoryMaterialByID", inventoryUID).Return(inventoryServiceResult)

	date := strings.ToLower(time.Now().Format("2Jan"))
	batchID := fmt.Sprintf("%s%s", "tom-sup-one-", date)
	cropServiceMock.On("FindByBatchID", batchID).Return(domain.ServiceResult{})

	containerType := domain.Tray{Cell: 15}

	crop1, cropErr1 := domain.CreateCropBatch(cropServiceMock, areaAUID, "SEEDING", inventoryUID, 20, containerType)
	crop2, cropErr2 := domain.CreateCropBatch(cropServiceMock, areaBUID, "GROWING", inventoryUID, 20, containerType)

	var result, foundOne RepositoryResult
	go func() {
		// Given
		<-repo.Save(&crop1)
		<-repo.Save(&crop2)

		// When
		result = <-repo.FindAll()

		val := result.Result.([]domain.Crop)
		foundOne = <-repo.FindByID(val[0].UID.String())

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, cropErr1)
	assert.Nil(t, cropErr2)

	val1, ok := result.Result.([]domain.Crop)
	assert.Equal(t, ok, true)
	assert.Equal(t, 2, len(val1))

	val2, ok := foundOne.Result.(domain.Crop)
	assert.Equal(t, ok, true)
	assert.Equal(t, val2.UID, val1[0].UID)
}

func TestCropInMemoryFindByID(t *testing.T) {
	// Given
	done := make(chan bool)

	rwMutex := deadlock.RWMutex{}
	cropStorage := storage.CropStorage{CropMap: make(map[uuid.UUID]domain.Crop), Lock: &rwMutex}
	repo := NewCropRepositoryInMemory(&cropStorage)

	// Prepare mock
	cropServiceMock := new(CropServiceMock)
	areaAUID, _ := uuid.NewV4()
	areaBUID, _ := uuid.NewV4()
	areaAServiceResult := domain.ServiceResult{
		Result: domain.CropArea{UID: areaAUID, Type: "seeding"},
	}
	areaBServiceResult := domain.ServiceResult{
		Result: domain.CropArea{UID: areaBUID, Type: "growing"},
	}
	cropServiceMock.On("FindAreaByID", areaAUID).Return(areaAServiceResult)
	cropServiceMock.On("FindAreaByID", areaBUID).Return(areaBServiceResult)

	inventoryUID, _ := uuid.NewV4()
	inventoryServiceResult := domain.ServiceResult{
		Result: domain.CropInventory{
			UID:     inventoryUID,
			Variety: "Tomato Super One",
		},
	}
	cropServiceMock.On("FindInventoryMaterialByID", inventoryUID).Return(inventoryServiceResult)

	date := strings.ToLower(time.Now().Format("2Jan"))
	batchID := fmt.Sprintf("%s%s", "tom-sup-one-", date)
	cropServiceMock.On("FindByBatchID", batchID).Return(domain.ServiceResult{})

	containerType := domain.Tray{Cell: 15}

	crop1, cropErr1 := domain.CreateCropBatch(cropServiceMock, areaAUID, "SEEDING", inventoryUID, 20, containerType)
	crop2, cropErr2 := domain.CreateCropBatch(cropServiceMock, areaBUID, "GROWING", inventoryUID, 20, containerType)

	var found1, found2 RepositoryResult
	go func() {
		// Given
		<-repo.Save(&crop1)
		<-repo.Save(&crop2)

		// When
		found1 = <-repo.FindByID(crop1.UID.String())
		found2 = <-repo.FindByID(crop2.UID.String())

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, cropErr1)
	assert.Nil(t, cropErr2)

	cropResult1 := found1.Result.(domain.Crop)
	assert.NotNil(t, cropResult1.UID)
	assert.Equal(t, cropResult1.InitialArea.AreaUID, areaAUID)

	cropResult2 := found2.Result.(domain.Crop)
	assert.NotNil(t, cropResult2.UID)
	assert.Equal(t, cropResult2.InitialArea.AreaUID, areaBUID)
}
