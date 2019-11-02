package mysql

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/Tanibox/tania-core/src/helper/paginationhelper"
	"github.com/Tanibox/tania-core/src/tasks/domain"
	"github.com/Tanibox/tania-core/src/tasks/query"
	"github.com/Tanibox/tania-core/src/tasks/storage"
	uuid "github.com/satori/go.uuid"
)

type TaskReadQueryMysql struct {
	DB *sql.DB
}

func NewTaskReadQueryMysql(s *sql.DB) query.TaskReadQuery {
	return &TaskReadQueryMysql{DB: s}
}

type taskReadQueryResult struct {
	UID                  []byte
	Title                string
	Description          string
	CreatedDate          time.Time
	DueDate              *time.Time
	CompletedDate        *time.Time
	CancelledDate        *time.Time
	Priority             string
	Status               string
	DomainCode           string
	DomainDataMaterialID uuid.NullUUID
	DomainDataAreaID     uuid.NullUUID
	DomainDataCropID     uuid.NullUUID
	Category             string
	IsDue                int
	AssetID              uuid.NullUUID
}

func (r TaskReadQueryMysql) FindAll(page, limit int) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		tasks := []storage.TaskRead{}

		sql := `SELECT * FROM TASK_READ ORDER BY CREATED_DATE DESC`
		var args []interface{}

		if page != 0 && limit != 0 {
			sql += " LIMIT ? OFFSET ?"
			offset := paginationhelper.CalculatePageToOffset(page, limit)
			args = append(args, limit, offset)
		}

		rows, err := r.DB.Query(sql, args...)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		for rows.Next() {
			taskRead, err := r.populateQueryResult(rows)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			tasks = append(tasks, taskRead)
		}

		result <- query.QueryResult{Result: tasks}

		close(result)
	}()

	return result
}

// FindByID is to find by ID
func (r TaskReadQueryMysql) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		task := storage.TaskRead{}

		rows, err := r.DB.Query(`SELECT * FROM TASK_READ WHERE UID = ?`, uid.Bytes())
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		for rows.Next() {
			taskRead, err := r.populateQueryResult(rows)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			task = taskRead
		}

		result <- query.QueryResult{Result: task}
		close(result)
	}()

	return result
}

func (s TaskReadQueryMysql) FindTasksWithFilter(params map[string]string, page, limit int) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		tasks := []storage.TaskRead{}

		sql := "SELECT * FROM TASK_READ WHERE 1 = 1"
		var args []interface{}

		if value, _ := params["is_due"]; value != "" {
			b, _ := strconv.ParseBool(value)
			sql += " AND IS_DUE = ? "
			args = append(args, b)
		}
		start, _ := params["due_start"]
		end, _ := params["due_end"]
		if start != "" && end != "" {
			sql += " AND DUE_DATE BETWEEN ? AND ? "
			args = append(args, start)
			args = append(args, end)
		}
		if value, _ := params["priority"]; value != "" {
			sql += " AND PRIORITY = ? "
			args = append(args, value)
		}
		if value, _ := params["status"]; value != "" {
			sql += " AND STATUS = ? "
			args = append(args, value)
		}
		if value, _ := params["domain"]; value != "" {
			sql += " AND DOMAIN_CODE = ? "
			args = append(args, value)
		}
		if value, _ := params["category"]; value != "" {
			sql += " AND CATEGORY = ? "
			args = append(args, value)
		}
		if value, _ := params["asset_id"]; value != "" {
			assetID, _ := uuid.FromString(value)
			sql += " AND ASSET_ID = ? "
			args = append(args, assetID.Bytes())
		}

		if page != 0 && limit != 0 {
			sql += " LIMIT ? OFFSET ?"
			offset := paginationhelper.CalculatePageToOffset(page, limit)
			args = append(args, limit, offset)
		}

		rows, err := s.DB.Query(sql, args...)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		for rows.Next() {
			taskRead, err := s.populateQueryResult(rows)
			if err != nil {
				result <- query.QueryResult{Error: err}
			}

			tasks = append(tasks, taskRead)
		}

		result <- query.QueryResult{Result: tasks}

		close(result)
	}()

	return result
}

func (q TaskReadQueryMysql) CountAll() <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		total := 0
		var params []interface{}

		sql := "SELECT COUNT(UID) FROM TASK_READ"

		err := q.DB.QueryRow(sql, params...).Scan(&total)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		result <- query.QueryResult{Result: total}
		close(result)
	}()

	return result
}

func (q TaskReadQueryMysql) CountTasksWithFilter(params map[string]string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		total := 0

		sql := "SELECT COUNT(UID) FROM TASK_READ WHERE 1 = 1"
		var args []interface{}

		if value, _ := params["is_due"]; value != "" {
			b, _ := strconv.ParseBool(value)
			sql += " AND IS_DUE = ? "
			args = append(args, b)
		}
		start, _ := params["due_start"]
		end, _ := params["due_end"]
		if start != "" && end != "" {
			sql += " AND DUE_DATE BETWEEN ? AND ? "
			args = append(args, start)
			args = append(args, end)
		}
		if value, _ := params["priority"]; value != "" {
			sql += " AND PRIORITY = ? "
			args = append(args, value)
		}
		if value, _ := params["status"]; value != "" {
			sql += " AND STATUS = ? "
			args = append(args, value)
		}
		if value, _ := params["domain"]; value != "" {
			sql += " AND DOMAIN_CODE = ? "
			args = append(args, value)
		}
		if value, _ := params["category"]; value != "" {
			sql += " AND CATEGORY = ? "
			args = append(args, value)
		}
		if value, _ := params["asset_id"]; value != "" {
			assetID, _ := uuid.FromString(value)
			sql += " AND ASSET_ID = ? "
			args = append(args, assetID.Bytes())
		}

		err := q.DB.QueryRow(sql, args...).Scan(&total)
		if err != nil {
			result <- query.QueryResult{Error: err}
		}

		result <- query.QueryResult{Result: total}
		close(result)
	}()

	return result
}

func (s TaskReadQueryMysql) populateQueryResult(rows *sql.Rows) (storage.TaskRead, error) {
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
