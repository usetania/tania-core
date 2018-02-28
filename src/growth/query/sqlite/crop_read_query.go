package sqlite

import (
	"database/sql"
	"time"

	"github.com/Tanibox/tania-server/src/growth/domain"
	"github.com/Tanibox/tania-server/src/growth/query"
	"github.com/Tanibox/tania-server/src/growth/storage"
	uuid "github.com/satori/go.uuid"
)

type CropReadQuerySqlite struct {
	DB *sql.DB
}

func NewCropReadQuerySqlite(db *sql.DB) query.CropReadQuery {
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

type cropReadPhotoResult struct {
	UID         string
	CropUID     string
	Filename    string
	Mimetype    string
	Size        int
	Width       int
	Height      int
	Description string
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

type cropReadHarvestedStorageResult struct {
	ID                   int
	CropUID              string
	Quantity             int
	ProducedGramQuantity float32
	SourceAreaUID        string
	SourceAreaName       string
	CreatedDate          string
	LastUpdated          string
}

type cropReadTrashResult struct {
	ID             int
	CropUID        string
	Quantity       int
	SourceAreaUID  string
	SourceAreaName string
	CreatedDate    string
	LastUpdated    string
}

type cropReadNotesResult struct {
	UID         string
	CropUID     string
	Content     string
	CreatedDate string
}

func (s CropReadQuerySqlite) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		cropRead := storage.CropRead{}
		rowsData := cropReadResult{}
		photoRowsData := cropReadPhotoResult{}
		movedRowsData := cropReadMovedAreaResult{}
		harvestedRowsData := cropReadHarvestedStorageResult{}
		trashRowsData := cropReadTrashResult{}
		notesRowsData := cropReadNotesResult{}

		err := s.DB.QueryRow(`SELECT UID, BATCH_ID, STATUS, TYPE, CONTAINER_QUANTITY, CONTAINER_TYPE, CONTAINER_CELL,
			INVENTORY_UID, INVENTORY_PLANT_TYPE, INVENTORY_NAME,
			AREA_STATUS_SEEDING, AREA_STATUS_GROWING, AREA_STATUS_DUMPED,
			FARM_UID,
			INITIAL_AREA_UID, INITIAL_AREA_NAME,
			INITIAL_AREA_INITIAL_QUANTITY, INITIAL_AREA_CURRENT_QUANTITY,
			INITIAL_AREA_LAST_WATERED, INITIAL_AREA_LAST_FERTILIZED, INITIAL_AREA_LAST_PESTICIDED,
			INITIAL_AREA_LAST_PRUNED, INITIAL_AREA_CREATED_DATE, INITIAL_AREA_LAST_UPDATED
			FROM CROP_READ WHERE UID = ?`, uid).Scan(
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

		if err != nil && err != sql.ErrNoRows {
			result <- query.QueryResult{Error: err}
		}

		if err == sql.ErrNoRows {
			result <- query.QueryResult{Result: cropRead}
		}

		cropUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		farmUID, err := uuid.FromString(rowsData.FarmUID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		inventoryUID, err := uuid.FromString(rowsData.InventoryUID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		initialAreaUID, err := uuid.FromString(rowsData.InitialAreaUID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		var initialAreaLastWatered *time.Time
		if rowsData.InitialAreaLastWatered.Valid && rowsData.InitialAreaLastWatered.String != "" {
			date, err := time.Parse(time.RFC3339, rowsData.InitialAreaLastWatered.String)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			initialAreaLastWatered = &date
		}

		var initialAreaLastFertilized *time.Time
		if rowsData.InitialAreaLastFertilized.Valid && rowsData.InitialAreaLastFertilized.String != "" {
			date, err := time.Parse(time.RFC3339, rowsData.InitialAreaLastFertilized.String)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			initialAreaLastFertilized = &date
		}

		var initialAreaLastPesticided *time.Time
		if rowsData.InitialAreaLastPesticided.Valid && rowsData.InitialAreaLastPesticided.String != "" {
			date, err := time.Parse(time.RFC3339, rowsData.InitialAreaLastPesticided.String)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			initialAreaLastPesticided = &date
		}

		var initialAreaLastPruned *time.Time
		if rowsData.InitialAreaLastPruned.Valid && rowsData.InitialAreaLastPruned.String != "" {
			date, err := time.Parse(time.RFC3339, rowsData.InitialAreaLastPruned.String)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			initialAreaLastPruned = &date
		}

		initialAreaCreatedDate, err := time.Parse(time.RFC3339, rowsData.InitialAreaCreatedDate)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		initialAreaLastUpdated, err := time.Parse(time.RFC3339, rowsData.InitialAreaLastUpdated)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		rows, err := s.DB.Query("SELECT * FROM CROP_READ_PHOTO WHERE CROP_UID = ?", uid)
		if err != nil {
			result <- query.QueryResult{Error: err}
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
				result <- query.QueryResult{Error: err}
			}

			photoUID, err := uuid.FromString(photoRowsData.UID)
			if err != nil {
				result <- query.QueryResult{Error: err}
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

		rows, err = s.DB.Query("SELECT * FROM CROP_READ_MOVED_AREA WHERE CROP_UID = ?", uid)
		if err != nil {
			result <- query.QueryResult{Error: err}
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
				result <- query.QueryResult{Error: err}
			}

			var lw *time.Time
			if movedRowsData.LastWatered.Valid && movedRowsData.LastWatered.String != "" {
				date, err := time.Parse(time.RFC3339, movedRowsData.LastWatered.String)
				if err != nil {
					result <- query.QueryResult{Error: err}
				}

				lw = &date
			}

			var lf *time.Time
			if movedRowsData.LastFertilized.Valid && movedRowsData.LastFertilized.String != "" {
				date, err := time.Parse(time.RFC3339, movedRowsData.LastFertilized.String)
				if err != nil {
					result <- query.QueryResult{Error: err}
				}

				lf = &date
			}

			var lp *time.Time
			if movedRowsData.LastPesticided.Valid && movedRowsData.LastPesticided.String != "" {
				date, err := time.Parse(time.RFC3339, movedRowsData.LastPesticided.String)
				if err != nil {
					result <- query.QueryResult{Error: err}
				}

				lp = &date
			}

			var lpr *time.Time
			if movedRowsData.LastPruned.Valid && movedRowsData.LastPruned.String != "" {
				date, err := time.Parse(time.RFC3339, movedRowsData.LastPruned.String)
				if err != nil {
					result <- query.QueryResult{Error: err}
				}

				lpr = &date
			}

			areaUID, err := uuid.FromString(movedRowsData.AreaUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			createdDate, err := time.Parse(time.RFC3339, movedRowsData.CreatedDate)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			lastUpdated, err := time.Parse(time.RFC3339, movedRowsData.LastUpdated)
			if err != nil {
				result <- query.QueryResult{Error: err}
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

		rows, err = s.DB.Query("SELECT * FROM CROP_READ_HARVESTED_STORAGE WHERE CROP_UID = ?", uid)
		if err != nil {
			result <- query.QueryResult{Error: err}
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

			sourceAreaUID, err := uuid.FromString(harvestedRowsData.SourceAreaUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			createdDate, err := time.Parse(time.RFC3339, harvestedRowsData.CreatedDate)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			lastUpdated, err := time.Parse(time.RFC3339, harvestedRowsData.LastUpdated)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			harvestedStorages = append(harvestedStorages, storage.HarvestedStorage{
				Quantity:             harvestedRowsData.Quantity,
				ProducedGramQuantity: harvestedRowsData.ProducedGramQuantity,
				SourceAreaUID:        sourceAreaUID,
				SourceAreaName:       harvestedRowsData.SourceAreaName,
				CreatedDate:          createdDate,
				LastUpdated:          lastUpdated,
			})
		}

		rows, err = s.DB.Query("SELECT * FROM CROP_READ_TRASH WHERE CROP_UID = ?", uid)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		trash := []storage.Trash{}
		for rows.Next() {
			err = rows.Scan(
				&trashRowsData.ID,
				&trashRowsData.CropUID,
				&trashRowsData.SourceAreaUID,
				&trashRowsData.SourceAreaName,
				&trashRowsData.CreatedDate,
				&trashRowsData.LastUpdated)

			sourceAreaUID, err := uuid.FromString(harvestedRowsData.SourceAreaUID)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			createdDate, err := time.Parse(time.RFC3339, harvestedRowsData.CreatedDate)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			lastUpdated, err := time.Parse(time.RFC3339, harvestedRowsData.LastUpdated)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			trash = append(trash, storage.Trash{
				Quantity:       harvestedRowsData.Quantity,
				SourceAreaUID:  sourceAreaUID,
				SourceAreaName: harvestedRowsData.SourceAreaName,
				CreatedDate:    createdDate,
				LastUpdated:    lastUpdated,
			})
		}

		rows, err = s.DB.Query("SELECT * FROM CROP_READ_NOTES WHERE CROP_UID = ?", uid)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		notes := []domain.CropNote{}
		for rows.Next() {
			rows.Scan(
				&notesRowsData.UID,
				&notesRowsData.CropUID,
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

			notes = append(notes, domain.CropNote{
				UID:         noteUID,
				Content:     notesRowsData.Content,
				CreatedDate: noteCreatedDate,
			})
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
		cropRead.Photos = photos
		cropRead.MovedArea = movedAreas
		cropRead.HarvestedStorage = harvestedStorages
		cropRead.Trash = trash
		cropRead.Notes = notes

		result <- query.QueryResult{Result: cropRead}
		close(result)
	}()

	return result
}

func (s CropReadQuerySqlite) FindByBatchID(batchID string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		cropRead := storage.CropRead{}
		rowsData := cropReadResult{}

		err := s.DB.QueryRow(`SELECT UID, BATCH_ID FROM CROP_READ WHERE BATCH_ID = ?`, batchID).Scan(
			&rowsData.UID,
			&rowsData.BatchID,
		)

		if err != nil && err != sql.ErrNoRows {
			result <- query.QueryResult{Error: err}
		}

		if err == sql.ErrNoRows {
			result <- query.QueryResult{Result: cropRead}
		}

		cropUID, err := uuid.FromString(rowsData.UID)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		cropRead.UID = cropUID
		cropRead.BatchID = rowsData.BatchID

		result <- query.QueryResult{Result: cropRead}
		close(result)
	}()

	return result
}

func (s CropReadQuerySqlite) FindAllCropsByFarm(farmUID uuid.UUID) <-chan query.QueryResult {
	return nil
}

func (s CropReadQuerySqlite) FindAllCropsArchives(farmUID uuid.UUID) <-chan query.QueryResult {
	return nil
}

func (s CropReadQuerySqlite) FindAllCropsByArea(areaUID uuid.UUID) <-chan query.QueryResult {
	return nil
}

func (s CropReadQuerySqlite) FindCropsInformation(farmUID uuid.UUID) <-chan query.QueryResult {
	return nil
}

func (s CropReadQuerySqlite) CountTotalBatch(farmUID uuid.UUID) <-chan query.QueryResult {
	return nil
}
