package server

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Tanibox/tania-core/config"
	assetsstorage "github.com/Tanibox/tania-core/src/assets/storage"
	"github.com/Tanibox/tania-core/src/eventbus"
	cropstorage "github.com/Tanibox/tania-core/src/growth/storage"
	"github.com/Tanibox/tania-core/src/helper/paginationhelper"
	"github.com/Tanibox/tania-core/src/helper/structhelper"
	"github.com/Tanibox/tania-core/src/tasks/domain"
	service "github.com/Tanibox/tania-core/src/tasks/domain/service"
	"github.com/Tanibox/tania-core/src/tasks/query"
	queryInMem "github.com/Tanibox/tania-core/src/tasks/query/inmemory"
	queryMysql "github.com/Tanibox/tania-core/src/tasks/query/mysql"
	querySqlite "github.com/Tanibox/tania-core/src/tasks/query/sqlite"
	"github.com/Tanibox/tania-core/src/tasks/repository"
	repoInMem "github.com/Tanibox/tania-core/src/tasks/repository/inmemory"
	repoMysql "github.com/Tanibox/tania-core/src/tasks/repository/mysql"
	repoSqlite "github.com/Tanibox/tania-core/src/tasks/repository/sqlite"
	"github.com/Tanibox/tania-core/src/tasks/storage"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

// TaskServer ties the routes and handlers with injected dependencies
type TaskServer struct {
	TaskEventRepo  repository.TaskEventRepository
	TaskReadRepo   repository.TaskReadRepository
	TaskEventQuery query.TaskEventQuery
	TaskReadQuery  query.TaskReadQuery
	TaskService    domain.TaskService
	EventBus       eventbus.TaniaEventBus
}

// NewTaskServer initializes TaskServer's dependencies and create new TaskServer struct
func NewTaskServer(
	db *sql.DB,
	bus eventbus.TaniaEventBus,
	cropStorage *cropstorage.CropReadStorage,
	areaStorage *assetsstorage.AreaReadStorage,
	materialStorage *assetsstorage.MaterialReadStorage,
	reservoirStorage *assetsstorage.ReservoirReadStorage,
	taskEventStorage *storage.TaskEventStorage,
	taskReadStorage *storage.TaskReadStorage) (*TaskServer, error) {

	taskServer := &TaskServer{
		EventBus: bus,
	}

	switch *config.Config.TaniaPersistenceEngine {
	case config.DB_INMEMORY:
		taskServer.TaskEventRepo = repoInMem.NewTaskEventRepositoryInMemory(taskEventStorage)
		taskServer.TaskReadRepo = repoInMem.NewTaskReadRepositoryInMemory(taskReadStorage)

		taskServer.TaskEventQuery = queryInMem.NewTaskEventQueryInMemory(taskEventStorage)
		taskServer.TaskReadQuery = queryInMem.NewTaskReadQueryInMemory(taskReadStorage)

		cropQuery := queryInMem.NewCropQueryInMemory(cropStorage)
		areaQuery := queryInMem.NewAreaQueryInMemory(areaStorage)
		materialReadQuery := queryInMem.NewMaterialQueryInMemory(materialStorage)
		reservoirQuery := queryInMem.NewReservoirQueryInMemory(reservoirStorage)

		taskServer.TaskService = service.TaskServiceSqlLite{
			CropQuery:      cropQuery,
			AreaQuery:      areaQuery,
			MaterialQuery:  materialReadQuery,
			ReservoirQuery: reservoirQuery,
		}

	case config.DB_SQLITE:
		taskServer.TaskEventRepo = repoSqlite.NewTaskEventRepositorySqlite(db)
		taskServer.TaskReadRepo = repoSqlite.NewTaskReadRepositorySqlite(db)

		taskServer.TaskEventQuery = querySqlite.NewTaskEventQuerySqlite(db)
		taskServer.TaskReadQuery = querySqlite.NewTaskReadQuerySqlite(db)

		cropQuery := querySqlite.NewCropQuerySqlite(db)
		areaQuery := querySqlite.NewAreaQuerySqlite(db)
		materialReadQuery := querySqlite.NewMaterialQuerySqlite(db)
		reservoirQuery := querySqlite.NewReservoirQuerySqlite(db)

		taskServer.TaskService = service.TaskServiceSqlLite{
			CropQuery:      cropQuery,
			AreaQuery:      areaQuery,
			MaterialQuery:  materialReadQuery,
			ReservoirQuery: reservoirQuery,
		}

	case config.DB_MYSQL:
		taskServer.TaskEventRepo = repoMysql.NewTaskEventRepositoryMysql(db)
		taskServer.TaskReadRepo = repoMysql.NewTaskReadRepositoryMysql(db)

		taskServer.TaskEventQuery = queryMysql.NewTaskEventQueryMysql(db)
		taskServer.TaskReadQuery = queryMysql.NewTaskReadQueryMysql(db)

		cropQuery := queryMysql.NewCropQueryMysql(db)
		areaQuery := queryMysql.NewAreaQueryMysql(db)
		materialReadQuery := queryMysql.NewMaterialQueryMysql(db)
		reservoirQuery := queryMysql.NewReservoirQueryMysql(db)

		taskServer.TaskService = service.TaskServiceSqlLite{
			CropQuery:      cropQuery,
			AreaQuery:      areaQuery,
			MaterialQuery:  materialReadQuery,
			ReservoirQuery: reservoirQuery,
		}

	}

	taskServer.InitSubscriber()

	return taskServer, nil
}

