package sqlite

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/tasks/domain"
	"github.com/usetania/tania-core/src/tasks/repository"
	"github.com/usetania/tania-core/src/tasks/storage"
)

type TaskReadRepositorySqlite struct {
	DB *sql.DB
}

func NewTaskReadRepositorySqlite(s *sql.DB) repository.TaskRead {
	return &TaskReadRepositorySqlite{DB: s}
}

func (f *TaskReadRepositorySqlite) Save(taskRead *storage.TaskRead) <-chan error {
	result := make(chan error)

	go func() {
		var dueDate *string

		if taskRead.DueDate != nil && !taskRead.DueDate.IsZero() {
			d := taskRead.DueDate.Format(time.RFC3339)
			dueDate = &d
		}

		var completedDate *string

		if taskRead.CompletedDate != nil && !taskRead.CompletedDate.IsZero() {
			d := taskRead.CompletedDate.Format(time.RFC3339)
			completedDate = &d
		}

		var cancelledDate *string

		if taskRead.CancelledDate != nil && !taskRead.CancelledDate.IsZero() {
			d := taskRead.CancelledDate.Format(time.RFC3339)
			cancelledDate = &d
		}

		var domainDataMaterialID, domainDataAreaID *uuid.UUID

		switch v := taskRead.DomainDetails.(type) {
		case domain.TaskDomainArea:
			domainDataMaterialID = v.MaterialID
		case domain.TaskDomainCrop:
			domainDataMaterialID = v.MaterialID
			domainDataAreaID = v.AreaID
		case domain.TaskDomainReservoir:
			domainDataMaterialID = v.MaterialID
		}

		res, err := f.DB.Exec(`UPDATE TASK_READ SET
			TITLE = ?, DESCRIPTION = ?, CREATED_DATE = ?, DUE_DATE = ?,
			COMPLETED_DATE = ?, CANCELLED_DATE = ?, PRIORITY = ?, STATUS = ?,
			DOMAIN_CODE = ?, DOMAIN_DATA_MATERIAL_ID = ?, DOMAIN_DATA_AREA_ID = ?,
			CATEGORY = ?, IS_DUE = ?, ASSET_ID = ?
			WHERE UID = ?`,
			taskRead.Title, taskRead.Description, taskRead.CreatedDate.Format(time.RFC3339), dueDate,
			completedDate, cancelledDate, taskRead.Priority, taskRead.Status,
			taskRead.Domain, domainDataMaterialID, domainDataAreaID, taskRead.Category, taskRead.IsDue, taskRead.AssetID,
			taskRead.UID)
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
				taskRead.UID, taskRead.Title, taskRead.Description, taskRead.CreatedDate.Format(time.RFC3339), dueDate,
				completedDate, cancelledDate, taskRead.Priority, taskRead.Status,
				taskRead.Domain, domainDataMaterialID, domainDataAreaID, taskRead.Category, taskRead.IsDue, taskRead.AssetID)
			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}
