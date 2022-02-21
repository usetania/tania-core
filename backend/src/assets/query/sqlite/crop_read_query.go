package sqlite

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/query"
	"github.com/usetania/tania-core/src/assets/storage"
)

type CropReadQuerySqlite struct {
	DB *sql.DB
}

func NewCropReadQuerySqlite(db *sql.DB) query.CropRead {
	return CropReadQuerySqlite{DB: db}
}

type cropReadResult struct {
	UID                        string
	BatchID                    string
	Status                     string
	Type                       string
	ContainerQuantity          int
	ContainerType              string
	ContainerCell              int
	InventoryUID               string
	InventoryPlantType         string
	InventoryName              string
	AreaStatusSeeding          int
	AreaStatusGrowing          int
	AreaStatusDumped           int
	FarmUID                    string
	InitialAreaUID             string
	InitialAreaName            string
	InitialAreaInitialQuantity int
	InitialAreaCurrentQuantity int
	InitialAreaLastWatered     sql.NullString
	InitialAreaLastFertilized  sql.NullString
	InitialAreaLastPesticided  sql.NullString
	InitialAreaLastPruned      sql.NullString
	InitialAreaCreatedDate     string
	InitialAreaLastUpdated     string
}

type cropReadMovedAreaResult struct {
	ID              int
	CropUID         string
	AreaUID         string
	Name            string
	InitialQuantity int
	CurrentQuantity int
	LastWatered     sql.NullString
	LastFertilized  sql.NullString
	LastPesticided  sql.NullString
	LastPruned      sql.NullString
	CreatedDate     string
	LastUpdated     string
}

func (q CropReadQuerySqlite) CountCropsByArea(areaUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		var totalCropBatchInitial, totalPlantInitial sql.NullInt64

		err := q.DB.QueryRow(`SELECT COUNT(UID), SUM(INITIAL_AREA_CURRENT_QUANTITY)
			FROM CROP_READ WHERE INITIAL_AREA_UID = ?`, areaUID).Scan(&totalCropBatchInitial, &totalPlantInitial)
		if err != nil {
			result <- query.Result{Error: err}
		}

		var totalCropBatchMoved, totalPlantMoved sql.NullInt64
		err = q.DB.QueryRow(`SELECT COUNT(CROP_UID), SUM(CURRENT_QUANTITY)
			FROM CROP_READ_MOVED_AREA WHERE AREA_UID = ?`, areaUID).Scan(&totalCropBatchMoved, &totalPlantMoved)

		if err != nil {
			result <- query.Result{Error: err}
		}

		result <- query.Result{Result: query.CountAreaCropResult{
			PlantQuantity:  int(totalPlantInitial.Int64) + int(totalPlantMoved.Int64),
			TotalCropBatch: int(totalCropBatchInitial.Int64) + int(totalCropBatchMoved.Int64),
		}}

		close(result)
	}()

	return result
}

