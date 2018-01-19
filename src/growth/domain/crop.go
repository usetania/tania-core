package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Crop struct {
	UID          uuid.UUID
	BatchID      string
	Status       CropStatus
	Type         CropType
	Container    CropContainer
	InventoryUID uuid.UUID
	CreatedDate  time.Time

	// Fields to maintain crop's movement
	InitialArea      InitialArea
	MovedArea        []MovedArea
	HarvestedStorage []HarvestedStorage
	Trash            []Trash

	// Fields to maintain care crop
	LastFertilized time.Time
	LastPruned     time.Time
	LastPesticided time.Time

	// Notes
	Notes map[uuid.UUID]CropNote
}

// We just recorded two crop status
// because other activities are recurring and parallel.
// So we can't have crop batch with a status, for example `HARVESTED`,
// because not all the crop has harvested
const (
	CropActive = "ACTIVE"
	CropEnd    = "END"
)

type CropStatus struct {
	Code  string
	Label string
}

func CropStatuses() []CropStatus {
	return []CropStatus{
		{Code: CropActive, Label: "Active"},
		{Code: CropEnd, Label: "End"},
	}
}

func GetCropStatus(code string) CropStatus {
	for _, v := range CropStatuses() {
		if code == v.Code {
			return v
		}
	}

	return CropStatus{}
}

const (
	CropTypeSeeding = "SEEDING"
	CropTypeGrowing = "GROWING"
)

type CropType struct {
	Code  string
	Label string
}

func CropTypes() []CropType {
	return []CropType{
		{Code: CropTypeSeeding, Label: "Seeding"},
		{Code: CropTypeGrowing, Label: "Growing"},
	}
}

func GetCropType(code string) CropType {
	for _, v := range CropTypes() {
		if code == v.Code {
			return v
		}
	}

	return CropType{}
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

type CropContainer struct {
	Quantity int
	Type     CropContainerType
}

type InitialArea struct {
	AreaUID         uuid.UUID
	InitialQuantity int
	CurrentQuantity int
}

type MovedArea struct {
	AreaUID         uuid.UUID
	SourceAreaUID   uuid.UUID
	InitialQuantity int
	CurrentQuantity int
	Date            time.Time
}

type HarvestedStorage struct {
	Quantity      int
	SourceAreaUID uuid.UUID
	Date          time.Time
}

type Trash struct {
	Quantity      int
	SourceAreaUID uuid.UUID
	Date          time.Time
}

type CropNote struct {
	UID         uuid.UUID
	Content     string
	CreatedDate time.Time
}

func CreateCropBatch() (Crop, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return Crop{}, err
	}

	return Crop{
		UID:         uid,
		CreatedDate: time.Now(),
	}, nil
	return Crop{}, nil
}

func (c *Crop) Fertilize() error {
	c.LastFertilized = time.Now()

	return nil
}

func (c *Crop) Prune() error {
	c.LastPruned = time.Now()

	return nil
}

func (c *Crop) Pesticide() error {
	c.LastPesticided = time.Now()

	return nil
}

func (c *Crop) ChangeCropType(cropType string) error {
	ct := GetCropType(cropType)
	if ct == (CropType{}) {
		return CropError{Code: CropErrorInvalidCropType}
	}

	c.Type = ct

	return nil
}

func (c *Crop) ChangeCropStatus(cropStatus string) error {
	cs := GetCropStatus(cropStatus)
	if cs == (CropStatus{}) {
		return CropError{Code: CropErrorInvalidCropStatus}
	}

	c.Status = cs

	return nil
}

func (c *Crop) ChangeContainer(quantity int, cell int, containerType CropContainerType) error {
	if quantity <= 0 {
		return CropError{Code: CropContainerErrorInvalidQuantity}
	}

	var t CropContainerType
	switch containerType.(type) {
	case Tray:
		if cell <= 0 {
			return CropError{Code: CropContainerErrorInvalidTrayCell}
		}

		t = Tray{Cell: cell}
	case Pot:
		t = Pot{}
	default:
		return CropError{Code: CropContainerErrorInvalidType}
	}

	c.Container = CropContainer{
		Quantity: quantity,
		Type:     t,
	}

	return nil
}

// ChangeInventory needs to be validated with external (ie: query to Inventory in Assets domain).
// BatchID will only be generated after inventory is assigned.
func (c *Crop) ChangeInventory(inventoryUID uuid.UUID) error {
	// Generate batchID from the inventory here

	return nil
}

func (c *Crop) AddNewNote(content string) error {
	if content == "" {
		return CropError{Code: CropNoteErrorInvalidContent}
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	cropNote := CropNote{
		UID:         uid,
		Content:     content,
		CreatedDate: time.Now(),
	}

	if len(c.Notes) == 0 {
		c.Notes = make(map[uuid.UUID]CropNote)
	}

	c.Notes[uid] = cropNote

	return nil
}

func (c *Crop) RemoveNote(uid uuid.UUID) error {
	if uid == (uuid.UUID{}) {
		return CropError{Code: CropNoteErrorNotFound}
	}

	found := false
	for _, v := range c.Notes {
		if v.UID == uid {
			delete(c.Notes, uid)
			found = true
		}
	}

	if !found {
		return CropError{Code: CropNoteErrorNotFound}
	}

	return nil
}

// CalculateDaysSinceSeeding will find how long since its been seeded
// It basically tell use the days since this crop is created.
func (c Crop) CalculateDaysSinceSeeding() int {
	now := time.Now()

	diff := now.Sub(c.CreatedDate)

	days := int(diff.Hours()) / 24

	return days
}