// InitSubscriber defines the mapping of which event this domain listen with their handler
func (s *TaskServer) InitSubscriber() {
	s.EventBus.Subscribe(domain.TaskCreatedCode, s.SaveToTaskReadModel)
	s.EventBus.Subscribe(domain.TaskTitleChangedCode, s.SaveToTaskReadModel)
	s.EventBus.Subscribe(domain.TaskDescriptionChangedCode, s.SaveToTaskReadModel)
	s.EventBus.Subscribe(domain.TaskPriorityChangedCode, s.SaveToTaskReadModel)
	s.EventBus.Subscribe(domain.TaskDueDateChangedCode, s.SaveToTaskReadModel)
	s.EventBus.Subscribe(domain.TaskCategoryChangedCode, s.SaveToTaskReadModel)
	s.EventBus.Subscribe(domain.TaskDetailsChangedCode, s.SaveToTaskReadModel)
	s.EventBus.Subscribe(domain.TaskCancelledCode, s.SaveToTaskReadModel)
	s.EventBus.Subscribe(domain.TaskCompletedCode, s.SaveToTaskReadModel)
	s.EventBus.Subscribe(domain.TaskDueCode, s.SaveToTaskReadModel)
}

// Mount defines the TaskServer's endpoints with its handlers
func (s *TaskServer) Mount(g *echo.Group) {
	g.POST("", s.SaveTask)

	g.GET("", s.FindAllTasks)
	g.GET("/search", s.FindFilteredTasks)
	g.GET("/:id", s.FindTaskByID)
	g.PUT("/:id", s.UpdateTask)
	g.PUT("/:id/cancel", s.CancelTask)
	g.PUT("/:id/complete", s.CompleteTask)
	// As we don't have an async task right now to check for Due state,
	// I'm adding a rest call to be able to manually do that. We can remove it in the future
	g.PUT("/:id/due", s.SetTaskAsDue)
}

func (s TaskServer) FindAllTasks(c echo.Context) error {
	data := make(map[string]interface{})

	page := c.QueryParam("page")
	limit := c.QueryParam("limit")

	pageInt, limitInt, err := paginationhelper.ParsePagination(page, limit)
	if err != nil {
		return Error(c, err)
	}

	result := <-s.TaskReadQuery.FindAll(pageInt, limitInt)
	if result.Error != nil {
		return result.Error
	}

	tasks, ok := result.Result.([]storage.TaskRead)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	taskList := []storage.TaskRead{}
	for _, v := range tasks {
		s.AppendTaskDomainDetails(&v)
		taskList = append(taskList, v)
	}

	// Return list of tasks
	data["data"] = taskList
	// Return number of tasks
	countResult := <-s.TaskReadQuery.CountAll()

	if countResult.Error != nil {
		return countResult.Error
	}
	count, ok := countResult.Result.(int)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}
	data["total_rows"] = count
	data["page"] = pageInt

	return c.JSON(http.StatusOK, data)
}

