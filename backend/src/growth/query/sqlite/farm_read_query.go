package sqlite

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/growth/query"
)

type FarmReadQuerySqlite struct {
	DB *sql.DB
}

func NewFarmReadQuerySqlite(db *sql.DB) query.FarmReadQuery {
	return FarmReadQuerySqlite{DB: db}
}

type farmReadResult struct {
	UID  string
	Name string
}

func (s FarmReadQuerySqlite) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		farmRead := query.CropFarmQueryResult{}
		rowsData := farmReadResult{}

		err := s.DB.QueryRow("SELECT UID, NAME FROM FARM_READ WHERE UID = ?", uid).Scan(
			&rowsData.UID,
			&rowsData.Name,
		)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Error: err}
		}

		if errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Result: farmRead}
		}

		farmUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		farmRead.UID = farmUID
		farmRead.Name = rowsData.Name

		result <- query.Result{Result: farmRead}
		close(result)
	}()

	return result
}
