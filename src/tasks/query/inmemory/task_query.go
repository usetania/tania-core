package inmemory

import (
	"github.com/Tanibox/tania-server/src/tasks/domain"
	"github.com/Tanibox/tania-server/src/tasks/query"
	"github.com/Tanibox/tania-server/src/tasks/storage"
)

type TaskQuery interface {
	FindTasksByAssetID(assetID string) <-chan query.QueryResult
}

type TaskQueryInMemory struct {
	Storage *storage.TaskStorage
}

func NewTaskQueryInMemory(s *storage.TaskStorage) TaskQuery {
	return TaskQueryInMemory{Storage: s}
}

func (s TaskQueryInMemory) FindTasksByAssetID(assetID string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		tasks := []domain.Task{}
		for _, val := range s.Storage.TaskMap {
			if val.AssetID.String() == assetID {
				tasks = append(tasks, val)
			}
		}

		result <- query.QueryResult{Result: tasks}

		close(result)
	}()

	return result
}
