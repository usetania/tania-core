package location

// LocationError is a custom error from Go built-in error
type LocationError struct {
	code int
}

const (
	LocationErrorInvalidCountryCode = iota // 0
	LocationErrorInvalidCityCode           // 1
)

func (e LocationError) Error() string {
	switch e.code {
	case LocationErrorInvalidCountryCode:
		return "Country code value is invalid."
	case LocationErrorInvalidCityCode:
		return "City code value is invalid."
	default:
		return "Unrecognized location error code"
	}
}
