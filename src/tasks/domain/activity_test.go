package domain

import (
	"testing"

	"github.com/Tanibox/tania-server/src/tasks/query"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	source_areaUID, _ := uuid.NewV4()
	source_areaUID_notexist, _ := uuid.NewV4()
	dest_areaUID, _ := uuid.NewV4()
	dest_areaUID_notexist, _ := uuid.NewV4()
	//cropUID, _ := uuid.NewV4()
	source_areaServiceResult := ServiceResult{
		Result: query.TaskAreaQueryResult{UID: source_areaUID},
	}
	dest_areaServiceResult := ServiceResult{
		Result: query.TaskAreaQueryResult{UID: dest_areaUID},
	}
	//cropServiceResult := ServiceResult{
	//	Result: query.TaskCropQueryResult{UID: cropUID},
	//}

	taskServiceMock.On("FindAreaByID", source_areaUID).Return(source_areaServiceResult)
	taskServiceMock.On("FindAreaByID", dest_areaUID).Return(dest_areaServiceResult)
	//taskServiceMock.On("FindCropByID", cropUID).Return(cropServiceResult)

	// Harvest Activity
	var tests_harvestactivity = []struct {
		source_area_id     string
		quantity           string
		eexpectedTaskError error
	}{
		// empty source_area_id
		{"", "67", TaskError{TaskErrorActivitySourceInvalidCode}},
		// invalid source_area_id
		{source_areaUID_notexist.String(), "67", TaskError{TaskErrorActivitySourceInvalidCode}},
		// invalid quantity
		{source_areaUID.String(), "nan", TaskError{TaskErrorActivityQuantityInvalidCode}},
	}

	for _, test := range tests_harvestactivity {

		_, err := CreateHarvestActivity(taskServiceMock, test.source_area_id, test.quantity)

		//taskServiceMock.AssertExpectations(t)

		assert.Equal(t, test.eexpectedTaskError, err)
	}

	// Dump Activity
	var tests_dumpactivity = []struct {
		source_area_id     string
		quantity           string
		eexpectedTaskError error
	}{
		// empty source_area_id
		{"", "67", TaskError{TaskErrorActivitySourceInvalidCode}},
		// invalid source_area_id
		{source_areaUID_notexist.String(), "67", TaskError{TaskErrorActivitySourceInvalidCode}},
		// invalid quantity
		{source_areaUID.String(), "nan", TaskError{TaskErrorActivityQuantityInvalidCode}},
	}

	for _, test := range tests_dumpactivity {

		_, err := CreateDumpActivity(taskServiceMock, test.source_area_id, test.quantity)

		//taskServiceMock.AssertExpectations(t)

		assert.Equal(t, test.eexpectedTaskError, err)
	}

	// Move To Area Activity
	var tests_movetoareaactivity = []struct {
		source_area_id     string
		dest_area_id       string
		quantity           string
		eexpectedTaskError error
	}{
		// empty source_area_id
		{"", dest_areaUID.String(), "67", TaskError{TaskErrorActivitySourceInvalidCode}},
		// invalid source_area_id
		{source_areaUID_notexist.String(), dest_areaUID.String(), "67", TaskError{TaskErrorActivitySourceInvalidCode}},
		// empty dest_area_id
		{source_areaUID.String(), "", "67", TaskError{TaskErrorActivityDestinationInvalidCode}},
		// invalid dest_area_id
		{source_areaUID.String(), dest_areaUID_notexist.String(), "67", TaskError{TaskErrorActivityDestinationInvalidCode}},
		// invalid quantity
		{source_areaUID.String(), dest_areaUID.String(), "nan", TaskError{TaskErrorActivityQuantityInvalidCode}},
	}

	for _, test := range tests_movetoareaactivity {

		_, err := CreateMoveToAreaActivity(taskServiceMock, test.source_area_id, test.dest_area_id, test.quantity)

		//taskServiceMock.AssertExpectations(t)

		assert.Equal(t, test.eexpectedTaskError, err)
	}
}
