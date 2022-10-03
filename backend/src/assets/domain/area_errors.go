package domain

const (
	AreaErrorNameEmptyCode = iota
	AreaErrorNameNotEnoughCharacterCode
	AreaErrorNameExceedMaximunCharacterCode
	AreaErrorNameAlphanumericOnlyCode
	AreaErrorFarmNotFound
	AreaErrorReservoirNotFound

	AreaErrorSizeEmptyCode
	AreaErrorInvalidSizeUnitCode

	AreaErrorTypeEmptyCode
	AreaErrorInvalidAreaTypeCode
	AreaErrorCropAlreadyCreated

	AreaErrorLocationEmptyCode
	AreaErrorInvalidAreaLocationCode

	AreaNoteErrorInvalidContent
	AreaNoteErrorInvalidID
	AreaNoteErrorNotFound
)

// AreaError is a custom error from Go built-in error.
type AreaError struct {
	Code int
}

func (e AreaError) Error() string {
	switch e.Code {
	case AreaErrorNameEmptyCode:
		return "Area name is required."
	case AreaErrorNameNotEnoughCharacterCode:
		return "Not enough character on Area Name"
	case AreaErrorNameExceedMaximunCharacterCode:
		return "Area name cannot more than 100 characters"
	case AreaErrorNameAlphanumericOnlyCode:
		return "Area name should be alphanumeric, space, hypen, or underscore"
	case AreaErrorFarmNotFound:
		return "Farm not found"
	case AreaErrorReservoirNotFound:
		return "Reservoir not found"
	case AreaErrorSizeEmptyCode:
		return "Area size cannot be empty"
	case AreaErrorInvalidSizeUnitCode:
		return "Area size unit is invalid"
	case AreaErrorTypeEmptyCode:
		return "Area type cannot be empty"
	case AreaErrorInvalidAreaTypeCode:
		return "Area type is invalid"
	case AreaErrorCropAlreadyCreated:
		return "Area type cannot be changed because there is already filled with crops"
	case AreaErrorLocationEmptyCode:
		return "Area location cannot be empty"
	case AreaErrorInvalidAreaLocationCode:
		return "Area location is invalid"
	case AreaNoteErrorInvalidID:
		return "Invalid note id"
	case AreaNoteErrorInvalidContent:
		return "Invalid crop note content"
	case AreaNoteErrorNotFound:
		return "Area note not found"
	default:
		return "Unrecognized Area Error Code"
	}
}
