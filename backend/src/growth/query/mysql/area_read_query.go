package sqlite

import (
	"database/sql"

	"github.com/Tanibox/tania-core/src/growth/query"
	uuid "github.com/satori/go.uuid"
)

type AreaReadQueryMysql struct {
	DB *sql.DB
}

func NewAreaReadQueryMysql(db *sql.DB) query.AreaReadQuery {
	return AreaReadQueryMysql{DB: db}
}

type areaReadResult struct {
	UID      []byte
	Name     string
	Size     float32
	SizeUnit string
	Type     string
	Location string
	FarmUID  []byte
}

func (s AreaReadQueryMysql) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		areaQueryResult := query.CropAreaQueryResult{}
		rowsData := areaReadResult{}

		err := s.DB.QueryRow(`SELECT UID, NAME, SIZE, SIZE_UNIT, TYPE, LOCATION, FARM_UID
			FROM AREA_READ WHERE UID = ?`, uid.Bytes()).Scan(
			&rowsData.UID,
			&rowsData.Name,
			&rowsData.Size,
			&rowsData.SizeUnit,
			&rowsData.Type,
			&rowsData.Location,
			&rowsData.FarmUID,
		)

		if err != nil && err != sql.ErrNoRows {
			result <- query.QueryResult{Error: err}
		}

		if err == sql.ErrNoRows {
			result <- query.QueryResult{Result: areaQueryResult}
		}

		areaUID, err := uuid.FromBytes(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		farmUID, err := uuid.FromBytes(rowsData.FarmUID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		areaQueryResult.UID = areaUID
		areaQueryResult.Name = rowsData.Name
		areaQueryResult.Size.Value = rowsData.Size
		areaQueryResult.Size.Symbol = rowsData.SizeUnit
		areaQueryResult.Type = rowsData.Type
		areaQueryResult.Location = rowsData.Location
		areaQueryResult.FarmUID = farmUID

		result <- query.QueryResult{Result: areaQueryResult}
		close(result)
	}()

	return result
}
