package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type CropActivity struct {
	UID         uuid.UUID
	CropUID     uuid.UUID
	Description string
	CreatedDate time.Time
	Activity    Activity
}

type Activity interface {
	Status() string
}

type Move struct {
	SourceAreaUID                  uuid.UUID
	CurrentSourceAreaQuantity      int
	DestinationAreaUID             uuid.UUID
	CurrentDestinationAreaQuantity int
	QualityToBeMoved               int
}

func (m Move) Status() string {
	return "MOVE"
}

// CheckArea checks if the source and destination area is valid based on business rules
func (m Move) CheckArea(sourceAreaType string, destinationAreaType string) error {
	if sourceAreaType == "seeding" && destinationAreaType == "growing" {
		return nil
	}
	if sourceAreaType == destinationAreaType {
		return nil
	}

	return CropActivityError{Code: CropActivityMoveErrorInvalidArea}
}

func (m Move) CheckQuantity(currentSourceAreaQuantity int, quantityToBeMoved int) error {
	if currentSourceAreaQuantity < quantityToBeMoved {
		return CropActivityError{Code: CropActivityMoveErrorInvalidQuantity}
	}

	return nil
}

type Dump struct {
	AreaUID        uuid.UUID
	DumpedQuantity int
}

func (d Dump) Status() string {
	return "DUMP"
}

type Harvest struct {
	AreaUID  uuid.UUID
	Type     HarvestType
	Quantity HarvestedQuantity
}

type HarvestType interface {
}

type Partial struct {
}

type All struct {
}

type HarvestedQuantity interface {
	Unit() string
}

type Gram struct {
	Quantity int
}

func (g Gram) Unit() string {
	return "gr"
}

type Kilogram struct {
	Quantity int
}

func (kg Kilogram) Unit() string {
	return "kg"
}

type UploadPhoto struct {
	Filename string `json:"filename"`
	MimeType string `json:"mime_type"`
	Size     int    `json:"size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

func (up UploadPhoto) Status() string {
	return "UPLOAD_PHOTO"
}

type Fertilize struct {
	AreaUID      uuid.UUID
	InventoryUID uuid.UUID
}

func (f Fertilize) Status() string {
	return "FERTILIZE"
}

type Pesticide struct {
	AreaUID     uuid.UUID
	InventoryID uuid.UUID
}

func (p Pesticide) Status() string {
	return "PESTICIDE"
}

type Prune struct {
	AreaUID     uuid.UUID
	InventoryID uuid.UUID
}

func (p Prune) Status() string {
	return "PRUNE"
}

type Water struct {
	AreaUID uuid.UUID
}

func (w Water) Status() string {
	return "WATER"
}
