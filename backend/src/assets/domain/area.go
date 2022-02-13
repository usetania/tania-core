package domain

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/helper/validationhelper"
)

type Area struct {
	UID          uuid.UUID              `json:"uid"`
	Name         string                 `json:"name"`
	Size         AreaSize               `json:"size"`
	Type         AreaType               `json:"type"`
	Location     AreaLocation           `json:"location"`
	Photo        AreaPhoto              `json:"photo"`
	CreatedDate  time.Time              `json:"created_date"`
	Notes        map[uuid.UUID]AreaNote `json:"-"`
	ReservoirUID uuid.UUID              `json:"-"`
	FarmUID      uuid.UUID              `json:"-"`

	// Events
	Version            int
	UncommittedChanges []interface{}
}

type AreaService interface {
	FindFarmByID(farmUID uuid.UUID) (AreaFarmServiceResult, error)
	FindReservoirByID(reservoirUID uuid.UUID) (AreaReservoirServiceResult, error)
	CountCropsByAreaID(areaUID uuid.UUID) (int, error)
}

type AreaFarmServiceResult struct {
	UID  uuid.UUID
	Name string
}

type AreaReservoirServiceResult struct {
	UID  uuid.UUID
	Name string
}

const (
	AreaTypeSeeding = "SEEDING"
	AreaTypeGrowing = "GROWING"
)

type AreaType struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

func AreaTypes() []AreaType {
	return []AreaType{
		{Code: AreaTypeSeeding, Label: "Seeding"},
		{Code: AreaTypeGrowing, Label: "Growing"},
	}
}

func GetAreaType(code string) AreaType {
	for _, v := range AreaTypes() {
		if v.Code == code {
			return v
		}
	}

	return AreaType{}
}

const (
	AreaLocationOutdoor = "OUTDOOR"
	AreaLocationIndoor  = "INDOOR"
)

type AreaLocation struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func AreaLocations() []AreaLocation {
	return []AreaLocation{
		{Code: AreaLocationOutdoor, Name: "Field (Outdoor)"},
		{Code: AreaLocationIndoor, Name: "Greenhouse (Indoor)"},
	}
}

func GetAreaLocation(code string) AreaLocation {
	for _, item := range AreaLocations() {
		if item.Code == code {
			return item
		}
	}

	return AreaLocation{}
}

const (
	SquareMeter = "m2"
	Hectare     = "Ha"
)

type AreaUnit struct {
	Label  string `json:"label"`
	Symbol string `json:"symbol"`
}

func AreaUnits() []AreaUnit {
	return []AreaUnit{
		{Symbol: SquareMeter, Label: "Square Meter"},
		{Symbol: Hectare, Label: "Hectare"},
	}
}

func GetAreaUnit(symbol string) AreaUnit {
	for _, v := range AreaUnits() {
		if v.Symbol == symbol {
			return v
		}
	}

	return AreaUnit{}
}

type AreaSize struct {
	Unit  AreaUnit `json:"unit"`
	Value float32  `json:"value"`
}

