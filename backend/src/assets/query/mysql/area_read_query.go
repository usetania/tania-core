package mysql

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/domain"
	"github.com/usetania/tania-core/src/assets/query"
	"github.com/usetania/tania-core/src/assets/storage"
)

type AreaReadQueryMysql struct {
	DB *sql.DB
}

func NewAreaReadQueryMysql(db *sql.DB) query.AreaRead {
	return AreaReadQueryMysql{DB: db}
}

type areaReadResult struct {
	UID           []byte
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
	CreatedDate   time.Time
	ReservoirUID  []byte
	ReservoirName string
	FarmUID       []byte
	FarmName      string
}

type areaNotesReadResult struct {
	UID         []byte
	AreaUID     []byte
	Content     string
	CreatedDate time.Time
}

func (s AreaReadQueryMysql) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		areaRead := storage.AreaRead{}
		rowsData := areaReadResult{}
		notesRowsData := areaNotesReadResult{}

		err := s.DB.QueryRow("SELECT * FROM AREA_READ WHERE UID = ?", uid.Bytes()).Scan(
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

		areaUID, err := uuid.FromBytes(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		reservoirUID, err := uuid.FromBytes(rowsData.ReservoirUID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		farmUID, err := uuid.FromBytes(rowsData.FarmUID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		rows, err := s.DB.Query("SELECT * FROM AREA_READ_NOTES WHERE AREA_UID = ?", uid.Bytes())
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

			noteUID, err := uuid.FromBytes(notesRowsData.UID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			notes = append(notes, storage.AreaNote{
				UID:         noteUID,
				Content:     notesRowsData.Content,
				CreatedDate: notesRowsData.CreatedDate,
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
			CreatedDate: rowsData.CreatedDate,
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

func (s AreaReadQueryMysql) FindAllByFarm(farmUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		areaReads := []storage.AreaRead{}

		rows, err := s.DB.Query("SELECT * FROM AREA_READ WHERE FARM_UID = ?", farmUID.Bytes())
		if err != nil {
			result <- query.Result{Error: err}
		}

		for rows.Next() {
			rowsData := areaReadResult{}
			if err := rows.Scan(
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

			areaUID, err := uuid.FromBytes(rowsData.UID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			reservoirUID, err := uuid.FromBytes(rowsData.ReservoirUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			farmUID, err := uuid.FromBytes(rowsData.FarmUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			rows, err := s.DB.Query("SELECT * FROM AREA_READ_NOTES WHERE AREA_UID = ?", areaUID.Bytes())
			if err != nil {
				result <- query.Result{Error: err}
			}

			notes := []storage.AreaNote{}

			for rows.Next() {
				notesRowsData := areaNotesReadResult{}
				if err := rows.Scan(
					&notesRowsData.UID,
					&notesRowsData.AreaUID,
					&notesRowsData.Content,
					&notesRowsData.CreatedDate,
				); err != nil {
					result <- query.Result{Error: err}
				}

				noteUID, err := uuid.FromBytes(notesRowsData.UID)
				if err != nil {
					result <- query.Result{Error: err}
				}

				notes = append(notes, storage.AreaNote{
					UID:         noteUID,
					Content:     notesRowsData.Content,
					CreatedDate: notesRowsData.CreatedDate,
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
				CreatedDate: rowsData.CreatedDate,
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

func (s AreaReadQueryMysql) FindByIDAndFarm(areaUID, farmUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		areaRead := storage.AreaRead{}
		rowsData := areaReadResult{}
		notesRowsData := areaNotesReadResult{}

		err := s.DB.QueryRow("SELECT * FROM AREA_READ WHERE UID = ? AND FARM_UID = ?", areaUID.Bytes(), farmUID.Bytes()).Scan(
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

		areaUID, err := uuid.FromBytes(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		reservoirUID, err := uuid.FromBytes(rowsData.ReservoirUID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		farmUID, err := uuid.FromBytes(rowsData.FarmUID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		rows, err := s.DB.Query("SELECT * FROM AREA_READ_NOTES WHERE AREA_UID = ?", areaUID.Bytes())
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

			noteUID, err := uuid.FromBytes(notesRowsData.UID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			notes = append(notes, storage.AreaNote{
				UID:         noteUID,
				Content:     notesRowsData.Content,
				CreatedDate: notesRowsData.CreatedDate,
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
			CreatedDate: rowsData.CreatedDate,
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

func (s AreaReadQueryMysql) FindAreasByReservoirID(reservoirUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		areaReads := []storage.AreaRead{}

		rows, err := s.DB.Query("SELECT * FROM AREA_READ WHERE RESERVOIR_UID = ?", reservoirUID.Bytes())
		if err != nil {
			result <- query.Result{Error: err}
		}

		for rows.Next() {
			rowsData := areaReadResult{}
			if err := rows.Scan(
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

			areaUID, err := uuid.FromBytes(rowsData.UID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			reservoirUID, err := uuid.FromBytes(rowsData.ReservoirUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			farmUID, err := uuid.FromBytes(rowsData.FarmUID)
			if err != nil {
				result <- query.Result{Error: err}
			}

			rows, err := s.DB.Query("SELECT * FROM AREA_READ_NOTES WHERE AREA_UID = ?", areaUID.Bytes())
			if err != nil {
				result <- query.Result{Error: err}
			}

			notes := []storage.AreaNote{}

			for rows.Next() {
				notesRowsData := areaNotesReadResult{}
				if err := rows.Scan(
					&notesRowsData.UID,
					&notesRowsData.AreaUID,
					&notesRowsData.Content,
					&notesRowsData.CreatedDate,
				); err != nil {
					result <- query.Result{Error: err}
				}

				noteUID, err := uuid.FromBytes(notesRowsData.UID)
				if err != nil {
					result <- query.Result{Error: err}
				}

				notes = append(notes, storage.AreaNote{
					UID:         noteUID,
					Content:     notesRowsData.Content,
					CreatedDate: notesRowsData.CreatedDate,
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
				CreatedDate: rowsData.CreatedDate,
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

func (s AreaReadQueryMysql) CountAreas(farmUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		total := 0

		err := s.DB.QueryRow(`SELECT COUNT(*) FROM AREA_READ WHERE FARM_UID = ?`, farmUID.Bytes()).Scan(&total)
		if err != nil {
			result <- query.Result{Error: err}
		}

		result <- query.Result{Result: total}

		close(result)
	}()

	return result
}
