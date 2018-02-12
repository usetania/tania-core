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
func (m TaskServiceMock) FindReservoirByID(uid uuid.UUID) domain.ServiceResult {
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

	assetID1, _ := uuid.NewV4()
	assetID2, _ := uuid.NewV4()
	assetID3, _ := uuid.NewV4()
	assetID4, _ := uuid.NewV4()

	taskServiceMock.On("FindCropByID", assetID1).Return(domain.ServiceResult{Result: query.TaskCropQueryResult{UID: inventoryID}})
	task1, taskErr1 := domain.CreateTask(
		taskServiceMock, title, description, due_ptr, priority, taskdomain_crop,
		category, &assetID1)

	taskServiceMock.On("FindCropByID", assetID2).Return(domain.ServiceResult{Result: query.TaskCropQueryResult{UID: inventoryID}})
	task2, taskErr2 := domain.CreateTask(
		taskServiceMock, title, description, due_ptr, priority, taskdomain_crop,
		category, &assetID2)
	task2.SetTaskAsDue()

	taskServiceMock.On("FindAreaByID", assetID3).Return(domain.ServiceResult{Result: query.TaskCropQueryResult{UID: inventoryID}})
	task3, taskErr3 := domain.CreateTask(
		taskServiceMock, title, description, due_ptr, priority, taskdomain_area,
		category, &assetID3)

	taskServiceMock.On("FindAreaByID", assetID4).Return(domain.ServiceResult{Result: query.TaskCropQueryResult{UID: inventoryID}})
	task4, taskErr4 := domain.CreateTask(
		taskServiceMock, title, description, due_ptr, "NORMAL", taskdomain_area,
		category, &assetID4)
	task4.SetTaskAsDue()

	taskServiceMock.On("FindAreaByID", assetID3).Return(domain.ServiceResult{Result: query.TaskCropQueryResult{UID: inventoryID}})
	task5, taskErr5 := domain.CreateTask(
		taskServiceMock, title, description, due_ptr, priority, taskdomain_area,
		category, &assetID3)

	var result, tc1_result, tc2_result, tc3_result, tc4_result, tc5_result RepositoryResult
	go func() {
		// Given
		<-repo.Save(&task1)
		<-repo.Save(&task2)
		<-repo.Save(&task3)
		<-repo.Save(&task4)
		<-repo.Save(&task5)

		// When
		result = <-repo.FindAll()

		// Test Case 1

		queryparams := make(map[string]string)
		queryparams["is_due"] = "false"
		queryparams["priority"] = "URGENT"
		queryparams["status"] = "CREATED"
		queryparams["domain"] = "CROP"
		queryparams["asset_id"] = assetID1.String()

		tc1_result = <-repo.FindTasksWithFilter(queryparams)
		fmt.Println(tc1_result)

		// Test Case 2 Get all "Due" tasks

		queryparams2 := make(map[string]string)
		queryparams2["is_due"] = "true"

		tc2_result = <-repo.FindTasksWithFilter(queryparams2)
		fmt.Println(tc2_result)

		// Test Case 3 Get all tasks under a specific asset

		queryparams3 := make(map[string]string)
		queryparams3["domain"] = "AREA"
		queryparams3["asset_id"] = assetID3.String()

		tc3_result = <-repo.FindTasksWithFilter(queryparams3)
		fmt.Println(tc3_result)

		// Test Case 4 Get all "urgent" tasks

		queryparams4 := make(map[string]string)
		queryparams4["priority"] = "URGENT"

		tc4_result = <-repo.FindTasksWithFilter(queryparams4)
		fmt.Println(tc4_result)

		// Test Case 5 Get all "CROP" tasks

		queryparams5 := make(map[string]string)
		queryparams5["domain"] = "CROP"

		tc5_result = <-repo.FindTasksWithFilter(queryparams5)
		fmt.Println(tc5_result)
		done <- true

	}()

	// Then
	<-done
	assert.Nil(t, taskErr1)
	assert.Nil(t, taskErr2)
	assert.Nil(t, taskErr3)
	assert.Nil(t, taskErr4)
	assert.Nil(t, taskErr5)

	val1, ok := tc1_result.Result.([]domain.Task)
	assert.Equal(t, ok, true)
	assert.Equal(t, 1, len(val1))

	val2, ok := tc2_result.Result.([]domain.Task)
	assert.Equal(t, ok, true)
	assert.Equal(t, 2, len(val2))

	val3, ok := tc3_result.Result.([]domain.Task)
	assert.Equal(t, ok, true)
	assert.Equal(t, 2, len(val3))

	val4, ok := tc4_result.Result.([]domain.Task)
	assert.Equal(t, ok, true)
	assert.Equal(t, 4, len(val4))

	val5, ok := tc5_result.Result.([]domain.Task)
	assert.Equal(t, ok, true)
	assert.Equal(t, 2, len(val5))
}
