package mysql

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/tasks/query"
)

type MaterialQueryMysql struct {
	DB *sql.DB
}

func NewMaterialQueryMysql(db *sql.DB) query.Material {
	return MaterialQueryMysql{DB: db}
}

func (s MaterialQueryMysql) FindMaterialByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		rowsData := struct {
			UID      []byte
			Name     string
			Type     string
			TypeData string
		}{}
		material := query.TaskMaterialResult{}

		s.DB.QueryRow(`SELECT UID, NAME, TYPE, TYPE_DATA
			FROM MATERIAL_READ WHERE UID = ?`, uid.Bytes()).Scan(
			&rowsData.UID,
			&rowsData.Name,
			&rowsData.Type,
			&rowsData.TypeData,
		)

		materialUID, err := uuid.FromBytes(rowsData.UID)
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
