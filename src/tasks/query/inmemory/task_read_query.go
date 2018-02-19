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
func (r TaskReadQueryInMemory) FindByID(uid string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		r.Storage.Lock.RLock()
		defer r.Storage.Lock.RUnlock()

		uid, err := uuid.FromString(uid)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		result <- query.QueryResult{Result: r.Storage.TaskReadMap[uid]}

		close(result)
	}()

	return result
}

func (s TaskReadQueryInMemory) QueryTasksWithFilter(params map[string]string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		tasks := []storage.TaskRead{}
		for _, val := range s.Storage.TaskReadMap {

			is_match := true

			// Is Due
			if value, ok := params["is_due"]; ok {
				b, _ := strconv.ParseBool(value)
				if val.IsDue != b {
					is_match = false
				}
			}
			if is_match {
				// Priority
				if value, ok := params["priority"]; ok {
					if val.Priority != value {
						is_match = false
					}
				}
				if is_match {
					// Status
					if value, ok := params["status"]; ok {
						if val.Status != value {
							is_match = false
						}
					}
					if is_match {
						// Domain
						if value, ok := params["domain"]; ok {
							if val.Domain != value {
								is_match = false
							}
						}
						if is_match {
							// Asset ID
							if value, ok := params["asset_id"]; ok {

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
