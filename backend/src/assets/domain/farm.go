// Package domain provides the operation that farm holder can do
// to their farm
package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Farm struct {
	UID         uuid.UUID `json:"uid"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	IsActive    bool      `json:"is_active"`
	CreatedDate time.Time `json:"created_date"`

	// Events
	Version            int
	UncommittedChanges []interface{}
}

type FarmService interface {
	GetCountryNameByCode() string
}

func (state *Farm) TrackChange(event interface{}) {
	state.UncommittedChanges = append(state.UncommittedChanges, event)
	state.Transition(event)
}

func (state *Farm) Transition(event interface{}) {
	switch e := event.(type) {
	case FarmCreated:
		state.UID = e.UID
		state.UserID = e.UserID
		state.Name = e.Name
		state.IsActive = e.IsActive
		state.CreatedDate = e.CreatedDate

	case FarmNameChanged:
		state.Name = e.Name
	}
}

// CreateFarm registers a new farm to Tania
func CreateFarm(name string, userID string) (*Farm, error) {
	err := validateFarmName(name)
	if err != nil {
		return nil, err
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	parsedUserID, err := uuid.FromString(userID)
	if err != nil {
		return nil, err
	}

	initial := &Farm{}

	initial.TrackChange(FarmCreated{
		UID:         uid,
		UserID:      parsedUserID,
		Name:        name,
		IsActive:    true,
		CreatedDate: time.Now(),
	})

	return initial, nil
}

func (f *Farm) ChangeName(name string) error {
	err := validateFarmName(name)
	if err != nil {
		return err
	}

	f.TrackChange(FarmNameChanged{
		FarmUID: f.UID,
		Name:    name,
	})

	return nil
}
