package domain_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	. "github.com/usetania/tania-core/src/tasks/domain"
	"github.com/usetania/tania-core/src/tasks/query"
)

type TaskServiceMock struct {
	mock.Mock
}

func (m *TaskServiceMock) FindAreaByID(uid uuid.UUID) ServiceResult {
	args := m.Called(uid)

	return args.Get(0).(ServiceResult)
}

func (m *TaskServiceMock) FindCropByID(uid uuid.UUID) ServiceResult {
	args := m.Called(uid)

	return args.Get(0).(ServiceResult)
}

func (m *TaskServiceMock) FindMaterialByID(uid uuid.UUID) ServiceResult {
	args := m.Called(uid)

	return args.Get(0).(ServiceResult)
}

func (m *TaskServiceMock) FindReservoirByID(uid uuid.UUID) ServiceResult {
	args := m.Called(uid)

	return args.Get(0).(ServiceResult)
}

func TestCreateTask(t *testing.T) {
	t.Parallel()

	taskServiceMock := new(TaskServiceMock)

	assetID, _ := uuid.NewV4()
	assetIDNotExist, _ := uuid.NewV4()

	dueDateInvalid := time.Now().Add(-1 * time.Hour)
	duePtrInvalid := &dueDateInvalid
	dueDate := time.Now().Add(1 * time.Hour)
	duePtr := &dueDate

	tasktitle := "My Task"
	taskdescription := "My Description"
	taskcategory := "SANITATION"

	cropID, _ := uuid.NewV4()
	batchID := "bro-sup-gre-3-4mar"

	taskServiceMock.On("FindCropByID", cropID).Return(ServiceResult{
		Result: query.TaskCropResult{
			UID:     cropID,
			BatchID: batchID,
		},
	})

	materialID, _ := uuid.NewV4()
	materialName := "A Good One Fertilizer For All"

	taskServiceMock.On("FindMaterialByID", materialID).Return(ServiceResult{
		Result: query.TaskMaterialResult{
			UID:              materialID,
			Name:             materialName,
			TypeCode:         "AGROCHEMICAL",
			DetailedTypeCode: "FERTILIZER",
		},
	})

	areaID, _ := uuid.NewV4()
	areaName := "MY AREA SEEDING"

	taskServiceMock.On("FindAreaByID", areaID).Return(ServiceResult{
		Result: query.TaskAreaResult{
			UID:  areaID,
			Name: areaName,
		},
	})

	taskdomain, _ := CreateTaskDomainCrop(taskServiceMock, taskcategory, &materialID, &areaID)

	tests := []struct {
		title              string
		description        string
		duedate            *time.Time
		priority           string
		domain             TaskDomain
		category           string
		assetid            *uuid.UUID
		eexpectedTaskError error
	}{
		// empty title
		{"", taskdescription, duePtr, "NORMAL", taskdomain, taskcategory, &assetID, TaskError{TaskErrorTitleEmptyCode}},
		// emptyduedate
		{tasktitle, taskdescription, nil, "URGENT", taskdomain, taskcategory, &assetID, nil},
		// duedatepassed
		{tasktitle, taskdescription, duePtrInvalid, "NORMAL", taskdomain, taskcategory, &assetID, TaskError{TaskErrorDueDateInvalidCode}}, //nolint:lll
		// empty priority
		{tasktitle, taskdescription, duePtr, "", taskdomain, taskcategory, &assetID, TaskError{TaskErrorPriorityEmptyCode}},
		// invalidpriority
		{tasktitle, taskdescription, duePtr, "urgent", taskdomain, taskcategory, &assetID, TaskError{TaskErrorInvalidPriorityCode}}, //nolint:lll
		// empty category
		{tasktitle, taskdescription, duePtr, "URGENT", taskdomain, "", &assetID, TaskError{TaskErrorCategoryEmptyCode}},
		// invalid category
		{tasktitle, taskdescription, duePtr, "NORMAL", taskdomain, "VEGETABLE", &assetID, TaskError{TaskErrorInvalidCategoryCode}}, //nolint:lll
	}

	for _, test := range tests {
		taskServiceMock.On("FindCropByID", *test.assetid).Return(ServiceResult{Result: query.TaskCropResult{}})
		taskServiceMock.On("FindMaterialByID", materialID).Return(ServiceResult{
			Result: query.TaskMaterialResult{
				UID:              materialID,
				Name:             materialName,
				TypeCode:         "AGROCHEMICAL",
				DetailedTypeCode: "FERTILIZER",
			},
		})

		_, err := CreateTask(
			taskServiceMock, test.title, test.description, test.priority, test.category, test.duedate, test.domain, test.assetid)

		assert.Equal(t, test.eexpectedTaskError, err)
	}

	// nil assetid
	taskServiceMock.On("FindMaterialByID", materialID).Return(ServiceResult{
		Result: query.TaskMaterialResult{
			UID:              materialID,
			Name:             materialName,
			TypeCode:         "AGROCHEMICAL",
			DetailedTypeCode: "FERTILIZER",
		},
	})

	_, err := CreateTask(
		taskServiceMock, tasktitle, taskdescription, "URGENT", taskcategory, duePtr, taskdomain, nil)

	assert.Equal(t, nil, err)

	// assetid doesn't exist
	taskServiceMock.On("FindCropByID", assetIDNotExist).Return(ServiceResult{
		Result: query.TaskCropResult{},
		Error:  TaskError{TaskErrorInvalidAssetIDCode},
	})
	taskServiceMock.On("FindMaterialByID", materialID).Return(ServiceResult{
		Result: query.TaskMaterialResult{
			UID:              materialID,
			Name:             materialName,
			TypeCode:         "AGROCHEMICAL",
			DetailedTypeCode: "FERTILIZER",
		},
	})

	_, err = CreateTask(
		taskServiceMock, tasktitle, taskdescription, "NORMAL", taskcategory, duePtr, taskdomain, &assetIDNotExist)

	assert.Equal(t, TaskError{TaskErrorInvalidAssetIDCode}, err)
}
