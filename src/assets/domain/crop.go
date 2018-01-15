package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Crop struct {
	UID          uuid.UUID
	BatchID      string
	InitialArea  Area
	CurrentAreas []Area
	Type         CropType
	Inventory    InventoryMaterial
	Container    CropContainer
	Notes        []CropNote
	CreatedDate  time.Time
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

type CropNote struct {
	Content     string    `json:"content"`
	CreatedDate time.Time `json:"created_date"`
}

func CreateCropBatch(area Area) (Crop, error) {
	if area.UID == (uuid.UUID{}) {
		return Crop{}, CropError{Code: CropErrorInvalidArea}
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return Crop{}, err
	}

	return Crop{
		UID:          uid,
		InitialArea:  area,
		CurrentAreas: []Area{area},
		CreatedDate:  time.Now(),
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

func (c *Crop) ChangeContainer(container CropContainer) error {
	err := validateCropContainer(container)
	if err != nil {
		return err
	}

	c.Container = container

	return nil
}

func (c *Crop) AddNewNote(content string) error {
	if content == "" {
		return CropError{Code: CropNoteErrorInvalidContent}
	}

	cropNote := CropNote{
		Content:     content,
		CreatedDate: time.Now(),
	}

	c.Notes = append(c.Notes, cropNote)

	return nil
}

func (c *Crop) RemoveNote(content string) error {
	if content == "" {
		return CropError{Code: CropNoteErrorInvalidContent}
	}

	for i, v := range c.Notes {
		if v.Content == content {
			copy(c.Notes[i:], c.Notes[i+1:])
			c.Notes[len(c.Notes)-1] = CropNote{}
			c.Notes = c.Notes[:len(c.Notes)-1]
		}
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
