package domain

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MyMockedArea struct{
  mock.Mock
}

func TestCreateFarm(t *testing.T) {
	var tests = []struct {
		description         string
		duedate             string
		priority            string
		status             	string
		tasktype          	string
		assetid             string
		eexpectedTaskError	error
	}{
		//valid values
		//{"MyDescription", "2019-01-23T17:37:39.697328206+01:00", "urgent", "inprogress", "area", "46b054ab-a080-4c0b-ada9-ce920b585512", nil},
		//empty description
		{"", "2019-01-23T17:37:39.697328206+01:00", "urgent", "inprogress", "area", "46b054ab-a080-4c0b-ada9-ce920b585512", TaskError{TaskErrorDescriptionEmptyCode}},
		//invalid description
		{"!#@$MyDescription", "2019-01-23T17:37:39.697328206+01:00", "urgent", "inprogress", "area", "46b054ab-a080-4c0b-ada9-ce920b585512", TaskError{TaskErrorDescriptionAlphanumericOnlyCode}},
		//emptyduedate
		{"MyDescription", "", "urgent", "inprogress", "area", "46b054ab-a080-4c0b-ada9-ce920b585512", TaskError{TaskErrorDueDateEmptyCode}},
		//duedatepassed
		{"MyDescription", "2017-01-23T17:37:39.697328206+01:00", "urgent", "inprogress", "area", "46b054ab-a080-4c0b-ada9-ce920b585512", TaskError{TaskErrorDueDateInvalidCode}},
		//empty priority
		{"MyDescription", "2019-01-23T17:37:39.697328206+01:00", "", "inprogress", "area", "46b054ab-a080-4c0b-ada9-ce920b585512", TaskError{TaskErrorPriorityEmptyCode}},
		//invalidpriority
		{"MyDescription", "2019-01-23T17:37:39.697328206+01:00", "later", "inprogress", "area", "46b054ab-a080-4c0b-ada9-ce920b585512", TaskError{TaskErrorInvalidPriorityCode}},
		//empty status
		{"MyDescription", "2019-01-23T17:37:39.697328206+01:00", "urgent", "", "area", "46b054ab-a080-4c0b-ada9-ce920b585512", TaskError{TaskErrorStatusEmptyCode}},
		//invalid status
		{"MyDescription", "2019-01-23T17:37:39.697328206+01:00", "urgent", "done", "area", "46b054ab-a080-4c0b-ada9-ce920b585512", TaskError{TaskErrorInvalidStatusCode}},
		//empty type
		{"MyDescription", "2019-01-23T17:37:39.697328206+01:00", "urgent", "inprogress", "", "46b054ab-a080-4c0b-ada9-ce920b585512", TaskError{TaskErrorTypeEmptyCode}},
		//invalid type
		{"MyDescription", "2019-01-23T17:37:39.697328206+01:00", "urgent", "inprogress", "vegetable", "46b054ab-a080-4c0b-ada9-ce920b585512", TaskError{TaskErrorInvalidTypeCode}},
		//empty assetid
		{"MyDescription", "2019-01-23T17:37:39.697328206+01:00", "urgent", "inprogress", "area", "", TaskError{TaskErrorAssetIDEmptyCode}},
		//assetid doesn't exist
		{"MyDescription", "2019-01-23T17:37:39.697328206+01:00", "urgent", "inprogress", "area", "46b054ab-a080-4c0b-ada9-ce920b585512", TaskError{TaskErrorInvalidAssetIDCode}},
	}

	for _, test := range tests {
		task, err := CreateTask(test.description, test.duedate, test.priority, test.status, test.tasktype, test.assetid)

		assert.Equal(t, test.eexpectedTaskError, err)
	}
}