func (s TaskServer) FindFilteredTasks(c echo.Context) error {
	data := make(map[string]interface{})

	queryparams := make(map[string]string)
	queryparams["is_due"] = c.QueryParam("is_due")
	queryparams["priority"] = c.QueryParam("priority")
	queryparams["status"] = c.QueryParam("status")
	queryparams["domain"] = c.QueryParam("domain")
	queryparams["asset_id"] = c.QueryParam("asset_id")
	queryparams["category"] = c.QueryParam("category")
	queryparams["due_start"] = c.QueryParam("due_start")
	queryparams["due_end"] = c.QueryParam("due_end")

	page := c.QueryParam("page")
	limit := c.QueryParam("limit")

	pageInt, limitInt, err := paginationhelper.ParsePagination(page, limit)
	if err != nil {
		return Error(c, err)
	}

	result := <-s.TaskReadQuery.FindTasksWithFilter(queryparams, pageInt, limitInt)
	if result.Error != nil {
		return result.Error
	}

	tasks, ok := result.Result.([]storage.TaskRead)

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	taskList := []storage.TaskRead{}
	for _, v := range tasks {
		s.AppendTaskDomainDetails(&v)
		taskList = append(taskList, v)
	}

	// Return list of tasks
	data["data"] = taskList
	// Return number of tasks
	countResult := <-s.TaskReadQuery.CountTasksWithFilter(queryparams)

	if countResult.Error != nil {
		return countResult.Error
	}
	count, ok := countResult.Result.(int)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}
	data["total_rows"] = count
	data["page"] = pageInt
	return c.JSON(http.StatusOK, data)
}

// SaveTask is a TaskServer's handler to save new Task
func (s *TaskServer) SaveTask(c echo.Context) error {

	data := make(map[string]storage.TaskRead)

	form_date := c.FormValue("due_date")
	due_ptr := (*time.Time)(nil)
	if len(form_date) != 0 {
		due_date, err := time.Parse(time.RFC3339Nano, form_date)

		if err != nil {
			return Error(c, err)
		}
		due_ptr = &due_date
	}

	asset_id := c.FormValue("asset_id")
	asset_id_ptr := (*uuid.UUID)(nil)
	if len(asset_id) != 0 {
		asset_id, err := uuid.FromString(asset_id)
		if err != nil {
			return Error(c, err)
		}
		asset_id_ptr = &asset_id
	}

	domaincode := c.FormValue("domain")

	domaintask, err := s.CreateTaskDomainByCode(domaincode, c)

	if err != nil {
		return Error(c, err)
	}

	task, err := domain.CreateTask(
		s.TaskService,
		c.FormValue("title"),
		c.FormValue("description"),
		due_ptr,
		c.FormValue("priority"),
		domaintask,
		c.FormValue("category"),
		asset_id_ptr)

	if err != nil {
		return Error(c, err)
	}

	err = <-s.TaskEventRepo.Save(task.UID, 0, task.UncommittedChanges)
	if err != nil {
		return Error(c, err)
	}

	// Trigger Events
	s.publishUncommittedEvents(task)

	taskRead := MapTaskToTaskRead(task)
	s.AppendTaskDomainDetails(taskRead)

	data["data"] = *taskRead

	return c.JSON(http.StatusOK, data)
}

