package inmemory

import (
	"github.com/Tanibox/tania-server/src/tasks/query"
	"github.com/Tanibox/tania-server/src/tasks/storage"
	uuid "github.com/satori/go.uuid"
	"strconv"
)

type TaskReadQueryInMemory struct {
	Storage *storage.TaskReadStorage
}

func NewTaskReadQueryInMemory(s *storage.TaskReadStorage) query.TaskReadQuery {
	return &TaskReadQueryInMemory{Storage: s}
}

func (r TaskReadQueryInMemory) FindAll() <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		r.Storage.Lock.RLock()
		defer r.Storage.Lock.RUnlock()

		tasks := []storage.TaskRead{}
		for _, val := range r.Storage.TaskReadMap {
			tasks = append(tasks, val)
		}

		result <- query.QueryResult{Result: tasks}

		close(result)
	}()

	return result
}

// FindByID is to find by ID
func (r TaskReadQueryInMemory) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		r.Storage.Lock.RLock()
		defer r.Storage.Lock.RUnlock()

		result <- query.QueryResult{Result: r.Storage.TaskReadMap[uid]}

		close(result)
	}()

	return result
}

func (s TaskReadQueryInMemory) FindTasksWithFilter(params map[string]string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		tasks := []storage.TaskRead{}
		for _, val := range s.Storage.TaskReadMap {
			is_match := true

			// Is Due
			if value, _ := params["is_due"]; value != "" {
				b, _ := strconv.ParseBool(value)
				if val.IsDue != b {
					is_match = false
				}
			}
			if is_match {
				// Priority
				if value, _ := params["priority"]; value != "" {
					if val.Priority != value {
						is_match = false
					}
				}
				if is_match {
					// Status
					if value, _ := params["status"]; value != "" {
						if val.Status != value {
							is_match = false
						}
					}
					if is_match {
						// Domain
						if value, _ := params["domain"]; value != "" {
							if val.Domain != value {
								is_match = false
							}
						}
						if is_match {
							// Asset ID
							if value, _ := params["asset_id"]; value != "" {
								asset_id, _ := uuid.FromString(value)
								if *val.AssetID != asset_id {
									is_match = false
								}
							}
						}
					}
				}
			}
			if is_match {
				tasks = append(tasks, val)
			}
		}

		result <- query.QueryResult{Result: tasks}

		close(result)
	}()

	return result
}
