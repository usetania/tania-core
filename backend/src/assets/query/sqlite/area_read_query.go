package sqlite

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/domain"
	"github.com/usetania/tania-core/src/assets/query"
	"github.com/usetania/tania-core/src/assets/storage"
)

type AreaReadQuerySqlite struct {
	DB *sql.DB
}

func NewAreaReadQuerySqlite(db *sql.DB) query.AreaRead {
	return AreaReadQuerySqlite{DB: db}
}

type areaReadResult struct {
	UID           string
	Name          string
	Size          float32
	SizeUnit      string
	Type          string
	Location      string
	PhotoFilename string
	PhotoMimetype string
	PhotoSize     int
	PhotoWidth    int
	PhotoHeight   int
	CreatedDate   string
	ReservoirUID  string
	ReservoirName string
	FarmUID       string
	FarmName      string
}

type areaNotesReadResult struct {
	UID         string
	AreaUID     string
	Content     string
	CreatedDate string
}

func (s AreaReadQuerySqlite) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		areaRead := storage.AreaRead{}
		rowsData := areaReadResult{}
		notesRowsData := areaNotesReadResult{}

		err := s.DB.QueryRow("SELECT * FROM AREA_READ WHERE UID = ?", uid).Scan(
			&rowsData.UID,
			&rowsData.Name,
			&rowsData.SizeUnit,
			&rowsData.Size,
			&rowsData.Type,
			&rowsData.Location,
			&rowsData.PhotoFilename,
			&rowsData.PhotoMimetype,
			&rowsData.PhotoSize,
			&rowsData.PhotoWidth,
			&rowsData.PhotoHeight,
			&rowsData.CreatedDate,
			&rowsData.ReservoirUID,
			&rowsData.ReservoirName,
			&rowsData.FarmUID,
			&rowsData.FarmName,
		)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Error: err}
		}

		if errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Result: areaRead}
		}

		areaUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		reservoirUID, err := uuid.FromString(rowsData.ReservoirUID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		farmUID, err := uuid.FromString(rowsData.FarmUID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		areaCreatedDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
		if err != nil {
			result <- query.Result{Error: err}
		}

		rows, err := s.DB.Query("SELECT * FROM AREA_READ_NOTES WHERE AREA_UID = ?", uid)
		if err != nil {
			result <- query.Result{Error: err}
		}

		notes := []storage.AreaNote{}

		for rows.Next() {
			err := rows.Scan(
				&notesRowsData.UID,
				&notesRowsData.AreaUID,
				&notesRowsData.Content,
				&notesRowsData.CreatedDate,
			)
			if err != nil {
				result <- query.Result{Error: err}
			}

			noteUID, err := uuid.FromString(notesRowsData.UID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			noteCreatedDate, err := time.Parse(time.RFC3339, notesRowsData.CreatedDate)
			if err != nil {
				result <- query.Result{Error: err}
			}

			notes = append(notes, storage.AreaNote{
				UID:         noteUID,
				Content:     notesRowsData.Content,
				CreatedDate: noteCreatedDate,
			})
		}

		sizeUnit := domain.GetAreaUnit(rowsData.SizeUnit)
		if sizeUnit == (domain.AreaUnit{}) {
			result <- query.Result{Error: domain.AreaError{Code: domain.AreaErrorInvalidSizeUnitCode}}
		}

		location := domain.GetAreaLocation(rowsData.Location)
		if location == (domain.AreaLocation{}) {
			result <- query.Result{Error: domain.AreaError{Code: domain.AreaErrorInvalidAreaLocationCode}}
		}

		areaRead = storage.AreaRead{
			UID:  areaUID,
			Name: rowsData.Name,
			Size: storage.AreaSize{
				Value: rowsData.Size,
				Unit:  sizeUnit,
			},
			Location: storage.AreaLocation(location),
			Type:     rowsData.Type,
			Photo: storage.AreaPhoto{
				Filename: rowsData.PhotoFilename,
				MimeType: rowsData.PhotoMimetype,
				Size:     rowsData.PhotoSize,
				Width:    rowsData.PhotoWidth,
				Height:   rowsData.PhotoHeight,
			},
			CreatedDate: areaCreatedDate,
			Notes:       notes,
			Farm: storage.AreaFarm{
				UID:  farmUID,
				Name: rowsData.FarmName,
			},
			Reservoir: storage.AreaReservoir{
				UID:  reservoirUID,
				Name: rowsData.ReservoirName,
			},
		}

		result <- query.Result{Result: areaRead}
		close(result)
	}()

	return result
}

func (s AreaReadQuerySqlite) FindAllByFarm(farmUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		areaReads := []storage.AreaRead{}

		rows, err := s.DB.Query("SELECT * FROM AREA_READ WHERE FARM_UID = ?", farmUID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		for rows.Next() {
			rowsData := areaReadResult{}
			if err = rows.Scan(
				&rowsData.UID,
				&rowsData.Name,
				&rowsData.SizeUnit,
				&rowsData.Size,
				&rowsData.Type,
				&rowsData.Location,
				&rowsData.PhotoFilename,
				&rowsData.PhotoMimetype,
				&rowsData.PhotoSize,
				&rowsData.PhotoWidth,
				&rowsData.PhotoHeight,
				&rowsData.CreatedDate,
				&rowsData.ReservoirUID,
				&rowsData.ReservoirName,
				&rowsData.FarmUID,
				&rowsData.FarmName,
			); err != nil {
				result <- query.Result{Error: err}
			}

			areaUID, err := uuid.FromString(rowsData.UID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			reservoirUID, err := uuid.FromString(rowsData.ReservoirUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			farmUID, err := uuid.FromString(rowsData.FarmUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			areaCreatedDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.Result{Error: err}
			}

			rows, err := s.DB.Query("SELECT * FROM AREA_READ_NOTES WHERE AREA_UID = ?", areaUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			notes := []storage.AreaNote{}

			for rows.Next() {
				notesRowsData := areaNotesReadResult{}
				if err = rows.Scan(
					&notesRowsData.UID,
					&notesRowsData.AreaUID,
					&notesRowsData.Content,
					&notesRowsData.CreatedDate,
				); err != nil {
					result <- query.Result{Error: err}
				}

				noteUID, err := uuid.FromString(notesRowsData.UID)
				if err != nil {
					result <- query.Result{Error: err}
				}

				noteCreatedDate, err := time.Parse(time.RFC3339, notesRowsData.CreatedDate)
				if err != nil {
					result <- query.Result{Error: err}
				}

				notes = append(notes, storage.AreaNote{
					UID:         noteUID,
					Content:     notesRowsData.Content,
					CreatedDate: noteCreatedDate,
				})
			}

			sizeUnit := domain.GetAreaUnit(rowsData.SizeUnit)
			if sizeUnit == (domain.AreaUnit{}) {
				result <- query.Result{Error: domain.AreaError{Code: domain.AreaErrorInvalidSizeUnitCode}}
			}

			location := domain.GetAreaLocation(rowsData.Location)
			if location == (domain.AreaLocation{}) {
				result <- query.Result{Error: domain.AreaError{Code: domain.AreaErrorInvalidAreaLocationCode}}
			}

			areaReads = append(areaReads, storage.AreaRead{
				UID:  areaUID,
				Name: rowsData.Name,
				Size: storage.AreaSize{
					Value: rowsData.Size,
					Unit:  sizeUnit,
				},
				Location: storage.AreaLocation(location),
				Type:     rowsData.Type,
				Photo: storage.AreaPhoto{
					Filename: rowsData.PhotoFilename,
					MimeType: rowsData.PhotoMimetype,
					Size:     rowsData.PhotoSize,
					Width:    rowsData.PhotoWidth,
					Height:   rowsData.PhotoHeight,
				},
				CreatedDate: areaCreatedDate,
				Notes:       notes,
				Farm: storage.AreaFarm{
					UID:  farmUID,
					Name: rowsData.FarmName,
				},
				Reservoir: storage.AreaReservoir{
					UID:  reservoirUID,
					Name: rowsData.ReservoirName,
				},
			})
		}

		result <- query.Result{Result: areaReads}
		close(result)
	}()

	return result
}

func (s AreaReadQuerySqlite) FindByIDAndFarm(areaUID, farmUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		areaRead := storage.AreaRead{}
		rowsData := areaReadResult{}
		notesRowsData := areaNotesReadResult{}

		err := s.DB.QueryRow("SELECT * FROM AREA_READ WHERE UID = ? AND FARM_UID = ?", areaUID, farmUID).Scan(
			&rowsData.UID,
			&rowsData.Name,
			&rowsData.SizeUnit,
			&rowsData.Size,
			&rowsData.Type,
			&rowsData.Location,
			&rowsData.PhotoFilename,
			&rowsData.PhotoMimetype,
			&rowsData.PhotoSize,
			&rowsData.PhotoWidth,
			&rowsData.PhotoHeight,
			&rowsData.CreatedDate,
			&rowsData.ReservoirUID,
			&rowsData.ReservoirName,
			&rowsData.FarmUID,
			&rowsData.FarmName,
		)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Error: err}
		}

		if errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Result: areaRead}
		}

		areaUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		reservoirUID, err := uuid.FromString(rowsData.ReservoirUID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		farmUID, err := uuid.FromString(rowsData.FarmUID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		areaCreatedDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
		if err != nil {
			result <- query.Result{Error: err}
		}

		rows, err := s.DB.Query("SELECT * FROM AREA_READ_NOTES WHERE AREA_UID = ?", areaUID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		notes := []storage.AreaNote{}

		for rows.Next() {
			err = rows.Scan(
				&notesRowsData.UID,
				&notesRowsData.AreaUID,
				&notesRowsData.Content,
				&notesRowsData.CreatedDate,
			)
			if err != nil {
				result <- query.Result{Error: err}
			}

			noteUID, err := uuid.FromString(notesRowsData.UID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			noteCreatedDate, err := time.Parse(time.RFC3339, notesRowsData.CreatedDate)
			if err != nil {
				result <- query.Result{Error: err}
			}

			notes = append(notes, storage.AreaNote{
				UID:         noteUID,
				Content:     notesRowsData.Content,
				CreatedDate: noteCreatedDate,
			})
		}

		sizeUnit := domain.GetAreaUnit(rowsData.SizeUnit)
		if sizeUnit == (domain.AreaUnit{}) {
			result <- query.Result{Error: domain.AreaError{Code: domain.AreaErrorInvalidSizeUnitCode}}
		}

		location := domain.GetAreaLocation(rowsData.Location)
		if location == (domain.AreaLocation{}) {
			result <- query.Result{Error: domain.AreaError{Code: domain.AreaErrorInvalidAreaLocationCode}}
		}

		areaRead = storage.AreaRead{
			UID:  areaUID,
			Name: rowsData.Name,
			Size: storage.AreaSize{
				Value: rowsData.Size,
				Unit:  sizeUnit,
			},
			Location: storage.AreaLocation(location),
			Type:     rowsData.Type,
			Photo: storage.AreaPhoto{
				Filename: rowsData.PhotoFilename,
				MimeType: rowsData.PhotoMimetype,
				Size:     rowsData.PhotoSize,
				Width:    rowsData.PhotoWidth,
				Height:   rowsData.PhotoHeight,
			},
			CreatedDate: areaCreatedDate,
			Notes:       notes,
			Farm: storage.AreaFarm{
				UID:  farmUID,
				Name: rowsData.FarmName,
			},
			Reservoir: storage.AreaReservoir{
				UID:  reservoirUID,
				Name: rowsData.ReservoirName,
			},
		}

		result <- query.Result{Result: areaRead}
		close(result)
	}()

	return result
}

func (s AreaReadQuerySqlite) FindAreasByReservoirID(reservoirUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		areaReads := []storage.AreaRead{}

		rows, err := s.DB.Query("SELECT * FROM AREA_READ WHERE RESERVOIR_UID = ?", reservoirUID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		for rows.Next() {
			rowsData := areaReadResult{}
			if err = rows.Scan(
				&rowsData.UID,
				&rowsData.Name,
				&rowsData.SizeUnit,
				&rowsData.Size,
				&rowsData.Type,
				&rowsData.Location,
				&rowsData.PhotoFilename,
				&rowsData.PhotoMimetype,
				&rowsData.PhotoSize,
				&rowsData.PhotoWidth,
				&rowsData.PhotoHeight,
				&rowsData.CreatedDate,
				&rowsData.ReservoirUID,
				&rowsData.ReservoirName,
				&rowsData.FarmUID,
				&rowsData.FarmName,
			); err != nil {
				result <- query.Result{Error: err}
			}

			areaUID, err := uuid.FromString(rowsData.UID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			reservoirUID, err := uuid.FromString(rowsData.ReservoirUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			farmUID, err := uuid.FromString(rowsData.FarmUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			areaCreatedDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.Result{Error: err}
			}

			rows, err := s.DB.Query("SELECT * FROM AREA_READ_NOTES WHERE AREA_UID = ?", areaUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			notes := []storage.AreaNote{}

			for rows.Next() {
				notesRowsData := areaNotesReadResult{}
				if err = rows.Scan(
					&notesRowsData.UID,
					&notesRowsData.AreaUID,
					&notesRowsData.Content,
					&notesRowsData.CreatedDate,
				); err != nil {
					result <- query.Result{Error: err}
				}

				noteUID, err := uuid.FromString(notesRowsData.UID)
				if err != nil {
					result <- query.Result{Error: err}
				}

				noteCreatedDate, err := time.Parse(time.RFC3339, notesRowsData.CreatedDate)
				if err != nil {
					result <- query.Result{Error: err}
				}

				notes = append(notes, storage.AreaNote{
					UID:         noteUID,
					Content:     notesRowsData.Content,
					CreatedDate: noteCreatedDate,
				})
			}

			sizeUnit := domain.GetAreaUnit(rowsData.SizeUnit)
			if sizeUnit == (domain.AreaUnit{}) {
				result <- query.Result{Error: domain.AreaError{Code: domain.AreaErrorInvalidSizeUnitCode}}
			}

			location := domain.GetAreaLocation(rowsData.Location)
			if location == (domain.AreaLocation{}) {
				result <- query.Result{Error: domain.AreaError{Code: domain.AreaErrorInvalidAreaLocationCode}}
			}

			areaReads = append(areaReads, storage.AreaRead{
				UID:  areaUID,
				Name: rowsData.Name,
				Size: storage.AreaSize{
					Value: rowsData.Size,
					Unit:  sizeUnit,
				},
				Location: storage.AreaLocation(location),
				Type:     rowsData.Type,
				Photo: storage.AreaPhoto{
					Filename: rowsData.PhotoFilename,
					MimeType: rowsData.PhotoMimetype,
					Size:     rowsData.PhotoSize,
					Width:    rowsData.PhotoWidth,
					Height:   rowsData.PhotoHeight,
				},
				CreatedDate: areaCreatedDate,
				Notes:       notes,
				Farm: storage.AreaFarm{
					UID:  farmUID,
					Name: rowsData.FarmName,
				},
				Reservoir: storage.AreaReservoir{
					UID:  reservoirUID,
					Name: rowsData.ReservoirName,
				},
			})
		}

		result <- query.Result{Result: areaReads}
		close(result)
	}()

	return result
}

func (s AreaReadQuerySqlite) CountAreas(farmUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		total := 0

		err := s.DB.QueryRow(`SELECT COUNT(*) FROM AREA_READ WHERE FARM_UID = ?`, farmUID).Scan(&total)
		if err != nil {
			result <- query.Result{Error: err}
		}

		result <- query.Result{Result: total}

		close(result)
	}()

	return result
}