func (s *TaskServer) CreateTaskDomainByCode(domaincode string, c echo.Context) (domain.TaskDomain, error) {
	domainvalue := domaincode
	if domainvalue == "" {
		return nil, NewRequestValidationError(REQUIRED, "domain")
	}

	switch domainvalue {
	case domain.TaskDomainAreaCode:

		category := c.FormValue("category")
		materialID := c.FormValue("material_id")

		materialPtr := (*uuid.UUID)(nil)
		if len(materialID) != 0 {
			uid, err := uuid.FromString(materialID)
			if err != nil {
				return domain.TaskDomainArea{}, err
			}
			materialPtr = &uid
		}

		return domain.CreateTaskDomainArea(s.TaskService, category, materialPtr)
	case domain.TaskDomainCropCode:

		category := c.FormValue("category")
		materialID := c.FormValue("material_id")
		areaID := c.FormValue("area_id")

		materialPtr := (*uuid.UUID)(nil)
		if len(materialID) != 0 {
			uid, err := uuid.FromString(materialID)
			if err != nil {
				return domain.TaskDomainCrop{}, err
			}
			materialPtr = &uid
		}

		areaPtr := (*uuid.UUID)(nil)
		if len(areaID) != 0 {
			uid, err := uuid.FromString(areaID)
			if err != nil {
				return domain.TaskDomainCrop{}, err
			}
			areaPtr = &uid
		}

		return domain.CreateTaskDomainCrop(s.TaskService, category, materialPtr, areaPtr)
	case domain.TaskDomainFinanceCode:
		return domain.CreateTaskDomainFinance()
	case domain.TaskDomainGeneralCode:
		return domain.CreateTaskDomainGeneral()
	case domain.TaskDomainInventoryCode:
		return domain.CreateTaskDomainInventory()
	case domain.TaskDomainReservoirCode:

		category := c.FormValue("category")
		materialID := c.FormValue("material_id")

		materialPtr := (*uuid.UUID)(nil)
		if len(materialID) != 0 {
			uid, err := uuid.FromString(materialID)
			if err != nil {
				return domain.TaskDomainReservoir{}, err
			}
			materialPtr = &uid
		}
		return domain.CreateTaskDomainReservoir(s.TaskService, category, materialPtr)
	default:
		return nil, NewRequestValidationError(INVALID_OPTION, "domain")
	}
}

func (s *TaskServer) FindTaskByID(c echo.Context) error {
	data := make(map[string]storage.TaskRead)
	uid, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	result := <-s.TaskReadQuery.FindByID(uid)
	if result.Error != nil {
		return Error(c, result.Error)
	}

	task, ok := result.Result.(storage.TaskRead)

	if task.UID != uid {
		return Error(c, NewRequestValidationError(NOT_FOUND, "id"))
	}
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	s.AppendTaskDomainDetails(&task)

	data["task"] = task
	return c.JSON(http.StatusOK, data)
}

func (s *TaskServer) AppendTaskDomainDetails(task *storage.TaskRead) error {

	switch task.Domain {
	case domain.TaskDomainAreaCode:
		materialID := task.DomainDetails.(domain.TaskDomainArea).MaterialID
		if materialID != nil {
			materialResult := s.TaskService.FindMaterialByID(*materialID)
			materialQueryResult, ok := materialResult.Result.(query.TaskMaterialQueryResult)

			if !ok {
				return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
			}
			task.DomainDetails = &storage.TaskDomainDetailedArea{
				MaterialID:           &materialQueryResult.UID,
				MaterialName:         materialQueryResult.Name,
				MaterialType:         materialQueryResult.TypeCode,
				MaterialDetailedType: materialQueryResult.DetailedTypeCode,
			}
		}

	case domain.TaskDomainCropCode:
		material := (*storage.TaskDomainCropMaterial)(nil)
		area := (*storage.TaskDomainCropArea)(nil)
		crop := (*storage.TaskDomainCropBatch)(nil)

		materialID := task.DomainDetails.(domain.TaskDomainCrop).MaterialID
		if materialID != nil {
			materialResult := s.TaskService.FindMaterialByID(*materialID)
			materialQueryResult, ok := materialResult.Result.(query.TaskMaterialQueryResult)

			if !ok {
				return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
			}
			material = &storage.TaskDomainCropMaterial{
				MaterialID:           &materialQueryResult.UID,
				MaterialName:         materialQueryResult.Name,
				MaterialType:         materialQueryResult.TypeCode,
				MaterialDetailedType: materialQueryResult.DetailedTypeCode,
			}
		}
		areaID := task.DomainDetails.(domain.TaskDomainCrop).AreaID
		if areaID != nil {
			areaResult := s.TaskService.FindAreaByID(*areaID)
			areaQueryResult, ok := areaResult.Result.(query.TaskAreaQueryResult)

			if !ok {
				return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
			}
			area = &storage.TaskDomainCropArea{
				AreaID:   &areaQueryResult.UID,
				AreaName: areaQueryResult.Name,
			}
		}
		cropID := task.AssetID
		if cropID != nil {
			cropResult := s.TaskService.FindCropByID(*cropID)
			cropQueryResult, ok := cropResult.Result.(query.TaskCropQueryResult)

			if !ok {
				return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
			}
			crop = &storage.TaskDomainCropBatch{
				CropID:      &cropQueryResult.UID,
				CropBatchID: cropQueryResult.BatchID,
			}
		}
		task.DomainDetails = &storage.TaskDomainDetailedCrop{
			Material: material,
			Area:     area,
			Crop:     crop,
		}
	case domain.TaskDomainReservoirCode:

		materialID := task.DomainDetails.(domain.TaskDomainReservoir).MaterialID
		if materialID != nil {
			materialResult := s.TaskService.FindMaterialByID(*materialID)
			materialQueryResult, ok := materialResult.Result.(query.TaskMaterialQueryResult)

			if !ok {
				return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
			}
			task.DomainDetails = &storage.TaskDomainDetailedReservoir{
				MaterialID:           &materialQueryResult.UID,
				MaterialName:         materialQueryResult.Name,
				MaterialType:         materialQueryResult.TypeCode,
				MaterialDetailedType: materialQueryResult.DetailedTypeCode,
			}
		}

	}
	return nil
}

