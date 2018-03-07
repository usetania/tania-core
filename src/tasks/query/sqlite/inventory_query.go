package sqlite

import (
	"database/sql"
	//"errors"
	"fmt"

	//assetsdomain "github.com/Tanibox/tania-server/src/assets/domain"
	//assetsstorage "github.com/Tanibox/tania-server/src/assets/storage"
	"github.com/Tanibox/tania-server/src/tasks/query"
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

		err := s.DB.QueryRow(`SELECT UID, NAME
			FROM MATERIAL_READ WHERE UID = ?`, uid).Scan(&rowsData.UID, &rowsData.Name)

		fmt.Println("MATERIALUID", rowsData.UID)
		materialUID, err := uuid.FromString(rowsData.UID)
		fmt.Println("MUID", materialUID)
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
