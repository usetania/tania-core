package inmemory

import (
	"errors"

	"github.com/Tanibox/tania-core/src/growth/query"
	tasksdomain "github.com/Tanibox/tania-core/src/tasks/domain"
	"github.com/Tanibox/tania-core/src/tasks/storage"
	uuid "github.com/satori/go.uuid"
)

type TaskReadQueryInMemory struct {
	Storage *storage.TaskReadStorage
}

func NewTaskReadQueryInMemory(s *storage.TaskReadStorage) query.TaskReadQuery {
	return TaskReadQueryInMemory{Storage: s}
}

func (s TaskReadQueryInMemory) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

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
						result <- query.QueryResult{Error: errors.New("Error type assertion")}
					}

					task.AreaUID = *tdc.AreaID
					task.MaterialUID = *tdc.MaterialID
				}
			}
		}

		result <- query.QueryResult{Result: task}

		close(result)
	}()

	return result
}
