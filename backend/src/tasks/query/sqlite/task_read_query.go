package sqlite

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

type TaskReadQuerySqlite struct {
	DB *sql.DB
}

func NewTaskReadQuerySqlite(s *sql.DB) query.TaskRead {
	return &TaskReadQuerySqlite{DB: s}
}

type taskReadQueryResult struct {
	UID                  string
	Title                string
	Description          string
	CreatedDate          string
	DueDate              sql.NullString
	CompletedDate        sql.NullString
	CancelledDate        sql.NullString
	Priority             string
	Status               string
	DomainCode           string
	DomainDataMaterialID sql.NullString
	DomainDataAreaID     sql.NullString
	Category             string
	IsDue                bool
	AssetID              sql.NullString
}

func (q TaskReadQuerySqlite) FindAll(page, limit int) <-chan query.Result {
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
func (q TaskReadQuerySqlite) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		task := storage.TaskRead{}

		rows, err := q.DB.Query(`SELECT * FROM TASK_READ WHERE UID = ?`, uid)
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

func (q TaskReadQuerySqlite) FindTasksWithFilter(params map[string]string, page, limit int) <-chan query.Result {
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

			args = append(args, assetID)
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

func (q TaskReadQuerySqlite) CountAll() <-chan query.Result {
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

func (q TaskReadQuerySqlite) CountTasksWithFilter(params map[string]string) <-chan query.Result {
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

			args = append(args, assetID)
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

func (TaskReadQuerySqlite) populateQueryResult(rows *sql.Rows) (storage.TaskRead, error) {
	rowsData := taskReadQueryResult{}

	err := rows.Scan(
		&rowsData.UID, &rowsData.Title, &rowsData.Description, &rowsData.CreatedDate,
		&rowsData.DueDate, &rowsData.CompletedDate, &rowsData.CancelledDate,
		&rowsData.Priority, &rowsData.Status, &rowsData.DomainCode, &rowsData.DomainDataMaterialID,
		&rowsData.DomainDataAreaID,
		&rowsData.Category, &rowsData.IsDue, &rowsData.AssetID,
	)
	if err != nil {
		return storage.TaskRead{}, err
	}

	taskUID, err := uuid.FromString(rowsData.UID)
	if err != nil {
		return storage.TaskRead{}, err
	}

	createdDate, err := time.Parse(time.RFC3339, rowsData.CreatedDate)
	if err != nil {
		return storage.TaskRead{}, err
	}

	var dueDate *time.Time

	if rowsData.DueDate.Valid && rowsData.DueDate.String != "" {
		d, err := time.Parse(time.RFC3339, rowsData.DueDate.String)
		if err != nil {
			return storage.TaskRead{}, err
		}

		dueDate = &d
	}

	var completedDate *time.Time

	if rowsData.CompletedDate.Valid && rowsData.CompletedDate.String != "" {
		d, err := time.Parse(time.RFC3339, rowsData.CompletedDate.String)
		if err != nil {
			return storage.TaskRead{}, err
		}

		completedDate = &d
	}

	var cancelledDate *time.Time

	if rowsData.CancelledDate.Valid && rowsData.CancelledDate.String != "" {
		d, err := time.Parse(time.RFC3339, rowsData.CancelledDate.String)
		if err != nil {
			return storage.TaskRead{}, err
		}

		cancelledDate = &d
	}

	var domainDetails domain.TaskDomain

	switch rowsData.DomainCode {
	case domain.TaskDomainAreaCode:
		materialID := (*uuid.UUID)(nil)

		if rowsData.DomainDataMaterialID.Valid && rowsData.DomainDataMaterialID.String != "" {
			uid, err := uuid.FromString(rowsData.DomainDataMaterialID.String)
			if err != nil {
				return storage.TaskRead{}, err
			}

			materialID = &uid
		}

		domainDetails = domain.TaskDomainArea{
			MaterialID: materialID,
		}

	case domain.TaskDomainCropCode:
		materialID := (*uuid.UUID)(nil)
		areaID := (*uuid.UUID)(nil)

		if rowsData.DomainDataMaterialID.Valid && rowsData.DomainDataMaterialID.String != "" {
			uid, err := uuid.FromString(rowsData.DomainDataMaterialID.String)
			if err != nil {
				return storage.TaskRead{}, err
			}

			materialID = &uid
		}

		if rowsData.DomainDataAreaID.Valid && rowsData.DomainDataAreaID.String != "" {
			uid, err := uuid.FromString(rowsData.DomainDataAreaID.String)
			if err != nil {
				return storage.TaskRead{}, err
			}

			areaID = &uid
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
		materialID := (*uuid.UUID)(nil)

		if rowsData.DomainDataMaterialID.Valid && rowsData.DomainDataMaterialID.String != "" {
			uid, err := uuid.FromString(rowsData.DomainDataMaterialID.String)
			if err != nil {
				return storage.TaskRead{}, err
			}

			materialID = &uid
		}

		domainDetails = domain.TaskDomainReservoir{
			MaterialID: materialID,
		}
	}

	var assetUID *uuid.UUID

	if rowsData.AssetID.Valid && rowsData.AssetID.String != "" {
		uid, err := uuid.FromString(rowsData.AssetID.String)
		if err != nil {
			return storage.TaskRead{}, err
		}

		assetUID = &uid
	}

	return storage.TaskRead{
		UID:           taskUID,
		Title:         rowsData.Title,
		Description:   rowsData.Description,
		CreatedDate:   createdDate,
		DueDate:       dueDate,
		CompletedDate: completedDate,
		CancelledDate: cancelledDate,
		Priority:      rowsData.Priority,
		Status:        rowsData.Status,
		Domain:        rowsData.DomainCode,
		DomainDetails: domainDetails,
		Category:      rowsData.Category,
		IsDue:         rowsData.IsDue,
		AssetID:       assetUID,
	}, nil
}
