package sqlite

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/tasks/query"
)

type ReservoirQuerySqlite struct {
	DB *sql.DB
}

func NewReservoirQuerySqlite(db *sql.DB) query.Reservoir {
	return ReservoirQuerySqlite{DB: db}
}

func (s ReservoirQuerySqlite) FindReservoirByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		rowsData := struct {
			UID  string
			Name string
		}{}
		reservoir := query.TaskReservoirResult{}

		s.DB.QueryRow(`SELECT UID, NAME
			FROM RESERVOIR_READ WHERE UID = ?`, uid).Scan(&rowsData.UID, &rowsData.Name)

		reservoirUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		reservoir.UID = reservoirUID
		reservoir.Name = rowsData.Name

		result <- query.Result{Result: reservoir}

		close(result)
	}()

	return result
}
