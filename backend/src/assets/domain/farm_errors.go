package domain

// FarmError is a custom error from Go built-in error.
type FarmError struct {
	Code int
}

const (
	FarmErrorInvalidFarmTypeCode = iota

	FarmErrorReservoirAlreadyAdded
	FarmErrorReservoirNotFound

	FarmErrorAreaAlreadyAdded
	FarmErrorAreaNotFound

	FarmErrorNameEmptyCode
	FarmErrorNameNotEnoughCharacterCode
	FarmErrorNameExceedMaximunCharacterCode
	FarmErrorNameAlphanumericOnlyCode

	FarmErrorInvalidLatitudeValueCode
	FarmErrorInvalidLongitudeValueCode
	FarmErrorInvalidCountry
	FarmErrorInvalidCity
)

func (e FarmError) Error() string {
	switch e.Code {
	case FarmErrorInvalidFarmTypeCode:
		return "Farm type code value is invalid."
	case FarmErrorReservoirAlreadyAdded:
		return "Reservoir is already added."
	case FarmErrorReservoirNotFound:
		return "Farm reservoir not found."
	case FarmErrorAreaAlreadyAdded:
		return "Area is already added."
	case FarmErrorAreaNotFound:
		return "Farm area not found."
	case FarmErrorNameEmptyCode:
		return "Farm name is required."
	case FarmErrorNameNotEnoughCharacterCode:
		return "Not enough character on farm name"
	case FarmErrorNameExceedMaximunCharacterCode:
		return "Farm name cannot more than 100 characters"
	case FarmErrorNameAlphanumericOnlyCode:
		return "Farm name should be alphanumeric, space, hypen, or underscore"
	case FarmErrorInvalidLatitudeValueCode:
		return "Latitude value is invalid"
	case FarmErrorInvalidLongitudeValueCode:
		return "Longitude value is invalid"
	case FarmErrorInvalidCountry:
		return "Invalid country"
	case FarmErrorInvalidCity:
		return "Invalid city"
	default:
		return "Unrecognized location error code"
	}
}
