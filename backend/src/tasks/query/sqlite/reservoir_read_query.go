package sqlite

import (
	"database/sql"

	"github.com/Tanibox/tania-core/src/tasks/query"
	uuid "github.com/satori/go.uuid"
)

type ReservoirQuerySqlite struct {
	DB *sql.DB
}

func NewReservoirQuerySqlite(db *sql.DB) query.ReservoirQuery {
	return ReservoirQuerySqlite{DB: db}
}

func (s ReservoirQuerySqlite) FindReservoirByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		rowsData := struct {
			UID  string
			Name string
		}{}
		reservoir := query.TaskReservoirQueryResult{}

		err := s.DB.QueryRow(`SELECT UID, NAME
			FROM RESERVOIR_READ WHERE UID = ?`, uid).Scan(&rowsData.UID, &rowsData.Name)

		reservoirUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		reservoir.UID = reservoirUID
		reservoir.Name = rowsData.Name

		result <- query.QueryResult{Result: reservoir}

		close(result)
	}()

	return result
}
