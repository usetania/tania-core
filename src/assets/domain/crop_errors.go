package domain

const (
	CropErrorInvalidArea = iota
	CropErrorInvalidCropType
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
	default:
		return "Unrecognized Crop Error Code"
	}
}
