package sqlite

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/tasks/query"
)

type AreaQuerySqlite struct {
	DB *sql.DB
}

func NewAreaQuerySqlite(db *sql.DB) query.Area {
	return AreaQuerySqlite{DB: db}
}

func (s AreaQuerySqlite) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		rowsData := struct {
			UID  string
			Name string
		}{}
		area := query.TaskAreaResult{}

		s.DB.QueryRow(`SELECT UID, NAME
			FROM AREA_READ WHERE UID = ?`, uid).Scan(&rowsData.UID, &rowsData.Name)

		areaUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		area.UID = areaUID
		area.Name = rowsData.Name

		result <- query.Result{Result: area}

		close(result)
	}()

	return result
}
