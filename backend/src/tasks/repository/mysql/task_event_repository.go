package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/helper/structhelper"
	"github.com/usetania/tania-core/src/tasks/decoder"
	"github.com/usetania/tania-core/src/tasks/repository"
)

type TaskEventRepositoryMysql struct {
	DB *sql.DB
}

func NewTaskEventRepositoryMysql(s *sql.DB) repository.TaskEvent {
	return &TaskEventRepositoryMysql{DB: s}
}

func (s *TaskEventRepositoryMysql) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
	result := make(chan error)

	go func() {
		for _, v := range events {
			latestVersion++

			stmt, err := s.DB.Prepare(`INSERT INTO TASK_EVENT
				(TASK_UID, VERSION, CREATED_DATE, EVENT)
				VALUES (?, ?, ?, ?)`)
			if err != nil {
				result <- err
			}

			e, err := json.Marshal(decoder.InterfaceWrapper{
				Name: structhelper.GetName(v),
				Data: v,
			})
			if err != nil {
				panic(err)
			}

			_, err = stmt.Exec(uid.Bytes(), latestVersion, time.Now(), e)
			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}
