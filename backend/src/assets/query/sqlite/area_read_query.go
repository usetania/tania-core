package sqlite

import (
	"database/sql"
	"time"

	"github.com/Tanibox/tania-core/src/assets/domain"
	"github.com/Tanibox/tania-core/src/assets/query"
	"github.com/Tanibox/tania-core/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type AreaReadQuerySqlite struct {
	DB *sql.DB
}

func NewAreaReadQuerySqlite(db *sql.DB) query.AreaReadQuery {
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

func (s AreaReadQuerySqlite) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

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

		if err != nil && err != sql.ErrNoRows {
			result <- query.QueryResult{Error: err}
		}

		if err == sql.ErrNoRows {
			result <- query.QueryResult{Result: areaRead}
		}

		areaUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		reservoirUID, err := uuid.FromString(rowsData.ReservoirUID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		farmUID, err := uuid.FromString(rowsData.FarmUID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		areaCreatedDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		rows, err := s.DB.Query("SELECT * FROM AREA_READ_NOTES WHERE AREA_UID = ?", uid)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		notes := []storage.AreaNote{}
		for rows.Next() {
			rows.Scan(
				&notesRowsData.UID,
				&notesRowsData.AreaUID,
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

			notes = append(notes, storage.AreaNote{
				UID:         noteUID,
				Content:     notesRowsData.Content,
				CreatedDate: noteCreatedDate,
			})
		}

		sizeUnit := domain.GetAreaUnit(rowsData.SizeUnit)
		if sizeUnit == (domain.AreaUnit{}) {
			result <- query.QueryResult{Error: domain.AreaError{domain.AreaErrorInvalidSizeUnitCode}}
		}

		location := domain.GetAreaLocation(rowsData.Location)
		if location == (domain.AreaLocation{}) {
			result <- query.QueryResult{Error: domain.AreaError{domain.AreaErrorInvalidAreaLocationCode}}
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

		result <- query.QueryResult{Result: areaRead}
		close(result)
	}()

	return result
}

func (s AreaReadQuerySqlite) FindAllByFarm(farmUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		areaReads := []storage.AreaRead{}

		rows, err := s.DB.Query("SELECT * FROM AREA_READ WHERE FARM_UID = ?", farmUID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		for rows.Next() {
			rowsData := areaReadResult{}
			rows.Scan(
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

			areaUID, err := uuid.FromString(rowsData.UID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			reservoirUID, err := uuid.FromString(rowsData.ReservoirUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			farmUID, err := uuid.FromString(rowsData.FarmUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			areaCreatedDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			rows, err := s.DB.Query("SELECT * FROM AREA_READ_NOTES WHERE AREA_UID = ?", areaUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			notes := []storage.AreaNote{}
			for rows.Next() {
				notesRowsData := areaNotesReadResult{}
				rows.Scan(
					&notesRowsData.UID,
					&notesRowsData.AreaUID,
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

				notes = append(notes, storage.AreaNote{
					UID:         noteUID,
					Content:     notesRowsData.Content,
					CreatedDate: noteCreatedDate,
				})
			}

			sizeUnit := domain.GetAreaUnit(rowsData.SizeUnit)
			if sizeUnit == (domain.AreaUnit{}) {
				result <- query.QueryResult{Error: domain.AreaError{domain.AreaErrorInvalidSizeUnitCode}}
			}

			location := domain.GetAreaLocation(rowsData.Location)
			if location == (domain.AreaLocation{}) {
				result <- query.QueryResult{Error: domain.AreaError{domain.AreaErrorInvalidAreaLocationCode}}
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

		result <- query.QueryResult{Result: areaReads}
		close(result)
	}()

	return result
}

func (s AreaReadQuerySqlite) FindByIDAndFarm(areaUID, farmUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

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

		if err != nil && err != sql.ErrNoRows {
			result <- query.QueryResult{Error: err}
		}

		if err == sql.ErrNoRows {
			result <- query.QueryResult{Result: areaRead}
		}

		areaUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		reservoirUID, err := uuid.FromString(rowsData.ReservoirUID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		farmUID, err := uuid.FromString(rowsData.FarmUID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		areaCreatedDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		rows, err := s.DB.Query("SELECT * FROM AREA_READ_NOTES WHERE AREA_UID = ?", areaUID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		notes := []storage.AreaNote{}
		for rows.Next() {
			rows.Scan(
				&notesRowsData.UID,
				&notesRowsData.AreaUID,
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

			notes = append(notes, storage.AreaNote{
				UID:         noteUID,
				Content:     notesRowsData.Content,
				CreatedDate: noteCreatedDate,
			})
		}

		sizeUnit := domain.GetAreaUnit(rowsData.SizeUnit)
		if sizeUnit == (domain.AreaUnit{}) {
			result <- query.QueryResult{Error: domain.AreaError{domain.AreaErrorInvalidSizeUnitCode}}
		}

		location := domain.GetAreaLocation(rowsData.Location)
		if location == (domain.AreaLocation{}) {
			result <- query.QueryResult{Error: domain.AreaError{domain.AreaErrorInvalidAreaLocationCode}}
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

		result <- query.QueryResult{Result: areaRead}
		close(result)
	}()

	return result
}

func (s AreaReadQuerySqlite) FindAreasByReservoirID(reservoirUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		areaReads := []storage.AreaRead{}

		rows, err := s.DB.Query("SELECT * FROM AREA_READ WHERE RESERVOIR_UID = ?", reservoirUID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		for rows.Next() {
			rowsData := areaReadResult{}
			rows.Scan(
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

			areaUID, err := uuid.FromString(rowsData.UID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			reservoirUID, err := uuid.FromString(rowsData.ReservoirUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			farmUID, err := uuid.FromString(rowsData.FarmUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			areaCreatedDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			rows, err := s.DB.Query("SELECT * FROM AREA_READ_NOTES WHERE AREA_UID = ?", areaUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			notes := []storage.AreaNote{}
			for rows.Next() {
				notesRowsData := areaNotesReadResult{}
				rows.Scan(
					&notesRowsData.UID,
					&notesRowsData.AreaUID,
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

				notes = append(notes, storage.AreaNote{
					UID:         noteUID,
					Content:     notesRowsData.Content,
					CreatedDate: noteCreatedDate,
				})
			}

			sizeUnit := domain.GetAreaUnit(rowsData.SizeUnit)
			if sizeUnit == (domain.AreaUnit{}) {
				result <- query.QueryResult{Error: domain.AreaError{domain.AreaErrorInvalidSizeUnitCode}}
			}

			location := domain.GetAreaLocation(rowsData.Location)
			if location == (domain.AreaLocation{}) {
				result <- query.QueryResult{Error: domain.AreaError{domain.AreaErrorInvalidAreaLocationCode}}
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

		result <- query.QueryResult{Result: areaReads}
		close(result)
	}()

	return result
}

func (s AreaReadQuerySqlite) CountAreas(farmUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		total := 0
		err := s.DB.QueryRow(`SELECT COUNT(*) FROM AREA_READ WHERE FARM_UID = ?`, farmUID).Scan(&total)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		result <- query.QueryResult{Result: total}

		close(result)
	}()

	return result
}
