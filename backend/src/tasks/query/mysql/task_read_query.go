package mysql

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/helper/paginationhelper"
	"github.com/usetania/tania-core/src/tasks/domain"
	"github.com/usetania/tania-core/src/tasks/query"
	"github.com/usetania/tania-core/src/tasks/storage"
)

type TaskReadQueryMysql struct {
	DB *sql.DB
}

func NewTaskReadQueryMysql(s *sql.DB) query.TaskRead {
	return &TaskReadQueryMysql{DB: s}
}

type taskReadQueryResult struct {
	UID                  []byte
	Title                string
	Description          string
	Priority             string
	Status               string
	DomainCode           string
	Category             string
	IsDue                int
	CreatedDate          time.Time
	DueDate              *time.Time
	CompletedDate        *time.Time
	CancelledDate        *time.Time
	DomainDataMaterialID uuid.NullUUID
	DomainDataAreaID     uuid.NullUUID
	DomainDataCropID     uuid.NullUUID
	AssetID              uuid.NullUUID
}

func (q TaskReadQueryMysql) FindAll(page, limit int) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		tasks := []storage.TaskRead{}

		sql := `SELECT * FROM TASK_READ ORDER BY CREATED_DATE DESC`

		var args []interface{}

		if page != 0 && limit != 0 {
			sql += " LIMIT ? OFFSET ?"
			offset := paginationhelper.CalculatePageToOffset(page, limit)
			args = append(args, limit, offset)
		}

		rows, err := q.DB.Query(sql, args...)
		if err != nil {
			result <- query.Result{Error: err}
		}

		for rows.Next() {
			taskRead, err := q.populateQueryResult(rows)
			if err != nil {
				result <- query.Result{Error: err}
			}

			tasks = append(tasks, taskRead)
		}

		result <- query.Result{Result: tasks}

		close(result)
	}()

	return result
}

// FindByID is to find by ID.
func (q TaskReadQueryMysql) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		task := storage.TaskRead{}

		rows, err := q.DB.Query(`SELECT * FROM TASK_READ WHERE UID = ?`, uid.Bytes())
		if err != nil {
			result <- query.Result{Error: err}
		}

		for rows.Next() {
			taskRead, err := q.populateQueryResult(rows)
			if err != nil {
				result <- query.Result{Error: err}
			}

			task = taskRead
		}

		result <- query.Result{Result: task}
		close(result)
	}()

	return result
}

func (q TaskReadQueryMysql) FindTasksWithFilter(params map[string]string, page, limit int) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		tasks := []storage.TaskRead{}

		sql := "SELECT * FROM TASK_READ WHERE 1 = 1"

		var args []interface{}

		if value := params["is_due"]; value != "" {
			b, _ := strconv.ParseBool(value)
			sql += " AND IS_DUE = ? "

			args = append(args, b)
		}

		start := params["due_start"]
		end := params["due_end"]

		if start != "" && end != "" {
			sql += " AND DUE_DATE BETWEEN ? AND ? "

			args = append(args, start, end)
		}

		if value := params["priority"]; value != "" {
			sql += " AND PRIORITY = ? "

			args = append(args, value)
		}

		if value := params["status"]; value != "" {
			sql += " AND STATUS = ? "

			args = append(args, value)
		}

		if value := params["domain"]; value != "" {
			sql += " AND DOMAIN_CODE = ? "

			args = append(args, value)
		}

		if value := params["category"]; value != "" {
			sql += " AND CATEGORY = ? "

			args = append(args, value)
		}

		if value := params["asset_id"]; value != "" {
			assetID, _ := uuid.FromString(value)
			sql += " AND ASSET_ID = ? "

			args = append(args, assetID.Bytes())
		}

		if page != 0 && limit != 0 {
			sql += " LIMIT ? OFFSET ?"
			offset := paginationhelper.CalculatePageToOffset(page, limit)

			args = append(args, limit, offset)
		}

		rows, err := q.DB.Query(sql, args...)
		if err != nil {
			result <- query.Result{Error: err}
		}

		for rows.Next() {
			taskRead, err := q.populateQueryResult(rows)
			if err != nil {
				result <- query.Result{Error: err}
			}

			tasks = append(tasks, taskRead)
		}

		result <- query.Result{Result: tasks}

		close(result)
	}()

	return result
}

func (q TaskReadQueryMysql) CountAll() <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		total := 0

		var params []interface{}

		sql := "SELECT COUNT(UID) FROM TASK_READ"

		err := q.DB.QueryRow(sql, params...).Scan(&total)
		if err != nil {
			result <- query.Result{Error: err}
		}

		result <- query.Result{Result: total}
		close(result)
	}()

	return result
}

