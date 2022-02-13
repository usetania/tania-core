package mysql

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/tasks/query"
)

type AreaQueryMysql struct {
	DB *sql.DB
}

func NewAreaQueryMysql(db *sql.DB) query.Area {
	return AreaQueryMysql{DB: db}
}

func (s AreaQueryMysql) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		rowsData := struct {
			UID  []byte
			Name string
		}{}
		area := query.TaskAreaResult{}

		s.DB.QueryRow(`SELECT UID, NAME
			FROM AREA_READ WHERE UID = ?`, uid.Bytes()).Scan(&rowsData.UID, &rowsData.Name)

		areaUID, err := uuid.FromBytes(rowsData.UID)
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
