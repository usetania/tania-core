package mysql

import (
	"database/sql"
	"time"

	"github.com/Tanibox/tania-core/src/assets/query"
	"github.com/Tanibox/tania-core/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type FarmReadQueryMysql struct {
	DB *sql.DB
}

func NewFarmReadQueryMysql(db *sql.DB) query.FarmReadQuery {
	return FarmReadQueryMysql{DB: db}
}

type farmReadResult struct {
	UID         []byte
	Name        string
	Latitude    string
	Longitude   string
	Type        string
	Country     string
	City        string
	IsActive    int
	CreatedDate time.Time
}

func (s FarmReadQueryMysql) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		farmRead := storage.FarmRead{}
		rowsData := farmReadResult{}

		err := s.DB.QueryRow("SELECT * FROM FARM_READ WHERE UID = ?", uid.Bytes()).Scan(
			&rowsData.UID,
			&rowsData.Name,
			&rowsData.Latitude,
			&rowsData.Longitude,
			&rowsData.Type,
			&rowsData.Country,
			&rowsData.City,
			&rowsData.IsActive,
			&rowsData.CreatedDate,
		)

		if err != nil && err != sql.ErrNoRows {
			result <- query.QueryResult{Error: err}
		}

		if err == sql.ErrNoRows {
			result <- query.QueryResult{Result: farmRead}
		}

		farmUID, err := uuid.FromBytes(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		farmRead = storage.FarmRead{
			UID:         farmUID,
			Name:        rowsData.Name,
			Latitude:    rowsData.Latitude,
			Longitude:   rowsData.Longitude,
			Type:        rowsData.Type,
			Country:     rowsData.Country,
			City:        rowsData.City,
			IsActive:    rowsData.IsActive != 0,
			CreatedDate: rowsData.CreatedDate,
		}

		result <- query.QueryResult{Result: farmRead}
		close(result)
	}()

	return result
}

func (s FarmReadQueryMysql) FindAll() <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		farmReads := []storage.FarmRead{}
		rowsData := farmReadResult{}

		rows, err := s.DB.Query("SELECT * FROM FARM_READ ORDER BY CREATED_DATE ASC")
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		for rows.Next() {
			err = rows.Scan(
				&rowsData.UID,
				&rowsData.Name,
				&rowsData.Latitude,
				&rowsData.Longitude,
				&rowsData.Type,
				&rowsData.Country,
				&rowsData.City,
				&rowsData.IsActive,
				&rowsData.CreatedDate,
			)

			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			farmUID, err := uuid.FromBytes(rowsData.UID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			farmReads = append(farmReads, storage.FarmRead{
				UID:         farmUID,
				Name:        rowsData.Name,
				Latitude:    rowsData.Latitude,
				Longitude:   rowsData.Longitude,
				Type:        rowsData.Type,
				Country:     rowsData.Country,
				City:        rowsData.City,
				IsActive:    rowsData.IsActive != 0,
				CreatedDate: rowsData.CreatedDate,
			})
		}

		result <- query.QueryResult{Result: farmReads}
		close(result)
	}()

	return result
}
