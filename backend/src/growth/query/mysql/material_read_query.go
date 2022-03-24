package sqlite

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/growth/query"
)

type MaterialReadQueryMysql struct {
	DB *sql.DB
}

func NewMaterialReadQueryMysql(db *sql.DB) query.MaterialReadQuery {
	return MaterialReadQueryMysql{DB: db}
}

type materialReadResult struct {
	UID      []byte
	Name     string
	Type     string
	TypeData string
}

func (q MaterialReadQueryMysql) FindByID(materialUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		materialQueryResult := query.CropMaterialQueryResult{}
		rowsData := materialReadResult{}

		err := q.DB.QueryRow(`SELECT UID, NAME, TYPE, TYPE_DATA FROM MATERIAL_READ
			WHERE UID = ?`, materialUID.Bytes()).Scan(
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

		materialUID, err := uuid.FromBytes(rowsData.UID)
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

func (q MaterialReadQueryMysql) FindMaterialByPlantTypeCodeAndName(plantTypeCode, name string) <-chan query.Result {
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

		materialUID, err := uuid.FromBytes(rowsData.UID)
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
