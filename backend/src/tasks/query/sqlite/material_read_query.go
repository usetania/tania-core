package sqlite

import (
	"database/sql"
	"github.com/Tanibox/tania-core/src/tasks/query"
	uuid "github.com/satori/go.uuid"
)

type MaterialQuerySqlite struct {
	DB *sql.DB
}

func NewMaterialQuerySqlite(db *sql.DB) query.MaterialQuery {
	return MaterialQuerySqlite{DB: db}
}

func (s MaterialQuerySqlite) FindMaterialByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		rowsData := struct {
			UID      string
			Name     string
			Type     string
			TypeData string
		}{}
		material := query.TaskMaterialQueryResult{}

		err := s.DB.QueryRow(`SELECT UID, NAME, TYPE, TYPE_DATA 
			FROM MATERIAL_READ WHERE UID = ?`, uid).Scan(&rowsData.UID, &rowsData.Name, &rowsData.Type, &rowsData.TypeData)

		materialUID, err := uuid.FromString(rowsData.UID)
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
