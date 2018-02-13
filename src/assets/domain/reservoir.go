package domain

import (
	"time"

	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/helper/validationhelper"
	uuid "github.com/satori/go.uuid"
)

// Reservoir is entity that provides the operation that farm owner or his/her staff
// can do with the reservoir in a farm.
type Reservoir struct {
	UID         uuid.UUID                   `json:"uid"`
	Name        string                      `json:"name"`
	WaterSource WaterSource                 `json:"water_source"`
	FarmUID     uuid.UUID                   `json:"-"`
	Notes       map[uuid.UUID]ReservoirNote `json:"-"`
	CreatedDate time.Time                   `json:"created_date"`

	// Events
	Version            int
	UncommittedChanges []interface{}
}

type ReservoirService interface {
	FindFarmByID(farmUID uuid.UUID) ServiceResult
}

type ServiceResult struct {
	Result interface{}
	Error  error
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
	Capacity float32 `json:"capacity"`
}

func (b Bucket) Type() string {
	return BucketType
}

// Tap is value object attached to the Reservoir.waterSource domain.
type Tap struct {
}

func (t Tap) Type() string {
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

func (state *Reservoir) TrackChange(event interface{}) {
	state.UncommittedChanges = append(state.UncommittedChanges, event)
	state.Transition(event)
}

func (state *Reservoir) Transition(event interface{}) {
	switch e := event.(type) {
	case ReservoirCreated:
		state.UID = e.UID
		state.Name = e.Name
		state.WaterSource = e.WaterSource
		state.FarmUID = e.FarmUID
		state.CreatedDate = e.CreatedDate

	case ReservoirWaterSourceChanged:
		state.WaterSource = e.WaterSource

	case ReservoirNoteAdded:
		if len(state.Notes) == 0 {
			state.Notes = make(map[uuid.UUID]ReservoirNote)
		}

		state.Notes[e.UID] = ReservoirNote{
			UID:         e.UID,
			Content:     e.Content,
			CreatedDate: e.CreatedDate,
		}

	case ReservoirNoteRemoved:
		delete(state.Notes, e.UID)

	}
}

// CreateReservoir registers a new Reservoir.
func CreateReservoir(reservoirService ReservoirService, farmUID uuid.UUID, name string, waterSourceType string, capacity float32) (*Reservoir, error) {
	serviceResult := reservoirService.FindFarmByID(farmUID)
	if serviceResult.Error != nil {
		return nil, serviceResult.Error
	}

	farm, ok := serviceResult.Result.(query.FarmReadQueryResult)
	if !ok {
		return nil, ReservoirError{ReservoirErrorFarmNotFound}
	}

	if farm.UID == (uuid.UUID{}) {
		return nil, ReservoirError{ReservoirErrorFarmNotFound}
	}

	err := validateReservoirName(name)
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
		FarmUID:     farm.UID,
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
		WaterSource: ws,
	})

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

// ChangeName is used to change Reservoir Name.
func (r *Reservoir) ChangeName(name string) error {
	err := validateReservoirName(name)
	if err != nil {
		return err
	}

	r.Name = name

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
		UID:         uid,
		Content:     content,
		CreatedDate: time.Now(),
	})

	return nil
}

func (r *Reservoir) RemoveNote(uid string) error {
	if uid == "" {
		return ReservoirError{Code: ReservoirNoteErrorInvalidContent}
	}

	uuid, err := uuid.FromString(uid)
	if err != nil {
		return ReservoirError{Code: ReservoirNoteErrorNotFound}
	}

	found := false
	for _, v := range r.Notes {
		if v.UID == uuid {
			found = true
		}
	}

	if !found {
		return ReservoirError{Code: ReservoirNoteErrorNotFound}
	}

	r.TrackChange(ReservoirNoteRemoved{
		UID: uuid,
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