func (s *TaskServer) UpdateTask(c echo.Context) error {

	data := make(map[string]storage.TaskRead)
	uid, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	// Get TaskRead By UID
	readResult := <-s.TaskReadQuery.FindByID(uid)

	taskRead, ok := readResult.Result.(storage.TaskRead)

	if taskRead.UID != uid {
		return Error(c, NewRequestValidationError(NOT_FOUND, "id"))
	}
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	// Get TaskEvent under Task UID
	eventQueryResult := <-s.TaskEventQuery.FindAllByTaskID(uid)
	events := eventQueryResult.Result.([]storage.TaskEvent)

	// Build TaskEvents from history
	task := repository.BuildTaskFromEventHistory(s.TaskService, events)

	updatedTask, err := s.updateTaskAttributes(s.TaskService, task, c)
	if err != nil {
		return Error(c, err)
	}

	// Save new TaskEvent
	err = <-s.TaskEventRepo.Save(updatedTask.UID, updatedTask.Version, updatedTask.UncommittedChanges)
	if err != nil {
		return Error(c, err)
	}

	// Trigger Events
	s.publishUncommittedEvents(updatedTask)
	read := MapTaskToTaskRead(updatedTask)

	s.AppendTaskDomainDetails(read)

	data["data"] = *read

	return c.JSON(http.StatusOK, data)
}

func (s *TaskServer) updateTaskAttributes(taskService domain.TaskService, task *domain.Task, c echo.Context) (*domain.Task, error) {

	// Change Task Title
	title := c.FormValue("title")
	if len(title) != 0 {
		task.ChangeTaskTitle(s.TaskService, title)
	}

	// Change Task Description
	description := c.FormValue("description")
	if len(description) != 0 {
		task.ChangeTaskDescription(s.TaskService, description)
	}

	// Change Task Due Date
	form_date := c.FormValue("due_date")
	due_ptr := (*time.Time)(nil)
	if len(form_date) != 0 {

		due_date, err := time.Parse(time.RFC3339Nano, form_date)

		if err != nil {
			return task, Error(c, err)
		}
		due_ptr = &due_date
		task.ChangeTaskDueDate(s.TaskService, due_ptr)
	}

	// Change Task Priority
	priority := c.FormValue("priority")
	if len(priority) != 0 {
		task.ChangeTaskPriority(s.TaskService, priority)
	}

	// Change Task Category & Domain Details
	category := c.FormValue("category")
	if len(category) != 0 {
		task.ChangeTaskCategory(s.TaskService, category)
		details, err := s.CreateTaskDomainByCode(task.Domain, c)

		if err != nil {
			return &domain.Task{}, Error(c, err)
		}
		task.ChangeTaskDetails(s.TaskService, details)
	}

	return task, nil
}

