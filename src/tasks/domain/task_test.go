package domain

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateTask(t *testing.T) {

	assetID, _ := uuid.NewV4()

	due_date_invalid, _ := time.Parse(time.RFC3339, "2017-01-23T17:37:39.697328206+01:00")
	due_ptr_invalid := &due_date_invalid
	due_date, _ := time.Parse(time.RFC3339, "2019-01-23T17:37:39.697328206+01:00")
	due_ptr := &due_date

	tasktype := "crop"

	act, _ := CreatePruneActivity()

	var tests = []struct {
		description        string
		duedate            *time.Time
		priority           string
		tasktype           string
		assetid            string
		eexpectedTaskError error
	}{
		//empty description
		{"", due_ptr, "urgent", tasktype, assetID.String(), TaskError{TaskErrorDescriptionEmptyCode}},
		//emptyduedate
		{"MyDescription", nil, "urgent", tasktype, assetID.String(), TaskError{TaskErrorDueDateEmptyCode}},
		//duedatepassed
		{"MyDescription", due_ptr_invalid, "urgent", tasktype, assetID.String(), TaskError{TaskErrorDueDateInvalidCode}},
		//empty priority
		{"MyDescription", due_ptr, "", tasktype, assetID.String(), TaskError{TaskErrorPriorityEmptyCode}},
		//invalidpriority
		{"MyDescription", due_ptr, "later", tasktype, assetID.String(), TaskError{TaskErrorInvalidPriorityCode}},
		//empty type
		{"MyDescription", due_ptr, "urgent", "", assetID.String(), TaskError{TaskErrorTypeEmptyCode}},
		//invalid type
		{"MyDescription", due_ptr, "urgent", "vegetable", assetID.String(), TaskError{TaskErrorInvalidTypeCode}},
		//empty assetid
		{"MyDescription", due_ptr, "urgent", tasktype, "", TaskError{TaskErrorAssetIDEmptyCode}},
		//assetid doesn't exist
		{"MyDescription", due_ptr, "urgent", tasktype, assetID.String(), TaskError{TaskErrorInvalidAssetIDCode}},
	}

	for _, test := range tests {

		_, err := CreateTask(
			nil, test.description, test.duedate, test.priority, test.tasktype, test.assetid, act)

		assert.Equal(t, test.eexpectedTaskError, err)
	}
}
