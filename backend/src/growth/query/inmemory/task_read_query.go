package inmemory

import (
	"errors"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/growth/query"
	tasksdomain "github.com/usetania/tania-core/src/tasks/domain"
	"github.com/usetania/tania-core/src/tasks/storage"
)

type TaskReadQueryInMemory struct {
	Storage *storage.TaskReadStorage
}

func NewTaskReadQueryInMemory(s *storage.TaskReadStorage) query.TaskReadQuery {
	return TaskReadQueryInMemory{Storage: s}
}

func (s TaskReadQueryInMemory) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		task := query.CropTaskQueryResult{}

		for _, val := range s.Storage.TaskReadMap {
			if val.UID == uid {
				task.UID = uid
				task.Title = val.Title
				task.Description = val.Description
				task.Category = val.Category
				task.Status = val.Status
				task.Domain = val.Domain

				if val.Domain == "CROP" {
					tdc, ok := val.DomainDetails.(tasksdomain.TaskDomainCrop)
					if !ok {
						result <- query.Result{Error: errors.New("error type assertion")}
					}

					task.AreaUID = *tdc.AreaID
					task.MaterialUID = *tdc.MaterialID
				}
			}
		}

		result <- query.Result{Result: task}

		close(result)
	}()

	return result
}