type AreaPhoto struct {
	Filename string `json:"filename"`
	MimeType string `json:"mime_type"`
	Size     int    `json:"size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type AreaNote struct {
	UID         uuid.UUID `json:"uid"`
	Content     string    `json:"content"`
	CreatedDate time.Time `json:"created_date"`
}

func (a *Area) TrackChange(event interface{}) {
	a.UncommittedChanges = append(a.UncommittedChanges, event)
	a.Transition(event)
}

func (a *Area) Transition(event interface{}) {
	switch e := event.(type) {
	case AreaCreated:
		a.UID = e.UID
		a.Name = e.Name
		a.Type = e.Type
		a.Location = e.Location
		a.Size = e.Size
		a.CreatedDate = e.CreatedDate
		a.FarmUID = e.FarmUID
		a.ReservoirUID = e.ReservoirUID

	case AreaNameChanged:
		a.Name = e.Name

	case AreaSizeChanged:
		a.Size = e.Size

	case AreaTypeChanged:
		a.Type = e.Type

	case AreaLocationChanged:
		a.Location = e.Location

	case AreaReservoirChanged:
		a.ReservoirUID = e.ReservoirUID

	case AreaPhotoAdded:
		a.Photo = AreaPhoto{
			Filename: e.Filename,
			MimeType: e.MimeType,
			Size:     e.Size,
			Width:    e.Width,
			Height:   e.Height,
		}

	case AreaNoteAdded:
		if len(a.Notes) == 0 {
			a.Notes = make(map[uuid.UUID]AreaNote)
		}

		a.Notes[e.UID] = AreaNote{
			UID:         e.UID,
			Content:     e.Content,
			CreatedDate: e.CreatedDate,
		}

	case AreaNoteRemoved:
		delete(a.Notes, e.UID)
	}
}

// CreateArea registers a new area to a farm.
func CreateArea(
	areaService AreaService,
	farmUID uuid.UUID,
	reservoirUID uuid.UUID,
	name string,
	areaType string,
	size AreaSize,
	locationCode string) (*Area, error,
) {
	err := validateAreaName(name)
	if err != nil {
		return nil, err
	}

	err = validateAreaType(areaType)
	if err != nil {
		return nil, err
	}

	err = validateSize(size)
	if err != nil {
		return nil, err
	}

	al := GetAreaLocation(locationCode)
	if al == (AreaLocation{}) {
		return nil, AreaError{Code: AreaErrorInvalidAreaLocationCode}
	}

	at := GetAreaType(areaType)
	if at == (AreaType{}) {
		return nil, AreaError{Code: AreaErrorInvalidAreaTypeCode}
	}

	farm, err := areaService.FindFarmByID(farmUID)
	if err != nil {
		return nil, err
	}

	reservoir, err := areaService.FindReservoirByID(reservoirUID)
	if err != nil {
		return nil, err
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	initial := &Area{
		UID:          uid,
		Name:         name,
		Type:         at,
		Location:     al,
		Size:         size,
		FarmUID:      farm.UID,
		ReservoirUID: reservoir.UID,
		CreatedDate:  time.Now(),
	}

	initial.TrackChange(AreaCreated{
		UID:          initial.UID,
		Name:         initial.Name,
		Type:         initial.Type,
		Location:     initial.Location,
		Size:         initial.Size,
		FarmUID:      initial.FarmUID,
		ReservoirUID: initial.ReservoirUID,
		CreatedDate:  initial.CreatedDate,
	})

	return initial, nil
}

func (a *Area) ChangeName(name string) error {
	if err := validateAreaName(name); err != nil {
		return err
	}

	a.TrackChange(AreaNameChanged{
		AreaUID: a.UID,
		Name:    name,
	})

	return nil
}

// ChangeSize changes an area size.
func (a *Area) ChangeSize(size AreaSize) error {
	if err := validateSize(size); err != nil {
		return err
	}

	a.TrackChange(AreaSizeChanged{
		AreaUID: a.UID,
		Size:    size,
	})

	return nil
}

func (a *Area) ChangeType(areaService AreaService, areaType string) error {
	at := GetAreaType(areaType)
	if at == (AreaType{}) {
		return AreaError{Code: AreaErrorInvalidAreaTypeCode}
	}

	count, err := areaService.CountCropsByAreaID(a.UID)
	if err != nil {
		return err
	}

	if count > 0 {
		return AreaError{Code: AreaErrorCropAlreadyCreated}
	}

	a.TrackChange(AreaTypeChanged{
		AreaUID: a.UID,
		Type:    at,
	})

	return nil
}

// ChangeLocation changes an area location.
func (a *Area) ChangeLocation(locationCode string) error {
	v := GetAreaLocation(locationCode)
	if v == (AreaLocation{}) {
		return AreaError{Code: AreaErrorInvalidAreaLocationCode}
	}

	a.TrackChange(AreaLocationChanged{
		AreaUID:  a.UID,
		Location: v,
	})

	return nil
}

func (a *Area) ChangeReservoir(reservoirUID uuid.UUID) error {
	a.ReservoirUID = reservoirUID

	a.TrackChange(AreaReservoirChanged{
		AreaUID:      a.UID,
		ReservoirUID: reservoirUID,
	})

	return nil
}

// TODO: Do file type validation here.
func (a *Area) ChangePhoto(photo AreaPhoto) error {
	a.TrackChange(AreaPhotoAdded{
		AreaUID:  a.UID,
		Filename: photo.Filename,
		MimeType: photo.MimeType,
		Size:     photo.Size,
		Width:    photo.Width,
		Height:   photo.Height,
	})

	return nil
}

func (a *Area) AddNewNote(content string) error {
	if content == "" {
		return AreaError{Code: AreaNoteErrorInvalidContent}
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	a.TrackChange(AreaNoteAdded{
		AreaUID:     a.UID,
		UID:         uid,
		Content:     content,
		CreatedDate: time.Now(),
	})

	return nil
}

func (a *Area) RemoveNote(uid uuid.UUID) error {
	if uid == (uuid.UUID{}) {
		return AreaError{Code: AreaNoteErrorInvalidID}
	}

	found := false

	for _, v := range a.Notes {
		if v.UID == uid {
			found = true
		}
	}

	if !found {
		return AreaError{Code: AreaNoteErrorNotFound}
	}

	a.TrackChange(AreaNoteRemoved{
		AreaUID: a.UID,
		UID:     uid,
	})

	return nil
}

func validateAreaName(name string) error {
	if name == "" {
		return AreaError{AreaErrorNameEmptyCode}
	}

	if !validationhelper.IsAlphanumSpaceHyphenUnderscore(name) {
		return AreaError{AreaErrorNameAlphanumericOnlyCode}
	}

	if len(name) < 5 {
		return AreaError{AreaErrorNameNotEnoughCharacterCode}
	}

	if len(name) > 100 {
		return AreaError{AreaErrorNameExceedMaximunCharacterCode}
	}

	return nil
}

func validateSize(size AreaSize) error {
	unit := GetAreaUnit(size.Unit.Symbol)
	if unit == (AreaUnit{}) {
		return AreaError{AreaErrorInvalidSizeUnitCode}
	}

	if size.Value <= 0 {
		return AreaError{AreaErrorSizeEmptyCode}
	}

	return nil
}

func validateAreaType(areaType string) error {
	if areaType == "" {
		return AreaError{AreaErrorTypeEmptyCode}
	}

	if v := GetAreaType(areaType); v == (AreaType{}) {
		return AreaError{AreaErrorInvalidAreaTypeCode}
	}

	return nil
}
