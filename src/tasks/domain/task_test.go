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

	tasktype := "crop"

	var tests = []struct {
		description        string
		duedate            *time.Time
		priority           string
		tasktype           string
		assetid            uuid.UUID
		eexpectedTaskError error
	}{
		//emptyduedate
		{"MyDescription", nil, "urgent", tasktype, assetID, nil},
		//duedatepassed
		{"MyDescription", due_ptr_invalid, "urgent", tasktype, assetID, TaskError{TaskErrorDueDateInvalidCode}},
		//empty priority
		{"MyDescription", due_ptr, "", tasktype, assetID, TaskError{TaskErrorPriorityEmptyCode}},
		//invalidpriority
		{"MyDescription", due_ptr, "later", tasktype, assetID, TaskError{TaskErrorInvalidPriorityCode}},
		//empty type
		{"MyDescription", due_ptr, "urgent", "", assetID, TaskError{TaskErrorTypeEmptyCode}},
		//invalid type
		{"MyDescription", due_ptr, "urgent", "vegetable", assetID, TaskError{TaskErrorInvalidTypeCode}},
	}

	for _, test := range tests {
		taskServiceMock.On("FindCropByID", test.assetid).Return(ServiceResult{Result: query.TaskCropQueryResult{}})

		_, err := CreateTask(
			taskServiceMock, test.description, test.duedate, test.priority, test.tasktype, test.assetid.String())

		assert.Equal(t, test.eexpectedTaskError, err)
	}

	//empty assetid
	_, err := CreateTask(
		taskServiceMock, "MyDescription", due_ptr, "urgent", "crop", "")

	assert.Equal(t, TaskError{TaskErrorAssetIDEmptyCode}, err)

	//assetid doesn't exist
	taskServiceMock.On("FindCropByID", assetID_notexist).Return(ServiceResult{Result: query.TaskCropQueryResult{}, Error: TaskError{TaskErrorInvalidAssetIDCode}})

	_, err = CreateTask(
		taskServiceMock, "MyDescription", due_ptr, "urgent", "crop", assetID_notexist.String())

	assert.Equal(t, TaskError{TaskErrorInvalidAssetIDCode}, err)
}
