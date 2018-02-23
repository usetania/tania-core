package sqlite

import (
	"database/sql"
	"time"

	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type FarmReadQuerySqlite struct {
	DB *sql.DB
}

func NewFarmReadQuerySqlite(db *sql.DB) query.FarmReadQuery {
	return FarmReadQuerySqlite{DB: db}
}

type farmReadResult struct {
	UID         string
	Name        string
	Latitude    string
	Longitude   string
	Type        string
	CountryCode string
	CityCode    string
	IsActive    int
	CreatedDate string
}

func (s FarmReadQuerySqlite) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		farmRead := storage.FarmRead{}
		rowsData := farmReadResult{}

		err := s.DB.QueryRow("SELECT * FROM FARM_READ WHERE UID = ?", uid).Scan(
			&rowsData.UID,
			&rowsData.Name,
			&rowsData.Latitude,
			&rowsData.Longitude,
			&rowsData.Type,
			&rowsData.CountryCode,
			&rowsData.CityCode,
			&rowsData.IsActive,
			&rowsData.CreatedDate,
		)

		if err != nil && err != sql.ErrNoRows {
			result <- query.QueryResult{Error: err}
			close(result)
		}

		if err == sql.ErrNoRows {
			result <- query.QueryResult{Result: farmRead}
			close(result)
		}

		farmUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
			close(result)
		}

		createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
		if err != nil {
			result <- query.QueryResult{Error: err}
			close(result)
		}

		farmRead = storage.FarmRead{
			UID:         farmUID,
			Name:        rowsData.Name,
			Latitude:    rowsData.Latitude,
			Longitude:   rowsData.Longitude,
			Type:        rowsData.Type,
			CountryCode: rowsData.CountryCode,
			CityCode:    rowsData.CityCode,
			IsActive:    rowsData.IsActive != 0,
			CreatedDate: createdDate,
		}

		result <- query.QueryResult{Result: farmRead}
		close(result)
	}()

	return result
}

func (s FarmReadQuerySqlite) FindAll() <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		farmReads := []storage.FarmRead{}
		rowsData := farmReadResult{}

		rows, err := s.DB.Query("SELECT * FROM FARM_READ ORDER BY CREATED_DATE ASC")
		if err != nil {
			result <- query.QueryResult{Error: err}
			close(result)
		}

		for rows.Next() {
			err = rows.Scan(
				&rowsData.UID,
				&rowsData.Name,
				&rowsData.Latitude,
				&rowsData.Longitude,
				&rowsData.Type,
				&rowsData.CountryCode,
				&rowsData.CityCode,
				&rowsData.IsActive,
				&rowsData.CreatedDate,
			)

			if err != nil {
				result <- query.QueryResult{Error: err}
				close(result)
			}

			farmUID, err := uuid.FromString(rowsData.UID)
			if err != nil {
				result <- query.QueryResult{Error: err}
				close(result)
			}

			createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.QueryResult{Error: err}
				close(result)
			}

			farmReads = append(farmReads, storage.FarmRead{
				UID:         farmUID,
				Name:        rowsData.Name,
				Latitude:    rowsData.Latitude,
				Longitude:   rowsData.Longitude,
				Type:        rowsData.Type,
				CountryCode: rowsData.CountryCode,
				CityCode:    rowsData.CityCode,
				IsActive:    rowsData.IsActive != 0,
				CreatedDate: createdDate,
			})
		}

		result <- query.QueryResult{Result: farmReads}
		close(result)
	}()

	return result
}