func (s *TaskServer) CancelTask(c echo.Context) error {

	data := make(map[string]storage.TaskRead)
	uid, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	// Get TaskRead By UID
	readResult := <-s.TaskReadQuery.FindByID(uid)

	taskRead, ok := readResult.Result.(storage.TaskRead)

	if taskRead.UID != uid {
		return Error(c, NewRequestValidationError(NOT_FOUND, "id"))
	}
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	// Get TaskEvent under Task UID
	eventQueryResult := <-s.TaskEventQuery.FindAllByTaskID(uid)
	events := eventQueryResult.Result.([]storage.TaskEvent)

	// Build TaskEvents from history
	task := repository.BuildTaskFromEventHistory(s.TaskService, events)

	updatedTask, err := s.updateTaskAttributes(s.TaskService, task, c)
	if err != nil {
		return Error(c, err)
	}

	updatedTask.CancelTask(s.TaskService)

	// Save new TaskEvent
	err = <-s.TaskEventRepo.Save(updatedTask.UID, updatedTask.Version, updatedTask.UncommittedChanges)
	if err != nil {
		return Error(c, err)
	}

	// Trigger Events
	s.publishUncommittedEvents(updatedTask)

	read := MapTaskToTaskRead(updatedTask)

	s.AppendTaskDomainDetails(read)

	data["data"] = *read

	return c.JSON(http.StatusOK, data)
}

func (s *TaskServer) CompleteTask(c echo.Context) error {

	data := make(map[string]storage.TaskRead)
	uid, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	// Get TaskRead By UID
	readResult := <-s.TaskReadQuery.FindByID(uid)

	taskRead, ok := readResult.Result.(storage.TaskRead)

	if taskRead.UID != uid {
		return Error(c, NewRequestValidationError(NOT_FOUND, "id"))
	}
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	// Get TaskEvent under Task UID
	eventQueryResult := <-s.TaskEventQuery.FindAllByTaskID(uid)
	if eventQueryResult.Error != nil {
		return Error(c, eventQueryResult.Error)
	}

	events, ok := eventQueryResult.Result.([]storage.TaskEvent)
	if !ok {
		return Error(c, echo.NewHTTPError(http.StatusBadRequest, "Internal server error"))
	}

	// Build TaskEvents from history
	task := repository.BuildTaskFromEventHistory(s.TaskService, events)

	updatedTask, err := s.updateTaskAttributes(s.TaskService, task, c)
	if err != nil {
		return Error(c, err)
	}

	updatedTask.CompleteTask(s.TaskService)

	// Save new TaskEvent
	err = <-s.TaskEventRepo.Save(updatedTask.UID, updatedTask.Version, updatedTask.UncommittedChanges)
	if err != nil {
		return Error(c, err)
	}

	// Trigger Events
	s.publishUncommittedEvents(updatedTask)
	read := MapTaskToTaskRead(updatedTask)

	s.AppendTaskDomainDetails(read)

	data["data"] = *read

	return c.JSON(http.StatusOK, data)
}

func (s *TaskServer) SetTaskAsDue(c echo.Context) error {

	data := make(map[string]storage.TaskRead)
	uid, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	// Get TaskRead By UID
	readResult := <-s.TaskReadQuery.FindByID(uid)

	taskRead, ok := readResult.Result.(storage.TaskRead)

	if taskRead.UID != uid {
		return Error(c, NewRequestValidationError(NOT_FOUND, "id"))
	}
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	// Get TaskEvent under Task UID
	eventQueryResult := <-s.TaskEventQuery.FindAllByTaskID(uid)
	events := eventQueryResult.Result.([]storage.TaskEvent)

	// Build TaskEvents from history
	task := repository.BuildTaskFromEventHistory(s.TaskService, events)

	task.SetTaskAsDue(s.TaskService)

	// Save new TaskEvent
	err = <-s.TaskEventRepo.Save(task.UID, task.Version, task.UncommittedChanges)
	if err != nil {
		return Error(c, err)
	}

	// Trigger Events
	s.publishUncommittedEvents(task)

	read := MapTaskToTaskRead(task)

	s.AppendTaskDomainDetails(read)

	data["data"] = *read

	return c.JSON(http.StatusOK, data)
}

func (s *TaskServer) publishUncommittedEvents(entity interface{}) error {

	switch e := entity.(type) {
	case *domain.Task:
		for _, v := range e.UncommittedChanges {
			name := structhelper.GetName(v)
			s.EventBus.Publish(name, v)
		}
	default:
	}

	return nil
}
