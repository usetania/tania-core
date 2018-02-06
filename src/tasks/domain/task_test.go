package domain

import (
	"github.com/Tanibox/tania-server/src/tasks/query"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type TaskServiceMock struct {
	mock.Mock
}

func (m TaskServiceMock) FindAreaByID(uid uuid.UUID) ServiceResult {
	args := m.Called(uid)
	return args.Get(0).(ServiceResult)
}
func (m TaskServiceMock) FindCropByID(uid uuid.UUID) ServiceResult {
	args := m.Called(uid)
	return args.Get(0).(ServiceResult)
}

func TestCreateTask(t *testing.T) {
	taskServiceMock := new(TaskServiceMock)

	assetID, _ := uuid.NewV4()
	assetID_notexist, _ := uuid.NewV4()

	due_date_invalid, _ := time.Parse(time.RFC3339, "2017-01-23T17:37:39.697328206+01:00")
	due_ptr_invalid := &due_date_invalid
	due_date, _ := time.Parse(time.RFC3339, "2019-01-23T17:37:39.697328206+01:00")
	due_ptr := &due_date

	taskcategory := "CROP"

	var tests = []struct {
		description        string
		duedate            *time.Time
		priority           string
		taskcategory       string
		assetid            *uuid.UUID
		eexpectedTaskError error
	}{
		//emptyduedate
		{"MyDescription", nil, "URGENT", taskcategory, &assetID, nil},
		//duedatepassed
		{"MyDescription", due_ptr_invalid, "NORMAL", taskcategory, &assetID, TaskError{TaskErrorDueDateInvalidCode}},
		//empty priority
		{"MyDescription", due_ptr, "", taskcategory, &assetID, TaskError{TaskErrorPriorityEmptyCode}},
		//invalidpriority
		{"MyDescription", due_ptr, "urgent", taskcategory, &assetID, TaskError{TaskErrorInvalidPriorityCode}},
		//empty category
		{"MyDescription", due_ptr, "URGENT", "", &assetID, TaskError{TaskErrorCategoryEmptyCode}},
		//invalid category
		{"MyDescription", due_ptr, "NORMAL", "VEGETABLE", &assetID, TaskError{TaskErrorInvalidCategoryCode}},
	}

	for _, test := range tests {
		taskServiceMock.On("FindCropByID", *test.assetid).Return(ServiceResult{Result: query.TaskCropQueryResult{}})

		_, err := CreateTask(
			taskServiceMock, test.description, test.duedate, test.priority, test.taskcategory, test.assetid)

		assert.Equal(t, test.eexpectedTaskError, err)
	}

	//nil assetid
	_, err := CreateTask(
		taskServiceMock, "MyDescription", due_ptr, "URGENT", taskcategory, nil)

	assert.Equal(t, nil, err)

	//assetid doesn't exist
	taskServiceMock.On("FindCropByID", assetID_notexist).Return(ServiceResult{Result: query.TaskCropQueryResult{}, Error: TaskError{TaskErrorInvalidAssetIDCode}})

	_, err = CreateTask(
		taskServiceMock, "MyDescription", due_ptr, "NORMAL", taskcategory, &assetID_notexist)

	assert.Equal(t, TaskError{TaskErrorInvalidAssetIDCode}, err)
}
