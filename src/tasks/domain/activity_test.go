package domain

import (
	"github.com/Tanibox/tania-server/src/tasks/query"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
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

func TestCreateActivity(t *testing.T) {
	taskServiceMock := new(TaskServiceMock)

	source_area_id, _ := uuid.NewV4()
	source_area_id_not_exist, _ := uuid.NewV4()
	dest_area_id, _ := uuid.NewV4()
	dest_area_id_not_exist, _ := uuid.NewV4()

	// Harvest Activity
	var tests_harvestactivity = []struct {
		source_area_id     uuid.UUID
		quantity           string
		expectedMockError  error
		eexpectedTaskError error
	}{
		// empty source_area_id
		{uuid.UUID{}, "67", TaskError{TaskErrorActivitySourceInvalidCode}, nil},
		// invalid source_area_id
		{source_area_id_not_exist, "67", TaskError{TaskErrorActivitySourceInvalidCode}, nil},
		// invalid quantity
		{source_area_id, "nan", nil, TaskError{TaskErrorActivityQuantityInvalidCode}},
	}

	for _, test := range tests_harvestactivity {
		source_areaServiceResult := ServiceResult{
			Result: query.TaskAreaQueryResult{UID: test.source_area_id},
			Error:  test.expectedMockError,
		}
		taskServiceMock.On("FindAreaByID", test.source_area_id).Return(source_areaServiceResult)
		_, err := CreateHarvestActivity(taskServiceMock, test.source_area_id.String(), test.quantity)

		taskServiceMock.AssertExpectations(t)

		if test.expectedMockError != nil {
			assert.Equal(t, test.expectedMockError, err)
		} else {
			assert.Equal(t, test.eexpectedTaskError, err)
		}
	}

	// Dump Activity
	var tests_dumpactivity = []struct {
		source_area_id     uuid.UUID
		quantity           string
		query_result       query.TaskAreaQueryResult
		eexpectedTaskError error
	}{
		// empty source_area_id
		{uuid.UUID{}, "67", query.TaskAreaQueryResult{}, TaskError{TaskErrorActivitySourceInvalidCode}},
		// invalid source_area_id
		{source_area_id_not_exist, "67", query.TaskAreaQueryResult{}, TaskError{TaskErrorActivitySourceInvalidCode}},
		// invalid quantity
		{source_area_id, "nan", query.TaskAreaQueryResult{UID: source_area_id}, TaskError{TaskErrorActivityQuantityInvalidCode}},
	}

	for _, test := range tests_dumpactivity {
		taskServiceMock.On("FindAreaByID", test.source_area_id).Return(ServiceResult{Result: test.query_result})
		_, err := CreateDumpActivity(taskServiceMock, test.source_area_id.String(), test.quantity)

		assert.Equal(t, test.eexpectedTaskError, err)
	}

	// Move To Area Activity
	var tests_movetoareaactivity = []struct {
		source_area_id     uuid.UUID
		dest_area_id       uuid.UUID
		quantity           string
		query_result       query.TaskAreaQueryResult
		eexpectedTaskError error
	}{
		// empty source_area_id
		{uuid.UUID{}, dest_area_id, "67", query.TaskAreaQueryResult{}, TaskError{TaskErrorActivitySourceInvalidCode}},
		// invalid source_area_id
		{source_area_id_not_exist, dest_area_id, "67", query.TaskAreaQueryResult{}, TaskError{TaskErrorActivitySourceInvalidCode}},
	}

	for _, test := range tests_movetoareaactivity {

		taskServiceMock.On("FindAreaByID", test.source_area_id).Return(ServiceResult{Result: test.query_result})
		_, err := CreateMoveToAreaActivity(taskServiceMock, test.source_area_id.String(), test.dest_area_id.String(), test.quantity)

		assert.Equal(t, test.eexpectedTaskError, err)
	}

	var tests_movetoareaactivity2 = []struct {
		source_area_id     uuid.UUID
		dest_area_id       uuid.UUID
		quantity           string
		query_result       query.TaskAreaQueryResult
		eexpectedTaskError error
	}{
		// empty dest_area_id
		{source_area_id, uuid.UUID{}, "67", query.TaskAreaQueryResult{}, TaskError{TaskErrorActivityDestinationInvalidCode}},
		// invalid dest_area_id
		{source_area_id, dest_area_id_not_exist, "67", query.TaskAreaQueryResult{}, TaskError{TaskErrorActivityDestinationInvalidCode}},
		// invalid quantity
		//{source_area_id, dest_area_id, "nan", query.TaskAreaQueryResult{}, TaskError{TaskErrorActivityQuantityInvalidCode}},
	}

	for _, test := range tests_movetoareaactivity2 {

		taskServiceMock.On("FindAreaByID", test.source_area_id).Return(ServiceResult{Result: query.TaskAreaQueryResult{UID: test.source_area_id}})
		taskServiceMock.On("FindAreaByID", test.dest_area_id).Return(ServiceResult{Result: test.query_result})
		_, err := CreateMoveToAreaActivity(taskServiceMock, test.source_area_id.String(), test.dest_area_id.String(), test.quantity)

		assert.Equal(t, test.eexpectedTaskError, err)
	}
}
