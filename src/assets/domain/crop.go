package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Crop struct {
	UID         uuid.UUID
	BatchID     string
	InitialArea Area
	CurrentArea Area
	Type        CropType
	Inventory   InventoryMaterial
	Container   CropContainer
	CreatedDate time.Time
}

// CropType defines type of crop
type CropType interface {
	Code() string
}

// Nursery implements CropType
type Nursery struct{}

func (s Nursery) Code() string { return "nursery" }

// Growing implements CropType
type Growing struct{}

func (g Growing) Code() string { return "growing" }

// CropContainer defines the container of a crop
type CropContainer struct {
	Quantity int
	Type     CropContainerType
}

// CropContainerType defines the type of a container
type CropContainerType interface {
	Code() string
}

// Tray implements CropContainerType
type Tray struct {
	Cell int
}

func (t Tray) Code() string { return "tray" }

// Pot implements CropContainerType
type Pot struct{}

func (p Pot) Code() string { return "pot" }

func CreateCropBatch(area Area) (Crop, error) {
	if area.UID == (uuid.UUID{}) {
		return Crop{}, CropError{Code: CropErrorInvalidArea}
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return Crop{}, err
	}

	return Crop{
		UID:         uid,
		InitialArea: area,
		CurrentArea: area,
		CreatedDate: time.Now(),
	}, nil
}

func (c *Crop) ChangeCropType(cropType CropType) error {
	err := validateCropType(cropType)
	if err != nil {
		return err
	}

	c.Type = cropType

	return nil
}

func (c *Crop) ChangeInventory(inventory InventoryMaterial) error {
	err := validateInventory(inventory)
	if err != nil {
		return err
	}

	batchID, err := generateBatchID(inventory)

	c.Inventory = inventory
	c.BatchID = batchID

	return nil
}

func (c *Crop) ChangeContainer(container CropContainer) error {
	err := validateCropContainer(container)
	if err != nil {
		return err
	}

	return nil
}

func generateBatchID(inventory InventoryMaterial) (string, error) {
	batchID := "DUMMY-BATCH-ID"

	return batchID, nil
}

func validateInventory(inventory InventoryMaterial) error {
	err := validatePlantType(inventory.PlantType)
	if err != nil {
		return err
	}

	if inventory.Variety == "" {
		return InventoryMaterialError{Code: InventoryMaterialInvalidVariety}
	}

	return nil
}

func validateCropType(cropType CropType) error {
	switch cropType.(type) {
	case Nursery:
	case Growing:
	default:
		return CropError{Code: CropErrorInvalidCropType}
	}

	return nil
}

func validateCropContainer(container CropContainer) error {
	switch container.Type.(type) {
	case Tray:
	case Pot:
	default:
		return CropError{Code: CropContainerErrorInvalidType}
	}

	return nil
}
