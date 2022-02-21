package sqlite

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/growth/domain"
	"github.com/usetania/tania-core/src/growth/query"
	"github.com/usetania/tania-core/src/growth/storage"
	"github.com/usetania/tania-core/src/helper/paginationhelper"
)

type CropReadQueryMysql struct {
	DB *sql.DB
}

func NewCropReadQueryMysql(db *sql.DB) query.CropReadQuery {
	return CropReadQueryMysql{DB: db}
}

type cropReadResult struct {
	UID                        []byte
	BatchID                    string
	Status                     string
	Type                       string
	ContainerQuantity          int
	ContainerType              string
	ContainerCell              int
	InventoryUID               []byte
	InventoryType              string
	InventoryPlantType         string
	InventoryName              string
	AreaStatusSeeding          int
	AreaStatusGrowing          int
	AreaStatusDumped           int
	FarmUID                    []byte
	InitialAreaUID             []byte
	InitialAreaName            string
	InitialAreaInitialQuantity int
	InitialAreaCurrentQuantity int
	InitialAreaLastWatered     sql.NullString
	InitialAreaLastFertilized  sql.NullString
	InitialAreaLastPesticided  sql.NullString
	InitialAreaLastPruned      sql.NullString
	InitialAreaCreatedDate     time.Time
	InitialAreaLastUpdated     time.Time
}

type cropReadPhotoResult struct {
	UID         []byte
	CropUID     []byte
	Filename    string
	Mimetype    string
	Size        int
	Width       int
	Height      int
	Description string
}

type cropReadMovedAreaResult struct {
	ID              int
	CropUID         []byte
	AreaUID         []byte
	Name            string
	InitialQuantity int
	CurrentQuantity int
	LastWatered     sql.NullString
	LastFertilized  sql.NullString
	LastPesticided  sql.NullString
	LastPruned      sql.NullString
	CreatedDate     time.Time
	LastUpdated     time.Time
}

type cropReadHarvestedStorageResult struct {
	ID                   int
	CropUID              []byte
	Quantity             int
	ProducedGramQuantity float32
	SourceAreaUID        []byte
	SourceAreaName       string
	CreatedDate          time.Time
	LastUpdated          time.Time
}

type cropReadTrashResult struct {
	ID             int
	CropUID        []byte
	Quantity       int
	SourceAreaUID  []byte
	SourceAreaName string
	CreatedDate    time.Time
	LastUpdated    time.Time
}

type cropReadNotesResult struct {
	UID         []byte
	CropUID     []byte
	Content     string
	CreatedDate time.Time
}

func (s CropReadQueryMysql) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		cropRead := storage.CropRead{}

		err := s.populateCrop(uid, &cropRead)
		if err != nil {
			result <- query.Result{Error: err}
		}

		err = s.populateCropPhotos(uid, &cropRead)
		if err != nil {
			result <- query.Result{Error: err}
		}

		err = s.populateCropMovedArea(uid, &cropRead)
		if err != nil {
			result <- query.Result{Error: err}
		}

		err = s.populateCropHarvestedStorage(uid, &cropRead)
		if err != nil {
			result <- query.Result{Error: err}
		}

		err = s.populateCropTrash(uid, &cropRead)
		if err != nil {
			result <- query.Result{Error: err}
		}

		err = s.populateCropNotes(uid, &cropRead)
		if err != nil {
			result <- query.Result{Error: err}
		}

		result <- query.Result{Result: cropRead}
		close(result)
	}()

	return result
}

func (s CropReadQueryMysql) FindByBatchID(batchID string) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		cropRead := storage.CropRead{}
		rowsData := cropReadResult{}

		err := s.DB.QueryRow(`SELECT UID, BATCH_ID FROM CROP_READ WHERE BATCH_ID = ?`, batchID).Scan(
			&rowsData.UID,
			&rowsData.BatchID,
		)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Error: err}
		}

		if errors.Is(err, sql.ErrNoRows) {
			result <- query.Result{Result: cropRead}
		}

		cropUID, err := uuid.FromBytes(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		cropRead.UID = cropUID
		cropRead.BatchID = rowsData.BatchID

		result <- query.Result{Result: cropRead}
		close(result)
	}()

	return result
}

