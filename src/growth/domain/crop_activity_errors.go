package domain

const (
	CropActivityMoveErrorInvalidArea = iota
	CropActivityMoveErrorInvalidQuantity
)

// CropActivityError is a custom error from Go built-in error
type CropActivityError struct {
	Code int
}

func (e CropActivityError) Error() string {
	switch e.Code {
	case CropActivityMoveErrorAreaNotValid:
		return "Invalid Area. Cannot move the crop from selected area to destined area"
	case CropActivityMoveErrorInvalidQuantity:
		return "Invalid quantity. Make sure your quantity is enough to be moved"
	default:
		return "Unrecognized crop activity error code"
	}
}
