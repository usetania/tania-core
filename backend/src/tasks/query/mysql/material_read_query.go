package mysql

import (
	"database/sql"

	"github.com/Tanibox/tania-core/src/tasks/query"
	uuid "github.com/satori/go.uuid"
)

type MaterialQueryMysql struct {
	DB *sql.DB
}

func NewMaterialQueryMysql(db *sql.DB) query.MaterialQuery {
	return MaterialQueryMysql{DB: db}
}

func (s MaterialQueryMysql) FindMaterialByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		rowsData := struct {
			UID      []byte
			Name     string
			Type     string
			TypeData string
		}{}
		material := query.TaskMaterialQueryResult{}

		err := s.DB.QueryRow(`SELECT UID, NAME, TYPE, TYPE_DATA
			FROM MATERIAL_READ WHERE UID = ?`, uid.Bytes()).Scan(&rowsData.UID, &rowsData.Name, &rowsData.Type, &rowsData.TypeData)

		materialUID, err := uuid.FromBytes(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		material.UID = materialUID
		material.Name = rowsData.Name
		material.TypeCode = rowsData.Type
		material.DetailedTypeCode = rowsData.TypeData

		result <- query.QueryResult{Result: material}

		close(result)
	}()

	return result
}
