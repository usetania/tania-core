package domain

import (
	"time"

	"github.com/Tanibox/tania-server/src/helper/validationhelper"
	uuid "github.com/satori/go.uuid"
)

// Reservoir is entity that provides the operation that farm owner or his/her staff
// can do with the reservoir in a farm.
type Reservoir struct {
	UID         uuid.UUID       `json:"uid"`
	Name        string          `json:"name"`
	PH          float32         `json:"ph"`
	EC          float32         `json:"ec"`
	Temperature float32         `json:"temperature"`
	WaterSource WaterSource     `json:"water_source"`
	Farm        Farm            `json:"-"`
	Notes       []ReservoirNote `json:"-"`
	CreatedDate time.Time       `json:"created_date"`

	// This is for serialization purposes
	InstalledToArea []Area `json:"-"`
}

type ReservoirNote struct {
	Content     string    `json:"content"`
	CreatedDate time.Time `json:"created_date"`
}

// CreateReservoir registers a new Reservoir.
func CreateReservoir(farm Farm, name string) (Reservoir, error) {
	err := validateReservoirName(name)
	if err != nil {
		return Reservoir{}, err
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return Reservoir{}, err
	}

	return Reservoir{
		UID:         uid,
		Name:        name,
		PH:          0,
		EC:          0,
		Temperature: 0,
		Farm:        farm,
		CreatedDate: time.Now(),
	}, nil
}

// AttachBucket attach Bucket value object to Reservoir.waterSource.
func (r *Reservoir) AttachBucket(bucket Bucket) error {
	if r.IsAttachedToWaterSource() {
		return ReservoirError{ReservoirErrorWaterSourceAlreadyAttachedCode}
	}

	r.WaterSource = bucket
	return nil
}

// AttachTap attach Tap value object to Reservoir.WaterSource.
func (r *Reservoir) AttachTap(tap Tap) error {
	if r.IsAttachedToWaterSource() {
		return ReservoirError{ReservoirErrorWaterSourceAlreadyAttachedCode}
	}

	r.WaterSource = tap
	return nil
}

// IsAttachedToTap is used to check if Reservoir is attached to Tap WaterSource.
func (r Reservoir) IsAttachedToTap() bool {
	_, ok := r.WaterSource.(Tap)
	return ok
}

// IsAttachedToBucket is used to check if Reservoir is attached to Bucket WaterSource.
func (r Reservoir) IsAttachedToBucket() bool {
	_, ok := r.WaterSource.(Bucket)
	return ok
}

// IsAttachedToWaterSource is used to check if Reservoir is attached to WaterSource.
func (r Reservoir) IsAttachedToWaterSource() bool {
	return r.WaterSource != nil
}

// MeasureCondition will measure the Reservoir condition based on its properties.
func (r Reservoir) MeasureCondition() float32 {
	if !r.IsAttachedToBucket() {
		// We can't measure non bucket reservoir
		return 0
	}

	// Do measure here
	return 1
}

// ChangeName is used to change Reservoir Name.
func (r *Reservoir) ChangeName(name string) error {
	err := validateReservoirName(name)
	if err != nil {
		return err
	}

	r.Name = name

	return nil
}

// ChangeTemperature is used to change Reservoir Temperature.
//
// Temperature change can affect the value of pH and EC,
// so we will accept pH and EC value in arguments.
func (r *Reservoir) ChangeTemperature(temperature, ph, ec float32) error {
	err := validatePH(ph)
	if err != nil {
		return err
	}

	err = validateEC(ec)
	if err != nil {
		return err
	}

	r.Temperature = temperature
	r.PH = ph
	r.EC = ec

	return nil
}

func (r *Reservoir) AddNewNote(content string) error {
	if content == "" {
		return ReservoirError{Code: ReservoirNoteErrorInvalidContent}
	}

	reservoirNote := ReservoirNote{
		Content:     content,
		CreatedDate: time.Now(),
	}

	r.Notes = append(r.Notes, reservoirNote)

	return nil
}

func (r *Reservoir) RemoveNote(content string) error {
	if content == "" {
		return ReservoirError{Code: ReservoirNoteErrorInvalidContent}
	}

	for i, v := range r.Notes {
		if v.Content == content {
			copy(r.Notes[i:], r.Notes[i+1:])
			r.Notes[len(r.Notes)-1] = ReservoirNote{}
			r.Notes = r.Notes[:len(r.Notes)-1]
		}
	}

	return nil
}

func validateReservoirName(name string) error {
	if name == "" {
		return ReservoirError{ReservoirErrorNameEmptyCode}
	}
	if !validationhelper.IsAlphanumeric(name) {
		return ReservoirError{ReservoirErrorNameAlphanumericOnlyCode}
	}
	if len(name) < 5 {
		return ReservoirError{ReservoirErrorNameNotEnoughCharacterCode}
	}
	if len(name) > 100 {
		return ReservoirError{ReservoirErrorNameExceedMaximunCharacterCode}
	}

	return nil
}

func validatePH(ph float32) error {
	if ph < 0 {
		return ReservoirError{ReservoirErrorPHInvalidCode}
	}

	return nil
}

func validateEC(ec float32) error {
	if ec <= 0 {
		return ReservoirError{ReservoirErrorECInvalidCode}
	}

	return nil
}