func (q TaskReadQueryMysql) CountTasksWithFilter(params map[string]string) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		total := 0

		sql := "SELECT COUNT(UID) FROM TASK_READ WHERE 1 = 1"

		var args []interface{}

		if value := params["is_due"]; value != "" {
			b, _ := strconv.ParseBool(value)
			sql += " AND IS_DUE = ? "

			args = append(args, b)
		}

		start := params["due_start"]
		end := params["due_end"]

		if start != "" && end != "" {
			sql += " AND DUE_DATE BETWEEN ? AND ? "

			args = append(args, start, end)
		}

		if value := params["priority"]; value != "" {
			sql += " AND PRIORITY = ? "

			args = append(args, value)
		}

		if value := params["status"]; value != "" {
			sql += " AND STATUS = ? "

			args = append(args, value)
		}

		if value := params["domain"]; value != "" {
			sql += " AND DOMAIN_CODE = ? "

			args = append(args, value)
		}

		if value := params["category"]; value != "" {
			sql += " AND CATEGORY = ? "

			args = append(args, value)
		}

		if value := params["asset_id"]; value != "" {
			assetID, _ := uuid.FromString(value)
			sql += " AND ASSET_ID = ? "

			args = append(args, assetID.Bytes())
		}

		err := q.DB.QueryRow(sql, args...).Scan(&total)
		if err != nil {
			result <- query.Result{Error: err}
		}

		result <- query.Result{Result: total}
		close(result)
	}()

	return result
}

func (TaskReadQueryMysql) populateQueryResult(rows *sql.Rows) (storage.TaskRead, error) {
	rowsData := taskReadQueryResult{}

	err := rows.Scan(
		&rowsData.UID, &rowsData.Title, &rowsData.Description, &rowsData.CreatedDate,
		&rowsData.DueDate, &rowsData.CompletedDate, &rowsData.CancelledDate,
		&rowsData.Priority, &rowsData.Status, &rowsData.DomainCode, &rowsData.DomainDataMaterialID,
		&rowsData.DomainDataAreaID, &rowsData.DomainDataCropID, &rowsData.Category, &rowsData.IsDue, &rowsData.AssetID,
	)
	if err != nil {
		return storage.TaskRead{}, err
	}

	taskUID, err := uuid.FromBytes(rowsData.UID)
	if err != nil {
		return storage.TaskRead{}, err
	}

	var domainDetails domain.TaskDomain

	switch rowsData.DomainCode {
	case domain.TaskDomainAreaCode:
		domainDetails = domain.TaskDomainArea{}
	case domain.TaskDomainCropCode:
		var materialID *uuid.UUID

		var areaID *uuid.UUID

		if rowsData.DomainDataMaterialID.Valid {
			materialID = &rowsData.DomainDataMaterialID.UUID
		}

		if rowsData.DomainDataAreaID.Valid {
			areaID = &rowsData.DomainDataAreaID.UUID
		}

		domainDetails = domain.TaskDomainCrop{
			MaterialID: materialID,
			AreaID:     areaID,
		}

	case domain.TaskDomainFinanceCode:
		domainDetails = domain.TaskDomainFinance{}
	case domain.TaskDomainGeneralCode:
		domainDetails = domain.TaskDomainGeneral{}
	case domain.TaskDomainInventoryCode:
		domainDetails = domain.TaskDomainInventory{}
	case domain.TaskDomainReservoirCode:
		domainDetails = domain.TaskDomainReservoir{}
	}

	assetUID := &uuid.UUID{}
	if rowsData.AssetID.Valid {
		assetUID = &rowsData.AssetID.UUID
	}

	isDue := false
	if rowsData.IsDue == 1 {
		isDue = true
	}

	return storage.TaskRead{
		UID:           taskUID,
		Title:         rowsData.Title,
		Description:   rowsData.Description,
		CreatedDate:   rowsData.CreatedDate,
		DueDate:       rowsData.DueDate,
		CompletedDate: rowsData.CompletedDate,
		CancelledDate: rowsData.CancelledDate,
		Priority:      rowsData.Priority,
		Status:        rowsData.Status,
		Domain:        rowsData.DomainCode,
		DomainDetails: domainDetails,
		Category:      rowsData.Category,
		IsDue:         isDue,
		AssetID:       assetUID,
	}, nil
}
