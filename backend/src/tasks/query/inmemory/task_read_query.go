package inmemory

import (
	"github.com/Tanibox/tania-core/src/tasks/query"
	"github.com/Tanibox/tania-core/src/tasks/storage"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"time"
)

type TaskReadQueryInMemory struct {
	Storage *storage.TaskReadStorage
}

func NewTaskReadQueryInMemory(s *storage.TaskReadStorage) query.TaskReadQuery {
	return &TaskReadQueryInMemory{Storage: s}
}

func (r TaskReadQueryInMemory) FindAll(page, limit int) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		r.Storage.Lock.RLock()
		defer r.Storage.Lock.RUnlock()

		tasks := []storage.TaskRead{}

		for _, val := range r.Storage.TaskReadMap {
			tasks = append(tasks, val)
		}

		if limit != 0 {
			tasks = tasks[:limit]
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

func (s TaskReadQueryInMemory) FindTasksWithFilter(params map[string]string, page, limit int) <-chan query.QueryResult {
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
							if is_match {
								// Category
								if value, _ := params["category"]; value != "" {
									if val.Category != value {
										is_match = false
									}
								}
								if is_match {
									// Due Start Date & Due End Date
									start, _ := params["due_start"]
									end, _ := params["due_end"]

									if (start != "") && (end != "") {
										start_date, err := time.Parse(time.RFC3339Nano, start)

										if err == nil {
											end_date, err := time.Parse(time.RFC3339Nano, end)

											if err == nil {
												if !checkWithinTimeRange(start_date, end_date, *val.DueDate) {
													is_match = false
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
			if is_match {
				tasks = append(tasks, val)
			}
		}

		result <- query.QueryResult{Result: tasks}

		close(result)
	}()

	return result
}

func (q TaskReadQueryInMemory) CountAll() <-chan query.QueryResult {
  result := make(chan query.QueryResult)

  go func() {
    q.Storage.Lock.RLock()
    defer q.Storage.Lock.RUnlock()

    total := len(q.Storage.TaskReadMap)

    result <- query.QueryResult{Result: total}

    close(result)
  }()

  return result
}

func (s TaskReadQueryInMemory) CountTasksWithFilter(params map[string]string) <-chan query.QueryResult {
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
              if is_match {
                // Category
                if value, _ := params["category"]; value != "" {
                  if val.Category != value {
                    is_match = false
                  }
                }
                if is_match {
                  // Due Start Date & Due End Date
                  start, _ := params["due_start"]
                  end, _ := params["due_end"]

                  if (start != "") && (end != "") {
                    start_date, err := time.Parse(time.RFC3339Nano, start)

                    if err == nil {
                      end_date, err := time.Parse(time.RFC3339Nano, end)

                      if err == nil {
                        if !checkWithinTimeRange(start_date, end_date, *val.DueDate) {
                          is_match = false
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
      if is_match {
        tasks = append(tasks, val)
      }
    }

    result <- query.QueryResult{Result: len(tasks)}

    close(result)
  }()

  return result
}

func checkWithinTimeRange(start time.Time, end time.Time, check time.Time) bool {

	is_start := check.Equal(start)
	is_end := check.Equal(end)
	is_between := check.After(start) && check.Before(end)
	return is_start || is_end || is_between
}
