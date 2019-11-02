package sqlite

import (
	"database/sql"

	"github.com/Tanibox/tania-core/src/growth/query"
	uuid "github.com/satori/go.uuid"
)

type FarmReadQueryMysql struct {
	DB *sql.DB
}

func NewFarmReadQueryMysql(db *sql.DB) query.FarmReadQuery {
	return FarmReadQueryMysql{DB: db}
}

type farmReadResult struct {
	UID  []byte
	Name string
}

func (s FarmReadQueryMysql) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		farmRead := query.CropFarmQueryResult{}
		rowsData := farmReadResult{}

		err := s.DB.QueryRow("SELECT UID, NAME FROM FARM_READ WHERE UID = ?", uid.Bytes()).Scan(
			&rowsData.UID,
			&rowsData.Name,
		)

		if err != nil && err != sql.ErrNoRows {
			result <- query.QueryResult{Error: err}
		}

		if err == sql.ErrNoRows {
			result <- query.QueryResult{Result: farmRead}
		}

		farmUID, err := uuid.FromBytes(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		farmRead.UID = farmUID
		farmRead.Name = rowsData.Name

		result <- query.QueryResult{Result: farmRead}
		close(result)
	}()

	return result
}
