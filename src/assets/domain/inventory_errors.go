package domain

const (
	InventoryMaterialInvalidPlantType = iota
)

// InventoryMaterialError is a custom error from Go built-in error
type InventoryMaterialError struct {
	Code int
}

func (e InventoryMaterialError) Error() string {
	switch e.Code {
	case InventoryMaterialInvalidPlantType:
		return "Invalid plant type"
	default:
		return "Unrecognized Inventory Material Error Code"
	}
}
