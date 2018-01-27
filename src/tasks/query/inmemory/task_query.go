package query

import (
	"github.com/Tanibox/tania-server/src/tasks/domain"
	"github.com/Tanibox/tania-server/src/tasks/storage"
)

type TaskQuery interface {
	FindTasksByAssetID(assetID string) <-chan QueryResult
}

type TaskQueryInMemory struct {
	Storage *storage.TaskStorage
}

func NewTaskQueryInMemory(s *storage.TaskStorage) TaskQuery {
	return TaskQueryInMemory{Storage: s}
}

func (s TaskQueryInMemory) FindTasksByAssetID(assetID string) <-chan QueryResult {
	result := make(chan QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		tasks := []domain.Task{}
		for _, val := range s.Storage.TaskMap {
			if val.AssetID.String() == assetID {
				tasks = append(tasks, val)
			}
		}

		result <- QueryResult{Result: tasks}

		close(result)
	}()

	return result
}
