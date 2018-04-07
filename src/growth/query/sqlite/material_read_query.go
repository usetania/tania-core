package sqlite

import (
	"database/sql"

	"github.com/Tanibox/tania-core/src/growth/query"
	uuid "github.com/satori/go.uuid"
)

type MaterialReadQuerySqlite struct {
	DB *sql.DB
}

func NewMaterialReadQuerySqlite(db *sql.DB) query.MaterialReadQuery {
	return MaterialReadQuerySqlite{DB: db}
}

type materialReadResult struct {
	UID      string
	Name     string
	Type     string
	TypeData string
}

func (s MaterialReadQuerySqlite) FindByID(materialUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		materialQueryResult := query.CropMaterialQueryResult{}
		rowsData := materialReadResult{}

		err := s.DB.QueryRow(`SELECT UID, NAME, TYPE, TYPE_DATA FROM MATERIAL_READ
			WHERE UID = ?`, materialUID).Scan(
			&rowsData.UID,
			&rowsData.Name,
			&rowsData.Type,
			&rowsData.TypeData,
		)

		if err != nil && err != sql.ErrNoRows {
			result <- query.QueryResult{Error: err}
		}

		if err == sql.ErrNoRows {
			result <- query.QueryResult{Result: materialQueryResult}
		}

		materialUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		materialQueryResult.UID = materialUID
		materialQueryResult.Name = rowsData.Name
		materialQueryResult.TypeCode = rowsData.Type
		materialQueryResult.PlantTypeCode = rowsData.TypeData

		result <- query.QueryResult{Result: materialQueryResult}
		close(result)
	}()

	return result
}

func (q MaterialReadQuerySqlite) FindMaterialByPlantTypeCodeAndName(plantTypeCode string, name string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		materialQueryResult := query.CropMaterialQueryResult{}
		rowsData := materialReadResult{}

		err := q.DB.QueryRow(`SELECT UID, NAME, TYPE, TYPE_DATA FROM MATERIAL_READ
			WHERE TYPE_DATA = ? AND NAME = ?`, plantTypeCode, name).Scan(
			&rowsData.UID,
			&rowsData.Name,
			&rowsData.Type,
			&rowsData.TypeData,
		)

		if err != nil && err != sql.ErrNoRows {
			result <- query.QueryResult{Error: err}
		}

		if err == sql.ErrNoRows {
			result <- query.QueryResult{Result: materialQueryResult}
		}

		materialUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		materialQueryResult.UID = materialUID
		materialQueryResult.Name = rowsData.Name
		materialQueryResult.TypeCode = rowsData.Type
		materialQueryResult.PlantTypeCode = rowsData.TypeData

		result <- query.QueryResult{Result: materialQueryResult}
		close(result)
	}()

	return result
}
