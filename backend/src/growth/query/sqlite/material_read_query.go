package sqlite

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/growth/query"
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

func (q MaterialReadQuerySqlite) FindByID(materialUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		materialQueryResult := query.CropMaterialQueryResult{}
		rowsData := materialReadResult{}

		err := q.DB.QueryRow(`SELECT UID, NAME, TYPE, TYPE_DATA FROM MATERIAL_READ
			WHERE UID = ?`, materialUID).Scan(
			&rowsData.UID,
			&rowsData.Name,
			&rowsData.Type,
			&rowsData.TypeData,
		)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Error: err}
		}

		if errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Result: materialQueryResult}
		}

		materialUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		materialQueryResult.UID = materialUID
		materialQueryResult.Name = rowsData.Name
		materialQueryResult.TypeCode = rowsData.Type
		materialQueryResult.PlantTypeCode = rowsData.TypeData

		result <- query.Result{Result: materialQueryResult}
		close(result)
	}()

	return result
}

func (q MaterialReadQuerySqlite) FindMaterialByPlantTypeCodeAndName(plantTypeCode, name string) <-chan query.Result {
	result := make(chan query.Result)

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

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Error: err}
		}

		if errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Result: materialQueryResult}
		}

		materialUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		materialQueryResult.UID = materialUID
		materialQueryResult.Name = rowsData.Name
		materialQueryResult.TypeCode = rowsData.Type
		materialQueryResult.PlantTypeCode = rowsData.TypeData

		result <- query.Result{Result: materialQueryResult}
		close(result)
	}()

	return result
}
