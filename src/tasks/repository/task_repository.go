package repository

import (
	domain "github.com/Tanibox/tania-server/src/tasks/domain"
	storage "github.com/Tanibox/tania-server/src/tasks/storage"
	uuid "github.com/satori/go.uuid"
	"strconv"
)

type TaskRepository interface {
	Save(val *domain.Task) <-chan error
	FindAll() <-chan RepositoryResult
	FindByID(uid string) <-chan RepositoryResult
	FindTasksWithFilter(params map[string]string) <-chan RepositoryResult
}

// TaskRepositoryInMemory is in-memory TaskRepository db implementation
type TaskRepositoryInMemory struct {
	Storage *storage.TaskStorage
}

func NewTaskRepositoryInMemory(s *storage.TaskStorage) TaskRepository {
	return &TaskRepositoryInMemory{Storage: s}
}

func (r *TaskRepositoryInMemory) FindAll() <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		r.Storage.Lock.RLock()
		defer r.Storage.Lock.RUnlock()

		Tasks := []domain.Task{}
		for _, val := range r.Storage.TaskMap {
			Tasks = append(Tasks, val)
		}

		result <- RepositoryResult{Result: Tasks}

		close(result)
	}()

	return result
}

func (s *TaskRepositoryInMemory) FindTasksWithFilter(params map[string]string) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		tasks := []domain.Task{}
		for _, val := range s.Storage.TaskMap {

			is_match := true

			if len(params) > 0 {
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
						if value != "" && val.Priority != value {
							is_match = false
						}
					}
					if is_match {
						// Status
						if value, ok := params["status"]; ok {
							if value != "" && val.Status != value {
								is_match = false
							}
						}
						if is_match {
							// Domain
							if value, ok := params["domain"]; ok {
								if value != "" && val.Domain != value {
									is_match = false
								}
							}
							if is_match {
								// Asset ID
								if value, ok := params["asset_id"]; ok {
									if value != "" {
										asset_id, _ := uuid.FromString(value)
										if *val.AssetID != asset_id {
											is_match = false
										}
									}
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

		result <- RepositoryResult{Result: tasks}

		close(result)
	}()

	return result
}

// FindByID is to find by ID
func (r *TaskRepositoryInMemory) FindByID(uid string) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		r.Storage.Lock.RLock()
		defer r.Storage.Lock.RUnlock()

		uid, err := uuid.FromString(uid)
		if err != nil {
			result <- RepositoryResult{Error: err}
		}

		result <- RepositoryResult{Result: r.Storage.TaskMap[uid]}
	}()

	return result
}

// Save is to save
func (r *TaskRepositoryInMemory) Save(val *domain.Task) <-chan error {
	result := make(chan error)

	go func() {
		r.Storage.Lock.Lock()
		defer r.Storage.Lock.Unlock()

		r.Storage.TaskMap[val.UID] = *val

		result <- nil

		close(result)
	}()

	return result
}
