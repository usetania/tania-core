package inmemory

import (
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/tasks/query"
	"github.com/usetania/tania-core/src/tasks/storage"
)

type TaskReadQueryInMemory struct {
	Storage *storage.TaskReadStorage
}

func NewTaskReadQueryInMemory(s *storage.TaskReadStorage) query.TaskRead {
	return &TaskReadQueryInMemory{Storage: s}
}

func (q TaskReadQueryInMemory) FindAll(page, limit int) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		tasks := []storage.TaskRead{}

		for _, val := range q.Storage.TaskReadMap {
			tasks = append(tasks, val)
		}

		if limit != 0 {
			tasks = tasks[:limit]
		}

		result <- query.Result{Result: tasks}

		close(result)
	}()

	return result
}

// FindByID is to find by ID
func (q TaskReadQueryInMemory) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		result <- query.Result{Result: q.Storage.TaskReadMap[uid]}

		close(result)
	}()

	return result
}

func (q TaskReadQueryInMemory) FindTasksWithFilter(params map[string]string, page, limit int) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		tasks := []storage.TaskRead{}
		for _, val := range q.Storage.TaskReadMap {
			isMatch := true

			// Is Due
			if value := params["is_due"]; value != "" {
				b, _ := strconv.ParseBool(value)
				if val.IsDue != b {
					isMatch = false
				}
			}

			if isMatch {
				// Priority
				if value := params["priority"]; value != "" {
					if val.Priority != value {
						isMatch = false
					}
				}

				if isMatch {
					// Status
					if value := params["status"]; value != "" {
						if val.Status != value {
							isMatch = false
						}
					}

					if isMatch {
						// Domain
						if value := params["domain"]; value != "" {
							if val.Domain != value {
								isMatch = false
							}
						}

						if isMatch {
							// Asset ID
							if value := params["asset_id"]; value != "" {
								assetID, _ := uuid.FromString(value)
								if *val.AssetID != assetID {
									isMatch = false
								}
							}

							if isMatch {
								// Category
								if value := params["category"]; value != "" {
									if val.Category != value {
										isMatch = false
									}
								}

								if isMatch {
									// Due Start Date & Due End Date
									start := params["due_start"]
									end := params["due_end"]

									if (start != "") && (end != "") {
										startDate, err := time.Parse(time.RFC3339Nano, start)

										if err == nil {
											endDate, err := time.Parse(time.RFC3339Nano, end)

											if err == nil {
												if !checkWithinTimeRange(startDate, endDate, *val.DueDate) {
													isMatch = false
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}

			if isMatch {
				tasks = append(tasks, val)
			}
		}

		result <- query.Result{Result: tasks}

		close(result)
	}()

	return result
}

func (q TaskReadQueryInMemory) CountAll() <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		total := len(q.Storage.TaskReadMap)

		result <- query.Result{Result: total}

		close(result)
	}()

	return result
}

func (q TaskReadQueryInMemory) CountTasksWithFilter(params map[string]string) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		tasks := []storage.TaskRead{}
		for _, val := range q.Storage.TaskReadMap {
			isMatch := true

			// Is Due
			if value := params["is_due"]; value != "" {
				b, _ := strconv.ParseBool(value)
				if val.IsDue != b {
					isMatch = false
				}
			}

			if isMatch {
				// Priority
				if value := params["priority"]; value != "" {
					if val.Priority != value {
						isMatch = false
					}
				}

				if isMatch {
					// Status
					if value := params["status"]; value != "" {
						if val.Status != value {
							isMatch = false
						}
					}

					if isMatch {
						// Domain
						if value := params["domain"]; value != "" {
							if val.Domain != value {
								isMatch = false
							}
						}

						if isMatch {
							// Asset ID
							if value := params["asset_id"]; value != "" {
								assetID, _ := uuid.FromString(value)
								if *val.AssetID != assetID {
									isMatch = false
								}
							}

							if isMatch {
								// Category
								if value := params["category"]; value != "" {
									if val.Category != value {
										isMatch = false
									}
								}

								if isMatch {
									// Due Start Date & Due End Date
									start := params["due_start"]
									end := params["due_end"]

									if (start != "") && (end != "") {
										startDate, err := time.Parse(time.RFC3339Nano, start)

										if err == nil {
											endDate, err := time.Parse(time.RFC3339Nano, end)

											if err == nil {
												if !checkWithinTimeRange(startDate, endDate, *val.DueDate) {
													isMatch = false
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}

			if isMatch {
				tasks = append(tasks, val)
			}
		}

		result <- query.Result{Result: len(tasks)}

		close(result)
	}()

	return result
}

func checkWithinTimeRange(start, end, check time.Time) bool {
	isStart := check.Equal(start)
	isEnd := check.Equal(end)
	isBetween := check.After(start) && check.Before(end)

	return isStart || isEnd || isBetween
}
