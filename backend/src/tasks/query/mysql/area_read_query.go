package mysql

import (
	"database/sql"

	"github.com/Tanibox/tania-core/src/tasks/query"
	uuid "github.com/satori/go.uuid"
)

type AreaQueryMysql struct {
	DB *sql.DB
}

func NewAreaQueryMysql(db *sql.DB) query.AreaQuery {
	return AreaQueryMysql{DB: db}
}

func (s AreaQueryMysql) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		rowsData := struct {
			UID  []byte
			Name string
		}{}
		area := query.TaskAreaQueryResult{}

		err := s.DB.QueryRow(`SELECT UID, NAME
			FROM AREA_READ WHERE UID = ?`, uid.Bytes()).Scan(&rowsData.UID, &rowsData.Name)

		areaUID, err := uuid.FromBytes(rowsData.UID)
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
