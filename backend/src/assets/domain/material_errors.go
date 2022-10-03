package domain

const (
	MaterialErrorInvalidMaterialType = iota
)

// MaterialError is a custom error from Go built-in error.
type MaterialError struct {
	Code int
}

func (e MaterialError) Error() string {
	switch e.Code {
	case MaterialErrorInvalidMaterialType:
		return "Invalid material type"
	default:
		return "Unrecognized Material Error Code"
	}
}
