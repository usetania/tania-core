package domain

const (
	CropErrorInvalidArea = iota
	CropErrorInvalidCropType

	CropErrorInvalidBatchID
	CropErrorBatchIDAlreadyCreated

	CropContainerErrorInvalidType
)

// CropError is a custom error from Go built-in error
type CropError struct {
	Code int
}

func (e CropError) Error() string {
	switch e.Code {
	case CropErrorInvalidArea:
		return "Invalid area"
	case CropErrorInvalidCropType:
		return "Invalid crop type"
	case CropErrorInvalidBatchID:
		return "Invalid crop batch ID"
	case CropErrorBatchIDAlreadyCreated:
		return "Crop batch ID already created"
	case CropContainerErrorInvalidType:
		return "Invalid crop container type"
	default:
		return "Unrecognized Crop Error Code"
	}
}
