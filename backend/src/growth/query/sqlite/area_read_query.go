package sqlite

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/growth/query"
)

type AreaReadQuerySqlite struct {
	DB *sql.DB
}

func NewAreaReadQuerySqlite(db *sql.DB) query.AreaReadQuery {
	return AreaReadQuerySqlite{DB: db}
}

type areaReadResult struct {
	UID      string
	Name     string
	Size     float32
	SizeUnit string
	Type     string
	Location string
	FarmUID  string
}

func (s AreaReadQuerySqlite) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		areaQueryResult := query.CropAreaQueryResult{}
		rowsData := areaReadResult{}

		err := s.DB.QueryRow(`SELECT UID, NAME, SIZE, SIZE_UNIT, TYPE, LOCATION, FARM_UID
			FROM AREA_READ WHERE UID = ?`, uid).Scan(
			&rowsData.UID,
			&rowsData.Name,
			&rowsData.Size,
			&rowsData.SizeUnit,
			&rowsData.Type,
			&rowsData.Location,
			&rowsData.FarmUID,
		)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Error: err}
		}

		if errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Result: areaQueryResult}
		}

		areaUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		farmUID, err := uuid.FromString(rowsData.FarmUID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		areaQueryResult.UID = areaUID
		areaQueryResult.Name = rowsData.Name
		areaQueryResult.Size.Value = rowsData.Size
		areaQueryResult.Size.Symbol = rowsData.SizeUnit
		areaQueryResult.Type = rowsData.Type
		areaQueryResult.Location = rowsData.Location
		areaQueryResult.FarmUID = farmUID

		result <- query.Result{Result: areaQueryResult}
		close(result)
	}()

	return result
}
