package sqlite

import (
	"database/sql"

	"github.com/usetania/tania-core/src/tasks/domain"
	"github.com/usetania/tania-core/src/tasks/repository"
	"github.com/usetania/tania-core/src/tasks/storage"
)

type TaskReadRepositoryMysql struct {
	DB *sql.DB
}

func NewTaskReadRepositoryMysql(s *sql.DB) repository.TaskRead {
	return &TaskReadRepositoryMysql{DB: s}
}

func (f *TaskReadRepositoryMysql) Save(taskRead *storage.TaskRead) <-chan error {
	result := make(chan error)

	go func() {
		var domainDataMaterialID []byte

		var domainDataAreaID []byte

		switch v := taskRead.DomainDetails.(type) {
		case domain.TaskDomainCrop:
			if v.MaterialID != nil {
				domainDataMaterialID = v.MaterialID.Bytes()
			}

			if v.AreaID != nil {
				domainDataAreaID = v.AreaID.Bytes()
			}
		}

		var assetID []byte
		if taskRead.AssetID != nil {
			assetID = taskRead.AssetID.Bytes()
		}

		res, err := f.DB.Exec(`UPDATE TASK_READ SET
			TITLE = ?, DESCRIPTION = ?, CREATED_DATE = ?, DUE_DATE = ?,
			COMPLETED_DATE = ?, CANCELLED_DATE = ?, PRIORITY = ?, STATUS = ?,
			DOMAIN_CODE = ?, DOMAIN_DATA_MATERIAL_ID = ?, DOMAIN_DATA_AREA_ID = ?,
			CATEGORY = ?, IS_DUE = ?, ASSET_ID = ?
			WHERE UID = ?`,
			taskRead.Title, taskRead.Description, taskRead.CreatedDate, taskRead.DueDate,
			taskRead.CompletedDate, taskRead.CancelledDate, taskRead.Priority, taskRead.Status,
			taskRead.Domain, domainDataMaterialID, domainDataAreaID,
			taskRead.Category, taskRead.IsDue, assetID,
			taskRead.UID.Bytes())
		if err != nil {
			result <- err
		}

		rowsAffected := int64(0)
		if res != nil {
			rowsAffected, err = res.RowsAffected()
			if err != nil {
				result <- err
			}
		}

		if rowsAffected == 0 {
			_, err := f.DB.Exec(`INSERT INTO TASK_READ (
				UID, TITLE, DESCRIPTION, CREATED_DATE, DUE_DATE,
				COMPLETED_DATE, CANCELLED_DATE, PRIORITY, STATUS,
				DOMAIN_CODE, DOMAIN_DATA_MATERIAL_ID, DOMAIN_DATA_AREA_ID, CATEGORY, IS_DUE, ASSET_ID)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				taskRead.UID.Bytes(), taskRead.Title, taskRead.Description, taskRead.CreatedDate, taskRead.DueDate,
				taskRead.CompletedDate, taskRead.CancelledDate, taskRead.Priority, taskRead.Status,
				taskRead.Domain, domainDataMaterialID, domainDataAreaID,
				taskRead.Category, taskRead.IsDue, assetID)
			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}
