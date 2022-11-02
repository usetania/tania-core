package inmemory

import (
	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/growth/query"
	"github.com/usetania/tania-core/src/growth/storage"
)

type CropActivityQueryInMemory struct {
	Storage *storage.CropActivityStorage
}

func NewCropActivityQueryInMemory(s *storage.CropActivityStorage) query.CropActivityQuery {
	return CropActivityQueryInMemory{Storage: s}
}

func (s CropActivityQueryInMemory) FindAllByCropID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		activities := []storage.CropActivity{}

		for _, val := range s.Storage.CropActivityMap {
			if val.UID == uid {
				activities = append(activities, val)
			}
		}

		result <- query.Result{Result: activities}

		close(result)
	}()

	return result
}

func (s CropActivityQueryInMemory) FindByCropIDAndActivityType(
	uid uuid.UUID,
	activityType interface{},
) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		activity := storage.CropActivity{}

		for _, val := range s.Storage.CropActivityMap {
			var at storage.ActivityType
			switch v := activityType.(type) {
			case storage.SeedActivity:
				at = v
			}

			if at != nil && val.UID == uid {
				activity = val
			}
		}

		result <- query.Result{Result: activity}

		close(result)
	}()

	return result
}
