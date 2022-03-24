package domain

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/helper/validationhelper"
)

// Reservoir is entity that provides the operation that farm owner or his/her staff
// can do with the reservoir in a farm.
type Reservoir struct {
	UID         uuid.UUID
	Name        string
	WaterSource WaterSource
	FarmUID     uuid.UUID
	Notes       map[uuid.UUID]ReservoirNote
	CreatedDate time.Time

	// Events
	Version            int
	UncommittedChanges []interface{}
}

type ReservoirService interface {
	FindFarmByID(farmUID uuid.UUID) (ReservoirFarmServiceResult, error)
}

type ReservoirFarmServiceResult struct {
	UID  uuid.UUID
	Name string
}

const (
	BucketType = "BUCKET"
	TapType    = "TAP"
)

type WaterSource interface {
	Type() string
}

// Bucket is value object attached to the Reservoir.waterSource.
type Bucket struct {
	Capacity float32
}

func (Bucket) Type() string {
	return BucketType
}

// Tap is value object attached to the Reservoir.waterSource domain.
type Tap struct{}

func (Tap) Type() string {
	return TapType
}

// CreateBucket registers a new Bucket.
func CreateBucket(capacity float32) (Bucket, error) {
	if capacity <= 0 {
		return Bucket{}, ReservoirError{ReservoirErrorBucketCapacityInvalidCode}
	}

	return Bucket{Capacity: capacity}, nil
}

// CreateTap registers a new Tab.
func CreateTap() (Tap, error) {
	return Tap{}, nil
}

type ReservoirNote struct {
	UID         uuid.UUID `json:"uid"`
	Content     string    `json:"content"`
	CreatedDate time.Time `json:"created_date"`
}

func (r *Reservoir) TrackChange(event interface{}) {
	r.UncommittedChanges = append(r.UncommittedChanges, event)
	r.Transition(event)
}

func (r *Reservoir) Transition(event interface{}) {
	switch e := event.(type) {
	case ReservoirCreated:
		r.UID = e.UID
		r.Name = e.Name
		r.WaterSource = e.WaterSource
		r.FarmUID = e.FarmUID
		r.CreatedDate = e.CreatedDate

	case ReservoirWaterSourceChanged:
		r.WaterSource = e.WaterSource

	case ReservoirNameChanged:
		r.Name = e.Name

	case ReservoirNoteAdded:
		if len(r.Notes) == 0 {
			r.Notes = make(map[uuid.UUID]ReservoirNote)
		}

		r.Notes[e.UID] = ReservoirNote{
			UID:         e.UID,
			Content:     e.Content,
			CreatedDate: e.CreatedDate,
		}

	case ReservoirNoteRemoved:
		delete(r.Notes, e.UID)
	}
}

// CreateReservoir registers a new Reservoir.
func CreateReservoir(rs ReservoirService, farmUID uuid.UUID, name, waterSourceType string, capacity float32) (*Reservoir, error) {
	farmServiceResult, err := rs.FindFarmByID(farmUID)
	if err != nil {
		return nil, err
	}

	if farmServiceResult.UID == (uuid.UUID{}) {
		return nil, ReservoirError{ReservoirErrorFarmNotFound}
	}

	err = validateReservoirName(name)
	if err != nil {
		return nil, err
	}

	ws, err := validateWaterSource(waterSourceType, capacity)
	if err != nil {
		return nil, err
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	initial := &Reservoir{
		UID:         uid,
		Name:        name,
		WaterSource: ws,
		FarmUID:     farmServiceResult.UID,
		CreatedDate: time.Now(),
	}

	initial.TrackChange(ReservoirCreated{
		UID:         initial.UID,
		Name:        initial.Name,
		WaterSource: initial.WaterSource,
		FarmUID:     initial.FarmUID,
		CreatedDate: initial.CreatedDate,
	})

	return initial, nil
}

func (r *Reservoir) ChangeWaterSource(waterSourceType string, capacity float32) error {
	ws, err := validateWaterSource(waterSourceType, capacity)
	if err != nil {
		return err
	}

	r.TrackChange(ReservoirWaterSourceChanged{
		ReservoirUID: r.UID,
		WaterSource:  ws,
	})

	return nil
}

// ChangeName is used to change Reservoir Name.
func (r *Reservoir) ChangeName(name string) error {
	err := validateReservoirName(name)
	if err != nil {
		return err
	}

	r.TrackChange(ReservoirNameChanged{
		ReservoirUID: r.UID,
		Name:         name,
	})

	return nil
}

func (r *Reservoir) AddNewNote(content string) error {
	if content == "" {
		return ReservoirError{Code: ReservoirNoteErrorInvalidContent}
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	r.TrackChange(ReservoirNoteAdded{
		ReservoirUID: r.UID,
		UID:          uid,
		Content:      content,
		CreatedDate:  time.Now(),
	})

	return nil
}

func (r *Reservoir) RemoveNote(uid uuid.UUID) error {
	if uid == (uuid.UUID{}) {
		return ReservoirError{Code: ReservoirNoteErrorInvalidContent}
	}

	found := false

	for _, v := range r.Notes {
		if v.UID == uid {
			found = true
		}
	}

	if !found {
		return ReservoirError{Code: ReservoirNoteErrorNotFound}
	}

	r.TrackChange(ReservoirNoteRemoved{
		ReservoirUID: r.UID,
		UID:          uid,
	})

	return nil
}

func validateWaterSource(waterSourceType string, capacity float32) (WaterSource, error) {
	var ws WaterSource

	if waterSourceType == BucketType {
		b, err := CreateBucket(capacity)
		if err != nil {
			return nil, err
		}

		ws = b
	} else if waterSourceType == TapType {
		t, err := CreateTap()
		if err != nil {
			return nil, err
		}
		ws = t
	}

	return ws, nil
}

func validateReservoirName(name string) error {
	if name == "" {
		return ReservoirError{ReservoirErrorNameEmptyCode}
	}

	if !validationhelper.IsAlphanumSpaceHyphenUnderscore(name) {
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
