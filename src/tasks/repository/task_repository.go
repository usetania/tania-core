package repository

import (
	domain "github.com/Tanibox/tania-server/src/tasks/domain"
	storage "github.com/Tanibox/tania-server/src/tasks/storage"
	uuid "github.com/satori/go.uuid"
)

type TaskRepository interface {
	Save(val *domain.Task) <-chan error
	FindAll() <-chan RepositoryResult
	FindByID(uid string) <-chan RepositoryResult
}

// TaskRepositoryInMemory is in-memory TaskRepository db implementation
type TaskRepositoryInMemory struct {
	Storage *storage.TaskStorage
}

func NewTaskRepositoryInMemory(s *storage.TaskStorage) TaskRepository {
	return &TaskRepositoryInMemory{Storage: s}
}

func (f *TaskRepositoryInMemory) FindAll() <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.Storage.Lock.RLock()
		defer f.Storage.Lock.RUnlock()

		Tasks := []domain.Task{}
		for _, val := range f.Storage.TaskMap {
			Tasks = append(Tasks, val)
		}

		result <- RepositoryResult{Result: Tasks}

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

