package domain

const (
	CropErrorInvalidArea = iota
	CropErrorInvalidCropType
	CropErrorInvalidCropStatus

	CropMoveToAreaErrorInvalidQuantity
	CropMoveToAreaErrorInvalidArea
	CropMoveToAreaErrorInvalidExistingSourceArea
	CropMoveToAreaErrorCannotBeSame

	CropErrorInvalidBatchID
	CropErrorBatchIDAlreadyCreated

	CropAreaErrorInvalidSourceArea
	CropAreaErrorSourceAreaNotFound
	CropAreaErrorInvalidDestinationArea
	CropAreaErrorDestinationAreaNotFound

	CropContainerErrorInvalidType
	CropContainerErrorInvalidQuantity
	CropContainerErrorInvalidTrayCell

	CropInventoryErrorInvalidInventory
	CropInventoryErrorNotFound

	CropNoteErrorInvalidContent
	CropNoteErrorNotFound
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
	case CropErrorInvalidCropStatus:
		return "Invalid crop status"
	case CropErrorInvalidBatchID:
		return "Invalid crop batch ID"
	case CropErrorBatchIDAlreadyCreated:
		return "Crop batch ID already created"

	case CropMoveToAreaErrorInvalidQuantity:
		return "Invalid quantity. Make sure your quantity is not zero and enough to be moved"
	case CropMoveToAreaErrorInvalidArea:
		return "Invalid move crop to area. Crop can only be moved from Seeding to Growing, Seeding to Seeding or Growing to Growing"
	case CropMoveToAreaErrorInvalidExistingSourceArea:
		return "Invalid existing source area"
	case CropMoveToAreaErrorCannotBeSame:
		return "Invalid move crop to area. Area source and destination cannot be same"

	case CropAreaErrorInvalidSourceArea:
		return "Crop source area is invalid"
	case CropAreaErrorSourceAreaNotFound:
		return "Crop source area not found"
	case CropAreaErrorInvalidDestinationArea:
		return "Crop destination area is invalid"
	case CropAreaErrorDestinationAreaNotFound:
		return "Crop destination not found"

	case CropContainerErrorInvalidType:
		return "Invalid crop container type"
	case CropContainerErrorInvalidQuantity:
		return "Invalid crop container quantity"
	case CropContainerErrorInvalidTrayCell:
		return "Invalid crop container tray cell"

	case CropInventoryErrorInvalidInventory:
		return "Invalid crop inventory"
	case CropInventoryErrorNotFound:
		return "Crop inventory not found"

	case CropNoteErrorInvalidContent:
		return "Invalid crop note content"
	case CropNoteErrorNotFound:
		return "Crop note not found"
	default:
		return "Unrecognized Crop Error Code"
	}
}