func (q CropReadQuerySqlite) FindAllCropByArea(areaUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		crops := []query.AreaCropResult{}

		// TODO: REFACTOR TO REDUCE QUERY CALLS
		rows, err := q.DB.Query("SELECT UID FROM CROP_READ WHERE INITIAL_AREA_UID = ?", areaUID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		for rows.Next() {
			cropRead := storage.CropRead{}

			uid := ""

			err := rows.Scan(&uid)
			if err != nil {
				result <- query.Result{Error: err}
			}

			cropUID, err := uuid.FromString(uid)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = q.populateCrop(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			crops = append(crops, query.AreaCropResult{
				CropUID: cropRead.UID,
				BatchID: cropRead.BatchID,
				InitialArea: query.InitialArea{
					AreaUID: cropRead.InitialArea.AreaUID,
					Name:    cropRead.InitialArea.Name,
				},
				MovingDate:  cropRead.InitialArea.CreatedDate,
				CreatedDate: cropRead.InitialArea.CreatedDate,
				Inventory: query.Inventory{
					UID: cropRead.Inventory.UID,
				},
				Container: query.Container{
					Quantity: cropRead.Container.Quantity,
					Type: query.ContainerType{
						Code: cropRead.Container.Type,
						Cell: cropRead.Container.Cell,
					},
				},
			})
		}

		rows, err = q.DB.Query(`SELECT UID FROM CROP_READ
			LEFT JOIN CROP_READ_MOVED_AREA ON CROP_READ.UID = CROP_READ_MOVED_AREA.CROP_UID
			WHERE CROP_READ_MOVED_AREA.AREA_UID = ?`, areaUID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		for rows.Next() {
			cropRead := storage.CropRead{}

			uid := ""

			err := rows.Scan(&uid)
			if err != nil {
				result <- query.Result{Error: err}
			}

			cropUID, err := uuid.FromString(uid)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = q.populateCrop(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = q.populateCropMovedArea(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			for _, val := range cropRead.MovedArea {
				crops = append(crops, query.AreaCropResult{
					CropUID: cropRead.UID,
					BatchID: cropRead.BatchID,
					InitialArea: query.InitialArea{
						AreaUID: cropRead.InitialArea.AreaUID,
					},
					MovingDate:  val.CreatedDate,
					CreatedDate: val.CreatedDate,
					Inventory: query.Inventory{
						UID: cropRead.Inventory.UID,
					},
					Container: query.Container{
						Quantity: cropRead.Container.Quantity,
						Type: query.ContainerType{
							Code: cropRead.Container.Type,
							Cell: cropRead.Container.Cell,
						},
					},
				})
			}
		}

		result <- query.Result{Result: crops}
		close(result)
	}()

	return result
}

func (q CropReadQuerySqlite) populateCrop(cropUID uuid.UUID, cropRead *storage.CropRead) error {
	rowsData := cropReadResult{}

	err := q.DB.QueryRow(`SELECT UID, BATCH_ID, STATUS, TYPE, CONTAINER_QUANTITY, CONTAINER_TYPE, CONTAINER_CELL,
		INVENTORY_UID, INVENTORY_PLANT_TYPE, INVENTORY_NAME,
		AREA_STATUS_SEEDING, AREA_STATUS_GROWING, AREA_STATUS_DUMPED,
		FARM_UID,
		INITIAL_AREA_UID, INITIAL_AREA_NAME,
		INITIAL_AREA_INITIAL_QUANTITY, INITIAL_AREA_CURRENT_QUANTITY,
		INITIAL_AREA_LAST_WATERED, INITIAL_AREA_LAST_FERTILIZED, INITIAL_AREA_LAST_PESTICIDED,
		INITIAL_AREA_LAST_PRUNED, INITIAL_AREA_CREATED_DATE, INITIAL_AREA_LAST_UPDATED
		FROM CROP_READ WHERE UID = ?`, cropUID).Scan(
		&rowsData.UID,
		&rowsData.BatchID,
		&rowsData.Status,
		&rowsData.Type,
		&rowsData.ContainerQuantity,
		&rowsData.ContainerType,
		&rowsData.ContainerCell,
		&rowsData.InventoryUID,
		&rowsData.InventoryPlantType,
		&rowsData.InventoryName,
		&rowsData.AreaStatusSeeding,
		&rowsData.AreaStatusGrowing,
		&rowsData.AreaStatusDumped,
		&rowsData.FarmUID,
		&rowsData.InitialAreaUID,
		&rowsData.InitialAreaName,
		&rowsData.InitialAreaInitialQuantity,
		&rowsData.InitialAreaCurrentQuantity,
		&rowsData.InitialAreaLastWatered,
		&rowsData.InitialAreaLastFertilized,
		&rowsData.InitialAreaLastPesticided,
		&rowsData.InitialAreaLastPruned,
		&rowsData.InitialAreaCreatedDate,
		&rowsData.InitialAreaLastUpdated,
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return err
	}

	farmUID, err := uuid.FromString(rowsData.FarmUID)
	if err != nil {
		return err
	}

	inventoryUID, err := uuid.FromString(rowsData.InventoryUID)
	if err != nil {
		return err
	}

	initialAreaUID, err := uuid.FromString(rowsData.InitialAreaUID)
	if err != nil {
		return err
	}

	var initialAreaLastWatered *time.Time

	if rowsData.InitialAreaLastWatered.Valid && rowsData.InitialAreaLastWatered.String != "" {
		date, err := time.Parse(time.RFC3339, rowsData.InitialAreaLastWatered.String)
		if err != nil {
			return err
		}

		initialAreaLastWatered = &date
	}

	var initialAreaLastFertilized *time.Time

	if rowsData.InitialAreaLastFertilized.Valid && rowsData.InitialAreaLastFertilized.String != "" {
		date, err := time.Parse(time.RFC3339, rowsData.InitialAreaLastFertilized.String)
		if err != nil {
			return err
		}

		initialAreaLastFertilized = &date
	}

	var initialAreaLastPesticided *time.Time

	if rowsData.InitialAreaLastPesticided.Valid && rowsData.InitialAreaLastPesticided.String != "" {
		date, err := time.Parse(time.RFC3339, rowsData.InitialAreaLastPesticided.String)
		if err != nil {
			return err
		}

		initialAreaLastPesticided = &date
	}

	var initialAreaLastPruned *time.Time

	if rowsData.InitialAreaLastPruned.Valid && rowsData.InitialAreaLastPruned.String != "" {
		date, err := time.Parse(time.RFC3339, rowsData.InitialAreaLastPruned.String)
		if err != nil {
			return err
		}

		initialAreaLastPruned = &date
	}

	initialAreaCreatedDate, err := time.Parse(time.RFC3339, rowsData.InitialAreaCreatedDate)
	if err != nil {
		return err
	}

	initialAreaLastUpdated, err := time.Parse(time.RFC3339, rowsData.InitialAreaLastUpdated)
	if err != nil {
		return err
	}

	cropRead.UID = cropUID
	cropRead.BatchID = rowsData.BatchID
	cropRead.Status = rowsData.Status
	cropRead.Type = rowsData.Type
	cropRead.Container.Quantity = rowsData.ContainerQuantity
	cropRead.Container.Type = rowsData.ContainerType
	cropRead.Container.Cell = rowsData.ContainerCell
	cropRead.Inventory.UID = inventoryUID
	cropRead.Inventory.PlantType = rowsData.InventoryPlantType
	cropRead.Inventory.Name = rowsData.InventoryName
	cropRead.AreaStatus.Seeding = rowsData.AreaStatusSeeding
	cropRead.AreaStatus.Growing = rowsData.AreaStatusGrowing
	cropRead.AreaStatus.Dumped = rowsData.AreaStatusDumped
	cropRead.FarmUID = farmUID
	cropRead.InitialArea.AreaUID = initialAreaUID
	cropRead.InitialArea.Name = rowsData.InitialAreaName
	cropRead.InitialArea.InitialQuantity = rowsData.InitialAreaInitialQuantity
	cropRead.InitialArea.CurrentQuantity = rowsData.InitialAreaCurrentQuantity
	cropRead.InitialArea.LastWatered = initialAreaLastWatered
	cropRead.InitialArea.LastFertilized = initialAreaLastFertilized
	cropRead.InitialArea.LastPesticided = initialAreaLastPesticided
	cropRead.InitialArea.LastPruned = initialAreaLastPruned
	cropRead.InitialArea.CreatedDate = initialAreaCreatedDate
	cropRead.InitialArea.LastUpdated = initialAreaLastUpdated

	return nil
}

func (q CropReadQuerySqlite) populateCropMovedArea(uid uuid.UUID, cropRead *storage.CropRead) error {
	movedRowsData := cropReadMovedAreaResult{}

	rows, err := q.DB.Query("SELECT * FROM CROP_READ_MOVED_AREA WHERE CROP_UID = ?", uid)
	if err != nil {
		return err
	}

	movedAreas := []storage.MovedArea{}

	for rows.Next() {
		err = rows.Scan(
			&movedRowsData.ID,
			&movedRowsData.CropUID,
			&movedRowsData.AreaUID,
			&movedRowsData.Name,
			&movedRowsData.InitialQuantity,
			&movedRowsData.CurrentQuantity,
			&movedRowsData.LastWatered,
			&movedRowsData.LastFertilized,
			&movedRowsData.LastPesticided,
			&movedRowsData.LastPruned,
			&movedRowsData.CreatedDate,
			&movedRowsData.LastUpdated,
		)

		if err != nil {
			return err
		}

		var lw *time.Time

		if movedRowsData.LastWatered.Valid && movedRowsData.LastWatered.String != "" {
			date, err := time.Parse(time.RFC3339, movedRowsData.LastWatered.String)
			if err != nil {
				return err
			}

			lw = &date
		}

		var lf *time.Time

		if movedRowsData.LastFertilized.Valid && movedRowsData.LastFertilized.String != "" {
			date, err := time.Parse(time.RFC3339, movedRowsData.LastFertilized.String)
			if err != nil {
				return err
			}

			lf = &date
		}

		var lp *time.Time

		if movedRowsData.LastPesticided.Valid && movedRowsData.LastPesticided.String != "" {
			date, err := time.Parse(time.RFC3339, movedRowsData.LastPesticided.String)
			if err != nil {
				return err
			}

			lp = &date
		}

		var lpr *time.Time

		if movedRowsData.LastPruned.Valid && movedRowsData.LastPruned.String != "" {
			date, err := time.Parse(time.RFC3339, movedRowsData.LastPruned.String)
			if err != nil {
				return err
			}

			lpr = &date
		}

		areaUID, err := uuid.FromString(movedRowsData.AreaUID)
		if err != nil {
			return err
		}

		createdDate, err := time.Parse(time.RFC3339, movedRowsData.CreatedDate)
		if err != nil {
			return err
		}

		lastUpdated, err := time.Parse(time.RFC3339, movedRowsData.LastUpdated)
		if err != nil {
			return err
		}

		movedAreas = append(movedAreas, storage.MovedArea{
			AreaUID:         areaUID,
			Name:            movedRowsData.Name,
			InitialQuantity: movedRowsData.InitialQuantity,
			CurrentQuantity: movedRowsData.CurrentQuantity,
			LastWatered:     lw,
			LastFertilized:  lf,
			LastPesticided:  lp,
			LastPruned:      lpr,
			CreatedDate:     createdDate,
			LastUpdated:     lastUpdated,
		})
	}

	cropRead.MovedArea = movedAreas

	return nil
}
