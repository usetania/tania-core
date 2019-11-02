package mysql

import (
	"database/sql"

	"github.com/Tanibox/tania-core/src/tasks/query"
	uuid "github.com/satori/go.uuid"
)

type ReservoirQueryMysql struct {
	DB *sql.DB
}

func NewReservoirQueryMysql(db *sql.DB) query.ReservoirQuery {
	return ReservoirQueryMysql{DB: db}
}

func (s ReservoirQueryMysql) FindReservoirByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		rowsData := struct {
			UID  []byte
			Name string
		}{}
		reservoir := query.TaskReservoirQueryResult{}

		err := s.DB.QueryRow(`SELECT UID, NAME
			FROM RESERVOIR_READ WHERE UID = ?`, uid.Bytes()).Scan(&rowsData.UID, &rowsData.Name)

		reservoirUID, err := uuid.FromBytes(rowsData.UID)
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
