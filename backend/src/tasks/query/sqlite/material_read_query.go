package sqlite

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/tasks/query"
)

type MaterialQuerySqlite struct {
	DB *sql.DB
}

func NewMaterialQuerySqlite(db *sql.DB) query.Material {
	return MaterialQuerySqlite{DB: db}
}

func (s MaterialQuerySqlite) FindMaterialByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		rowsData := struct {
			UID      string
			Name     string
			Type     string
			TypeData string
		}{}
		material := query.TaskMaterialResult{}

		s.DB.QueryRow(`SELECT UID, NAME, TYPE, TYPE_DATA
			FROM MATERIAL_READ WHERE UID = ?`, uid).Scan(&rowsData.UID, &rowsData.Name, &rowsData.Type, &rowsData.TypeData)

		materialUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		material.UID = materialUID
		material.Name = rowsData.Name
		material.TypeCode = rowsData.Type
		material.DetailedTypeCode = rowsData.TypeData

		result <- query.Result{Result: material}

		close(result)
	}()

	return result
}
