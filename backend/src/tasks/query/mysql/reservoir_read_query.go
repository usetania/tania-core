package mysql

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/tasks/query"
)

type ReservoirQueryMysql struct {
	DB *sql.DB
}

func NewReservoirQueryMysql(db *sql.DB) query.Reservoir {
	return ReservoirQueryMysql{DB: db}
}

func (s ReservoirQueryMysql) FindReservoirByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		rowsData := struct {
			UID  []byte
			Name string
		}{}
		reservoir := query.TaskReservoirResult{}

		s.DB.QueryRow(`SELECT UID, NAME
			FROM RESERVOIR_READ WHERE UID = ?`, uid.Bytes()).Scan(&rowsData.UID, &rowsData.Name)

		reservoirUID, err := uuid.FromBytes(rowsData.UID)
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
