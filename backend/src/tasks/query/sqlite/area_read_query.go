package sqlite

import (
	"database/sql"

	"github.com/Tanibox/tania-core/src/tasks/query"
	uuid "github.com/satori/go.uuid"
)

type AreaQuerySqlite struct {
	DB *sql.DB
}

func NewAreaQuerySqlite(db *sql.DB) query.AreaQuery {
	return AreaQuerySqlite{DB: db}
}

func (s AreaQuerySqlite) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		rowsData := struct {
			UID  string
			Name string
		}{}
		area := query.TaskAreaQueryResult{}

		err := s.DB.QueryRow(`SELECT UID, NAME
			FROM AREA_READ WHERE UID = ?`, uid).Scan(&rowsData.UID, &rowsData.Name)

		areaUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		area.UID = areaUID
		area.Name = rowsData.Name

		result <- query.QueryResult{Result: area}

		close(result)
	}()

	return result
}
