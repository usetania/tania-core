package mysql

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/query"
	"github.com/usetania/tania-core/src/assets/storage"
)

type ReservoirReadQueryMysql struct {
	DB *sql.DB
}

func NewReservoirReadQueryMysql(db *sql.DB) query.ReservoirRead {
	return ReservoirReadQueryMysql{DB: db}
}

type reservoirReadResult struct {
	UID                 []byte
	Name                string
	WaterSourceType     string
	WaterSourceCapacity float32
	FarmUID             []byte
	FarmName            string
	CreatedDate         time.Time
}

type reservoirNotesReadResult struct {
	UID          []byte
	ReservoirUID []byte
	Content      string
	CreatedDate  time.Time
}

func (s ReservoirReadQueryMysql) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		reservoirRead := storage.ReservoirRead{}
		rowsData := reservoirReadResult{}
		notesRowsData := reservoirNotesReadResult{}

		err := s.DB.QueryRow("SELECT * FROM RESERVOIR_READ WHERE UID = ?", uid.Bytes()).Scan(
			&rowsData.UID,
			&rowsData.Name,
			&rowsData.WaterSourceType,
			&rowsData.WaterSourceCapacity,
			&rowsData.FarmUID,
			&rowsData.FarmName,
			&rowsData.CreatedDate,
		)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Error: err}
		}

		if errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Result: reservoirRead}
		}

		reservoirUID, err := uuid.FromBytes(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		farmUID, err := uuid.FromBytes(rowsData.FarmUID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		rows, err := s.DB.Query("SELECT * FROM RESERVOIR_READ_NOTES WHERE RESERVOIR_UID = ?", uid.Bytes())
		if err != nil {
			result <- query.Result{Error: err}
		}

		notes := []storage.ReservoirNote{}

		for rows.Next() {
			err := rows.Scan(
				&notesRowsData.UID,
				&notesRowsData.ReservoirUID,
				&notesRowsData.Content,
				&notesRowsData.CreatedDate,
			)
			if err != nil {
				result <- query.Result{Error: err}
			}

			noteUID, err := uuid.FromBytes(notesRowsData.UID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			notes = append(notes, storage.ReservoirNote{
				UID:         noteUID,
				Content:     notesRowsData.Content,
				CreatedDate: notesRowsData.CreatedDate,
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
			CreatedDate: rowsData.CreatedDate,
			Notes:       notes,
		}

		result <- query.Result{Result: reservoirRead}
		close(result)
	}()

	return result
}

func (s ReservoirReadQueryMysql) FindAllByFarm(farmUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		reservoirReads := []storage.ReservoirRead{}

		rows, err := s.DB.Query("SELECT * FROM RESERVOIR_READ WHERE FARM_UID = ?", farmUID.Bytes())
		if err != nil {
			result <- query.Result{Error: err}
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
				result <- query.Result{Error: err}
			}

			reservoirUID, err := uuid.FromBytes(rowsData.UID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			farmUID, err := uuid.FromBytes(rowsData.FarmUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			noteRows, err := s.DB.Query("SELECT * FROM RESERVOIR_READ_NOTES WHERE RESERVOIR_UID = ?", reservoirUID.Bytes())
			if err != nil {
				result <- query.Result{Error: err}
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
					result <- query.Result{Error: err}
				}

				noteUID, err := uuid.FromBytes(notesRowsData.UID)
				if err != nil {
					result <- query.Result{Error: err}
				}

				notes = append(notes, storage.ReservoirNote{
					UID:         noteUID,
					Content:     notesRowsData.Content,
					CreatedDate: notesRowsData.CreatedDate,
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
				CreatedDate: rowsData.CreatedDate,
				Notes:       notes,
			})
		}

		result <- query.Result{Result: reservoirReads}
		close(result)
	}()

	return result
}