func (s CropReadQueryMysql) FindAllCropsByFarm(farmUID uuid.UUID, status string, page, limit int) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		// TODO: REFACTOR TO REDUCE QUERY CALLS
		cropReads := []storage.CropRead{}
		params := []interface{}{}

		offset := paginationhelper.CalculatePageToOffset(page, limit)

		sql := `SELECT UID FROM CROP_READ WHERE FARM_UID = ?`

		params = append(params, farmUID.Bytes())

		if status != "" {
			sql += ` AND STATUS = ?`

			params = append(params, status)
		}

		sql += ` ORDER BY INITIAL_AREA_CREATED_DATE DESC LIMIT ? OFFSET ?`

		params = append(params, limit, offset)

		rows, err := s.DB.Query(sql, params...)
		if err != nil {
			result <- query.Result{Error: err}
		}

		for rows.Next() {
			cropRead := storage.CropRead{}

			uid := []byte{}

			err := rows.Scan(&uid)
			if err != nil {
				result <- query.Result{Error: err}
			}

			cropUID, err := uuid.FromBytes(uid)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = s.populateCrop(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = s.populateCropMovedArea(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = s.populateCropHarvestedStorage(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = s.populateCropTrash(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = s.populateCropNotes(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = s.populateCropPhotos(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			cropReads = append(cropReads, cropRead)
		}

		result <- query.Result{Result: cropReads}
		close(result)
	}()

	return result
}

func (s CropReadQueryMysql) CountAllCropsByFarm(farmUID uuid.UUID, status string) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		total := 0
		params := []interface{}{}

		sql := `SELECT COUNT(UID) FROM CROP_READ WHERE FARM_UID = ?`

		params = append(params, farmUID.Bytes())

		if status != "" {
			sql += `  AND STATUS = ?`

			params = append(params, status)
		}

		err := s.DB.QueryRow(sql, params...).Scan(&total)
		if err != nil {
			result <- query.Result{Error: err}
		}

		result <- query.Result{Result: total}
		close(result)
	}()

	return result
}

func (s CropReadQueryMysql) FindAllCropsArchives(farmUID uuid.UUID, page, limit int) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		cropReads := []storage.CropRead{}

		// TODO: REFACTOR TO REDUCE QUERY CALLS

		offset := paginationhelper.CalculatePageToOffset(page, limit)

		rows, err := s.DB.Query(`SELECT UID FROM CROP_READ
			WHERE FARM_UID = ? AND STATUS = ? ORDER BY INITIAL_AREA_CREATED_DATE DESC LIMIT ? OFFSET ?`,
			farmUID.Bytes(), domain.CropArchived, limit, offset)
		if err != nil {
			result <- query.Result{Error: err}
		}

		for rows.Next() {
			cropRead := storage.CropRead{}

			uid := []byte{}

			err := rows.Scan(&uid)
			if err != nil {
				result <- query.Result{Error: err}
			}

			cropUID, err := uuid.FromBytes(uid)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = s.populateCrop(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = s.populateCropMovedArea(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = s.populateCropHarvestedStorage(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = s.populateCropTrash(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = s.populateCropNotes(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = s.populateCropPhotos(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			cropReads = append(cropReads, cropRead)
		}

		result <- query.Result{Result: cropReads}
		close(result)
	}()

	return result
}

func (s CropReadQueryMysql) CountAllArchivedCropsByFarm(farmUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		total := 0

		err := s.DB.QueryRow(`SELECT COUNT(UID) FROM CROP_READ
			WHERE FARM_UID = ? AND STATUS = ?`, farmUID.Bytes(), domain.CropArchived).Scan(&total)
		if err != nil {
			result <- query.Result{Error: err}
		}

		result <- query.Result{Result: total}
		close(result)
	}()

	return result
}

func (s CropReadQueryMysql) FindAllCropsByArea(areaUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		crops := []query.CropAreaByAreaQueryResult{}

		rows, err := s.DB.Query("SELECT UID FROM CROP_READ WHERE INITIAL_AREA_UID = ?", areaUID.Bytes())
		if err != nil {
			result <- query.Result{Error: err}
		}

		for rows.Next() {
			cropRead := storage.CropRead{}

			uid := []byte{}

			err := rows.Scan(&uid)
			if err != nil {
				result <- query.Result{Error: err}
			}

			cropUID, err := uuid.FromBytes(uid)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = s.populateCrop(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			if cropRead.InitialArea.AreaUID == areaUID {
				crops = append(crops, query.CropAreaByAreaQueryResult{
					UID:         cropRead.UID,
					BatchID:     cropRead.BatchID,
					CreatedDate: cropRead.InitialArea.CreatedDate,
					Area: query.Area{
						UID:             cropRead.InitialArea.AreaUID,
						Name:            cropRead.InitialArea.Name,
						InitialQuantity: cropRead.InitialArea.InitialQuantity,
						CurrentQuantity: cropRead.InitialArea.CurrentQuantity,
						InitialArea: query.InitialArea{
							UID:         cropRead.InitialArea.AreaUID,
							Name:        cropRead.InitialArea.Name,
							CreatedDate: cropRead.InitialArea.CreatedDate,
						},
						LastWatered: cropRead.InitialArea.LastWatered,
						MovingDate:  cropRead.InitialArea.CreatedDate,
					},
					Container: query.Container{
						Type:     cropRead.Container.Type,
						Cell:     cropRead.Container.Cell,
						Quantity: cropRead.Container.Quantity,
					},
					Inventory: query.Inventory{
						UID:       cropRead.Inventory.UID,
						Name:      cropRead.Inventory.Name,
						PlantType: cropRead.Inventory.PlantType,
					},
				})
			}
		}

		rows, err = s.DB.Query(`SELECT UID FROM CROP_READ
			LEFT JOIN CROP_READ_MOVED_AREA ON CROP_READ.UID = CROP_READ_MOVED_AREA.CROP_UID
			WHERE CROP_READ_MOVED_AREA.AREA_UID = ?`, areaUID.Bytes())
		if err != nil {
			result <- query.Result{Error: err}
		}

		for rows.Next() {
			cropRead := storage.CropRead{}

			uid := []byte{}

			err := rows.Scan(&uid)
			if err != nil {
				result <- query.Result{Error: err}
			}

			cropUID, err := uuid.FromBytes(uid)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = s.populateCrop(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = s.populateCropMovedArea(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			for _, val := range cropRead.MovedArea {
				if val.AreaUID == areaUID {
					crops = append(crops, query.CropAreaByAreaQueryResult{
						UID:         cropRead.UID,
						BatchID:     cropRead.BatchID,
						CreatedDate: val.CreatedDate,
						Area: query.Area{
							UID:             val.AreaUID,
							Name:            val.Name,
							InitialQuantity: val.InitialQuantity,
							CurrentQuantity: val.CurrentQuantity,
							InitialArea: query.InitialArea{
								UID:         cropRead.InitialArea.AreaUID,
								Name:        cropRead.InitialArea.Name,
								CreatedDate: cropRead.InitialArea.CreatedDate,
							},
							LastWatered: val.LastWatered,
							MovingDate:  val.CreatedDate,
						},
						Container: query.Container{
							Type:     cropRead.Container.Type,
							Cell:     cropRead.Container.Cell,
							Quantity: cropRead.Container.Quantity,
						},
						Inventory: query.Inventory{
							UID:       cropRead.Inventory.UID,
							Name:      cropRead.Inventory.Name,
							PlantType: cropRead.Inventory.PlantType,
						},
					})
				}
			}
		}

		result <- query.Result{Result: crops}
		close(result)
	}()

	return result
}

func (s CropReadQueryMysql) FindCropsInformation(farmUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		cropInf := query.CropInformationQueryResult{}
		harvestProduced := float32(0)
		plantType := make(map[string]bool)
		totalPlantVariety := 0

		// TODO: REFACTOR TO REDUCE QUERY CALLS
		rows, err := s.DB.Query("SELECT UID FROM CROP_READ WHERE FARM_UID = ?", farmUID.Bytes())
		if err != nil {
			result <- query.Result{Error: err}
		}

		for rows.Next() {
			cropRead := storage.CropRead{}

			uid := []byte{}

			err := rows.Scan(&uid)
			if err != nil {
				result <- query.Result{Error: err}
			}

			cropUID, err := uuid.FromBytes(uid)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = s.populateCrop(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = s.populateCropHarvestedStorage(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			for _, val := range cropRead.HarvestedStorage {
				harvestProduced += val.ProducedGramQuantity
			}

			if _, ok := plantType[cropRead.Inventory.Name]; !ok {
				totalPlantVariety++

				plantType[cropRead.Inventory.Name] = true
			}
		}

		cropInf.TotalHarvestProduced = harvestProduced
		cropInf.TotalPlantVariety = totalPlantVariety

		result <- query.Result{Result: cropInf}

		close(result)
	}()

	return result
}

func (s CropReadQueryMysql) CountTotalBatch(farmUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		varQty := []query.CountTotalBatchQueryResult{}
		varietyName := make(map[string]int)

		// TODO: REFACTOR TO REDUCE QUERY CALLS
		rows, err := s.DB.Query("SELECT UID FROM CROP_READ WHERE FARM_UID = ?", farmUID.Bytes())
		if err != nil {
			result <- query.Result{Error: err}
		}

		for rows.Next() {
			cropRead := storage.CropRead{}

			uid := []byte{}

			err := rows.Scan(&uid)
			if err != nil {
				result <- query.Result{Error: err}
			}

			cropUID, err := uuid.FromBytes(uid)
			if err != nil {
				result <- query.Result{Error: err}
			}

			err = s.populateCrop(cropUID, &cropRead)
			if err != nil {
				result <- query.Result{Error: err}
			}

			varietyName[cropRead.Inventory.Name]++
		}

		for i, v := range varietyName {
			varQty = append(varQty, query.CountTotalBatchQueryResult{
				VarietyName: i,
				TotalBatch:  v,
			})
		}

		result <- query.Result{Result: varQty}
		close(result)
	}()

	return result
}

func (s CropReadQueryMysql) populateCrop(cropUID uuid.UUID, cropRead *storage.CropRead) error {
	rowsData := cropReadResult{}

	err := s.DB.QueryRow(`SELECT UID, BATCH_ID, STATUS, TYPE, CONTAINER_QUANTITY, CONTAINER_TYPE, CONTAINER_CELL,
		INVENTORY_UID, INVENTORY_TYPE, INVENTORY_PLANT_TYPE, INVENTORY_NAME,
		AREA_STATUS_SEEDING, AREA_STATUS_GROWING, AREA_STATUS_DUMPED,
		FARM_UID,
		INITIAL_AREA_UID, INITIAL_AREA_NAME,
		INITIAL_AREA_INITIAL_QUANTITY, INITIAL_AREA_CURRENT_QUANTITY,
		INITIAL_AREA_LAST_WATERED, INITIAL_AREA_LAST_FERTILIZED, INITIAL_AREA_LAST_PESTICIDED,
		INITIAL_AREA_LAST_PRUNED, INITIAL_AREA_CREATED_DATE, INITIAL_AREA_LAST_UPDATED
		FROM CROP_READ WHERE UID = ?`, cropUID.Bytes()).Scan(
		&rowsData.UID,
		&rowsData.BatchID,
		&rowsData.Status,
		&rowsData.Type,
		&rowsData.ContainerQuantity,
		&rowsData.ContainerType,
		&rowsData.ContainerCell,
		&rowsData.InventoryUID,
		&rowsData.InventoryType,
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

	farmUID, err := uuid.FromBytes(rowsData.FarmUID)
	if err != nil {
		return err
	}

	inventoryUID, err := uuid.FromBytes(rowsData.InventoryUID)
	if err != nil {
		return err
	}

	initialAreaUID, err := uuid.FromBytes(rowsData.InitialAreaUID)
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

	cropRead.UID = cropUID
	cropRead.BatchID = rowsData.BatchID
	cropRead.Status = rowsData.Status
	cropRead.Type = rowsData.Type
	cropRead.Container.Quantity = rowsData.ContainerQuantity
	cropRead.Container.Type = rowsData.ContainerType
	cropRead.Container.Cell = rowsData.ContainerCell
	cropRead.Inventory.UID = inventoryUID
	cropRead.Inventory.Type = rowsData.InventoryType
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
	cropRead.InitialArea.CreatedDate = rowsData.InitialAreaCreatedDate
	cropRead.InitialArea.LastUpdated = rowsData.InitialAreaLastUpdated

	return nil
}

func (s CropReadQueryMysql) populateCropPhotos(uid uuid.UUID, cropRead *storage.CropRead) error {
	photoRowsData := cropReadPhotoResult{}

	rows, err := s.DB.Query("SELECT * FROM CROP_READ_PHOTO WHERE CROP_UID = ?", uid.Bytes())
	if err != nil {
		return err
	}

	photos := []storage.CropPhoto{}

	for rows.Next() {
		err = rows.Scan(
			&photoRowsData.UID,
			&photoRowsData.CropUID,
			&photoRowsData.Filename,
			&photoRowsData.Mimetype,
			&photoRowsData.Size,
			&photoRowsData.Width,
			&photoRowsData.Height,
			&photoRowsData.Description,
		)

		if err != nil {
			return err
		}

		photoUID, err := uuid.FromBytes(photoRowsData.UID)
		if err != nil {
			return err
		}

		photos = append(photos, storage.CropPhoto{
			UID:         photoUID,
			Filename:    photoRowsData.Filename,
			MimeType:    photoRowsData.Mimetype,
			Size:        photoRowsData.Size,
			Width:       photoRowsData.Width,
			Height:      photoRowsData.Height,
			Description: photoRowsData.Description,
		})
	}

	cropRead.Photos = photos

	return nil
}

func (s CropReadQueryMysql) populateCropMovedArea(uid uuid.UUID, cropRead *storage.CropRead) error {
	movedRowsData := cropReadMovedAreaResult{}

	rows, err := s.DB.Query("SELECT * FROM CROP_READ_MOVED_AREA WHERE CROP_UID = ?", uid.Bytes())
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

		areaUID, err := uuid.FromBytes(movedRowsData.AreaUID)
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
			CreatedDate:     movedRowsData.CreatedDate,
			LastUpdated:     movedRowsData.LastUpdated,
		})
	}

	cropRead.MovedArea = movedAreas

	return nil
}

func (s CropReadQueryMysql) populateCropHarvestedStorage(uid uuid.UUID, cropRead *storage.CropRead) error {
	harvestedRowsData := cropReadHarvestedStorageResult{}

	rows, err := s.DB.Query("SELECT * FROM CROP_READ_HARVESTED_STORAGE WHERE CROP_UID = ?", uid.Bytes())
	if err != nil {
		return err
	}

	harvestedStorages := []storage.HarvestedStorage{}

	for rows.Next() {
		err = rows.Scan(
			&harvestedRowsData.ID,
			&harvestedRowsData.CropUID,
			&harvestedRowsData.Quantity,
			&harvestedRowsData.ProducedGramQuantity,
			&harvestedRowsData.SourceAreaUID,
			&harvestedRowsData.SourceAreaName,
			&harvestedRowsData.CreatedDate,
			&harvestedRowsData.LastUpdated)
		if err != nil {
			return err
		}

		sourceAreaUID, err := uuid.FromBytes(harvestedRowsData.SourceAreaUID)
		if err != nil {
			return err
		}

		harvestedStorages = append(harvestedStorages, storage.HarvestedStorage{
			Quantity:             harvestedRowsData.Quantity,
			ProducedGramQuantity: harvestedRowsData.ProducedGramQuantity,
			SourceAreaUID:        sourceAreaUID,
			SourceAreaName:       harvestedRowsData.SourceAreaName,
			CreatedDate:          harvestedRowsData.CreatedDate,
			LastUpdated:          harvestedRowsData.LastUpdated,
		})
	}

	cropRead.HarvestedStorage = harvestedStorages

	return nil
}

func (s CropReadQueryMysql) populateCropTrash(uid uuid.UUID, cropRead *storage.CropRead) error {
	trashRowsData := cropReadTrashResult{}

	rows, err := s.DB.Query("SELECT * FROM CROP_READ_TRASH WHERE CROP_UID = ?", uid.Bytes())
	if err != nil {
		return err
	}

	trash := []storage.Trash{}

	for rows.Next() {
		err = rows.Scan(
			&trashRowsData.ID,
			&trashRowsData.CropUID,
			&trashRowsData.Quantity,
			&trashRowsData.SourceAreaUID,
			&trashRowsData.SourceAreaName,
			&trashRowsData.CreatedDate,
			&trashRowsData.LastUpdated)
		if err != nil {
			return err
		}

		sourceAreaUID, err := uuid.FromBytes(trashRowsData.SourceAreaUID)
		if err != nil {
			return err
		}

		trash = append(trash, storage.Trash{
			Quantity:       trashRowsData.Quantity,
			SourceAreaUID:  sourceAreaUID,
			SourceAreaName: trashRowsData.SourceAreaName,
			CreatedDate:    trashRowsData.CreatedDate,
			LastUpdated:    trashRowsData.LastUpdated,
		})
	}

	cropRead.Trash = trash

	return nil
}

func (s CropReadQueryMysql) populateCropNotes(uid uuid.UUID, cropRead *storage.CropRead) error {
	notesRowsData := cropReadNotesResult{}

	rows, err := s.DB.Query("SELECT * FROM CROP_READ_NOTES WHERE CROP_UID = ?", uid.Bytes())
	if err != nil {
		return err
	}

	notes := []domain.CropNote{}

	for rows.Next() {
		rows.Scan(
			&notesRowsData.UID,
			&notesRowsData.CropUID,
			&notesRowsData.Content,
			&notesRowsData.CreatedDate,
		)

		noteUID, err := uuid.FromBytes(notesRowsData.UID)
		if err != nil {
			return err
		}

		notes = append(notes, domain.CropNote{
			UID:         noteUID,
			Content:     notesRowsData.Content,
			CreatedDate: notesRowsData.CreatedDate,
		})
	}

	cropRead.Notes = notes

	return nil
}
