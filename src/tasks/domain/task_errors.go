package domain

const (
	// Description Errors
	TaskErrorDescriptionEmptyCode = iota
	TaskErrorDescriptionAlphanumericOnlyCode

	// Date Errors
	TaskErrorDueDateEmptyCode
	TaskErrorDueDateInvalidCode

	// Priority Errors
	TaskErrorPriorityEmptyCode
	TaskErrorInvalidPriorityCode

	// Status Errors
	TaskErrorStatusEmptyCode
	TaskErrorInvalidStatusCode

	// Category Errors
	TaskErrorCategoryEmptyCode
	TaskErrorInvalidCategoryCode

	// Parent UID Errors
	TaskErrorAssetIDEmptyCode
	TaskErrorInvalidAssetIDCode

	// Task General Errors
	TaskErrorTaskNotFoundCode
)

// TaskError is a custom error from Go built-in error
type TaskError struct {
	Code int
}

func (e TaskError) Error() string {
	switch e.Code {
	case TaskErrorDescriptionEmptyCode:
		return "Task description is required."
	case TaskErrorDescriptionAlphanumericOnlyCode:
		return "Task description should be alphanumeric."
	case TaskErrorDueDateEmptyCode:
		return "Task due date is required."
	case TaskErrorDueDateInvalidCode:
		return "Task due date cannot be earlier than the current date."
	case TaskErrorPriorityEmptyCode:
		return "Task priority is required."
	case TaskErrorInvalidPriorityCode:
		return "Task priority is invalid."
	case TaskErrorStatusEmptyCode:
		return "Task status is required."
	case TaskErrorInvalidStatusCode:
		return "Task status is invalid."
	case TaskErrorCategoryEmptyCode:
		return "Task category is required."
	case TaskErrorInvalidCategoryCode:
		return "Task category is invalid."
	case TaskErrorAssetIDEmptyCode:
		return "Task must have a referenced asset."
	case TaskErrorInvalidAssetIDCode:
		return "Task asset reference is invalid."
	case TaskErrorTaskNotFoundCode:
		return "Task not found"
	default:
		return "Unrecognized Task Error Code"
	}
}
