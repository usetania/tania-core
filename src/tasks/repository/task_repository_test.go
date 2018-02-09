package repository

import (
	"fmt"
	"github.com/Tanibox/tania-server/src/tasks/domain"
	"github.com/Tanibox/tania-server/src/tasks/query"
	"github.com/Tanibox/tania-server/src/tasks/storage"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync"
	"testing"
	"time"
)

type TaskServiceMock struct {
	mock.Mock
}

func (m TaskServiceMock) FindAreaByID(uid uuid.UUID) domain.ServiceResult {
	args := m.Called(uid)
	return args.Get(0).(domain.ServiceResult)
}
func (m TaskServiceMock) FindCropByID(uid uuid.UUID) domain.ServiceResult {
	args := m.Called(uid)
	return args.Get(0).(domain.ServiceResult)
}
func (m TaskServiceMock) FindMaterialByID(uid uuid.UUID) domain.ServiceResult {
	args := m.Called(uid)
	return args.Get(0).(domain.ServiceResult)
}

func TestTaskInMemoryFindWithFilter(t *testing.T) {
	taskServiceMock := new(TaskServiceMock)

	// Given
	done := make(chan bool)

	rwMutex := sync.RWMutex{}
	taskStorage := storage.TaskStorage{TaskMap: make(map[uuid.UUID]domain.Task), Lock: rwMutex}
	repo := NewTaskRepositoryInMemory(&taskStorage)

	title := "My Title"
	description := "My Description"
	due_date, _ := time.Parse(time.RFC3339, "2019-01-23T17:37:39.697328206+01:00")
	due_ptr := &due_date
	priority := "URGENT"

	inventoryID, _ := uuid.NewV4()

	taskServiceMock.On("FindMaterialByID", inventoryID).Return(domain.ServiceResult{Result: query.TaskMaterialQueryResult{UID: inventoryID}})
	taskdomain_crop, _ := domain.CreateTaskDomainCrop(taskServiceMock, inventoryID)

	taskdomain_area, _ := domain.CreateTaskDomainArea()

	category := "SANITATION"

	assetID, _ := uuid.NewV4()

	taskServiceMock.On("FindCropByID", assetID).Return(domain.ServiceResult{Result: query.TaskCropQueryResult{UID: inventoryID}})
	task1, taskErr1 := domain.CreateTask(
		taskServiceMock, title, description, due_ptr, priority, taskdomain_crop,
		category, &assetID)

	taskServiceMock.On("FindCropByID", assetID).Return(domain.ServiceResult{Result: query.TaskCropQueryResult{UID: inventoryID}})
	task2, taskErr2 := domain.CreateTask(
		taskServiceMock, title, description, due_ptr, priority, taskdomain_crop,
		category, &assetID)

	taskServiceMock.On("FindAreaByID", assetID).Return(domain.ServiceResult{Result: query.TaskCropQueryResult{UID: inventoryID}})
	task3, taskErr3 := domain.CreateTask(
		taskServiceMock, title, description, due_ptr, priority, taskdomain_area,
		category, &assetID)

	taskServiceMock.On("FindAreaByID", assetID).Return(domain.ServiceResult{Result: query.TaskCropQueryResult{UID: inventoryID}})
	task4, taskErr4 := domain.CreateTask(
		taskServiceMock, title, description, due_ptr, priority, taskdomain_area,
		category, &assetID)

	taskServiceMock.On("FindAreaByID", assetID).Return(domain.ServiceResult{Result: query.TaskCropQueryResult{UID: inventoryID}})
	task5, taskErr5 := domain.CreateTask(
		taskServiceMock, title, description, due_ptr, priority, taskdomain_area,
		category, &assetID)

	var result, foundOne RepositoryResult
	go func() {
		// Given
		<-repo.Save(&task1)
		<-repo.Save(&task2)
		<-repo.Save(&task3)
		<-repo.Save(&task4)
		<-repo.Save(&task5)

		// When
		result = <-repo.FindAll()
		fmt.Println(result)

		queryparams := make(map[string]string)
		queryparams["is_due"] = "false"
		queryparams["priority"] = "URGENT"
		queryparams["status"] = "CREATED"
		queryparams["domain"] = "CROP"
		queryparams["asset_id"] = assetID.String()

		foundOne = <-repo.FindTasksWithFilter(queryparams)
		fmt.Println(foundOne)

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, taskErr1)
	assert.Nil(t, taskErr2)
	assert.Nil(t, taskErr3)
	assert.Nil(t, taskErr4)
	assert.Nil(t, taskErr5)

	val, ok := foundOne.Result.([]domain.Task)
	assert.Equal(t, ok, true)
	assert.Equal(t, 2, len(val))
}
