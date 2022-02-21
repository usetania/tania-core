package repository

import (
	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/growth/domain"
	"github.com/usetania/tania-core/src/growth/storage"
)

// Result is a struct to wrap repository result
// so its easy to use it in channel.
type Result struct {
	Result interface{}
	Error  error
}

type CropEvent interface {
	Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error
}

type CropRead interface {
	Save(cropRead *storage.CropRead) <-chan error
}

func NewCropBatchFromHistory(events []storage.CropEvent) *domain.Crop {
	state := &domain.Crop{}
	for _, v := range events {
		state.Transition(v.Event)
		state.Version++
	}

	return state
}

type CropActivity interface {
	Save(cropActivity *storage.CropActivity, isUpdate bool) <-chan error
}
