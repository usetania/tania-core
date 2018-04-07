package sqlite

import (
	"database/sql"
	"time"

	"github.com/Tanibox/tania-core/src/assets/query"
	"github.com/Tanibox/tania-core/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type ReservoirReadQuerySqlite struct {
	DB *sql.DB
}

func NewReservoirReadQuerySqlite(db *sql.DB) query.ReservoirReadQuery {
	return ReservoirReadQuerySqlite{DB: db}
}

type reservoirReadResult struct {
	UID                 string
	Name                string
	WaterSourceType     string
	WaterSourceCapacity float32
	FarmUID             string
	FarmName            string
	CreatedDate         string
}

type reservoirNotesReadResult struct {
	UID          string
	ReservoirUID string
	Content      string
	CreatedDate  string
}

func (s ReservoirReadQuerySqlite) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		reservoirRead := storage.ReservoirRead{}
		rowsData := reservoirReadResult{}
		notesRowsData := reservoirNotesReadResult{}

		err := s.DB.QueryRow("SELECT * FROM RESERVOIR_READ WHERE UID = ?", uid).Scan(
			&rowsData.UID,
			&rowsData.Name,
			&rowsData.WaterSourceType,
			&rowsData.WaterSourceCapacity,
			&rowsData.FarmUID,
			&rowsData.FarmName,
			&rowsData.CreatedDate,
		)

		if err != nil && err != sql.ErrNoRows {
			result <- query.QueryResult{Error: err}
		}

		if err == sql.ErrNoRows {
			result <- query.QueryResult{Result: reservoirRead}
		}

		reservoirUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		farmUID, err := uuid.FromString(rowsData.FarmUID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		resCreatedDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		rows, err := s.DB.Query("SELECT * FROM RESERVOIR_READ_NOTES WHERE RESERVOIR_UID = ?", uid)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		notes := []storage.ReservoirNote{}
		for rows.Next() {
			rows.Scan(
				&notesRowsData.UID,
				&notesRowsData.ReservoirUID,
				&notesRowsData.Content,
				&notesRowsData.CreatedDate,
			)

			noteUID, err := uuid.FromString(notesRowsData.UID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			noteCreatedDate, err := time.Parse(time.RFC3339, notesRowsData.CreatedDate)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			notes = append(notes, storage.ReservoirNote{
				UID:         noteUID,
				Content:     notesRowsData.Content,
				CreatedDate: noteCreatedDate,
			})
		}

		reservoirRead = storage.ReservoirRead{
			UID:  reservoirUID,
			Name: rowsData.Name,
			WaterSource: storage.WaterSource{
				Type:     rowsData.WaterSourceType,
				Capacity: rowsData.WaterSourceCapacity,
			},
			Farm: storage.ReservoirFarm{
				UID:  farmUID,
				Name: rowsData.FarmName,
			},
			CreatedDate: resCreatedDate,
			Notes:       notes,
		}

		result <- query.QueryResult{Result: reservoirRead}
		close(result)
	}()

	return result
}

func (s ReservoirReadQuerySqlite) FindAllByFarm(farmUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		reservoirReads := []storage.ReservoirRead{}

		rows, err := s.DB.Query("SELECT * FROM RESERVOIR_READ WHERE FARM_UID = ?", farmUID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		for rows.Next() {
			rowsData := reservoirReadResult{}
			err = rows.Scan(
				&rowsData.UID,
				&rowsData.Name,
				&rowsData.WaterSourceType,
				&rowsData.WaterSourceCapacity,
				&rowsData.FarmUID,
				&rowsData.FarmName,
				&rowsData.CreatedDate,
			)

			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			reservoirUID, err := uuid.FromString(rowsData.UID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			farmUID, err := uuid.FromString(rowsData.FarmUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			resCreatedDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			noteRows, err := s.DB.Query("SELECT * FROM RESERVOIR_READ_NOTES WHERE RESERVOIR_UID = ?", reservoirUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			notes := []storage.ReservoirNote{}
			for noteRows.Next() {
				notesRowsData := reservoirNotesReadResult{}
				err := noteRows.Scan(
					&notesRowsData.UID,
					&notesRowsData.ReservoirUID,
					&notesRowsData.Content,
					&notesRowsData.CreatedDate,
				)

				if err != nil {
					result <- query.QueryResult{Error: err}
				}

				noteUID, err := uuid.FromString(notesRowsData.UID)
				if err != nil {
					result <- query.QueryResult{Error: err}
				}

				noteCreatedDate, err := time.Parse(time.RFC3339, notesRowsData.CreatedDate)
				if err != nil {
					result <- query.QueryResult{Error: err}
				}

				notes = append(notes, storage.ReservoirNote{
					UID:         noteUID,
					Content:     notesRowsData.Content,
					CreatedDate: noteCreatedDate,
				})
			}

			reservoirReads = append(reservoirReads, storage.ReservoirRead{
				UID:  reservoirUID,
				Name: rowsData.Name,
				WaterSource: storage.WaterSource{
					Type:     rowsData.WaterSourceType,
					Capacity: rowsData.WaterSourceCapacity,
				},
				Farm: storage.ReservoirFarm{
					UID:  farmUID,
					Name: rowsData.FarmName,
				},
				CreatedDate: resCreatedDate,
				Notes:       notes,
			})
		}

		result <- query.QueryResult{Result: reservoirReads}
		close(result)
	}()

	return result
}
