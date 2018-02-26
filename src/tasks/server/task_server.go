package server

import (
	"net/http"
	"time"

	assetsstorage "github.com/Tanibox/tania-server/src/assets/storage"
	cropstorage "github.com/Tanibox/tania-server/src/growth/storage"
	"github.com/Tanibox/tania-server/src/helper/structhelper"
	"github.com/Tanibox/tania-server/src/tasks/domain"
	service "github.com/Tanibox/tania-server/src/tasks/domain/service"
	"github.com/Tanibox/tania-server/src/tasks/query"
	"github.com/Tanibox/tania-server/src/tasks/query/inmemory"
	"github.com/Tanibox/tania-server/src/tasks/repository"
	"github.com/Tanibox/tania-server/src/tasks/storage"
	"github.com/asaskevich/EventBus"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// TaskServer ties the routes and handlers with injected dependencies
type TaskServer struct {
	TaskRepo       repository.TaskRepository
	TaskEventRepo  repository.TaskEventRepository
	TaskReadRepo   repository.TaskReadRepository
	TaskEventQuery query.TaskEventQuery
	TaskReadQuery  query.TaskReadQuery
	TaskService    domain.TaskService
	EventBus       EventBus.Bus
}

// NewTaskServer initializes TaskServer's dependencies and create new TaskServer struct
func NewTaskServer(
	bus EventBus.Bus,
	cropStorage *cropstorage.CropReadStorage,
	areaStorage *assetsstorage.AreaReadStorage,
	materialStorage *assetsstorage.MaterialReadStorage,
	reservoirStorage *assetsstorage.ReservoirReadStorage,
	taskEventStorage *storage.TaskEventStorage,
	taskReadStorage *storage.TaskReadStorage) (*TaskServer, error) {

	taskEventRepo := repository.NewTaskEventRepositoryInMemory(taskEventStorage)
	taskReadRepo := repository.NewTaskReadRepositoryInMemory(taskReadStorage)

	taskEventQuery := inmemory.NewTaskEventQueryInMemory(taskEventStorage)
	taskReadQuery := inmemory.NewTaskReadQueryInMemory(taskReadStorage)

	taskStorage := storage.TaskStorage{TaskMap: make(map[uuid.UUID]domain.Task)}
	taskRepo := repository.NewTaskRepositoryInMemory(&taskStorage)

	cropQuery := inmemory.NewCropQueryInMemory(cropStorage)
	areaQuery := inmemory.NewAreaQueryInMemory(areaStorage)
	materialReadQuery := inmemory.NewMaterialQueryInMemory(materialStorage)
	reservoirQuery := inmemory.NewReservoirQueryInMemory(reservoirStorage)

	taskService := service.TaskServiceInMemory{
		CropQuery:      cropQuery,
		AreaQuery:      areaQuery,
		MaterialQuery:  materialReadQuery,
		ReservoirQuery: reservoirQuery,
	}

	taskServer := &TaskServer{
		TaskRepo:       taskRepo,
		TaskEventRepo:  taskEventRepo,
		TaskReadRepo:   taskReadRepo,
		TaskEventQuery: taskEventQuery,
		TaskReadQuery:  taskReadQuery,
		TaskService:    taskService,
		EventBus:       bus,
	}

	taskServer.InitSubscriber()

	return taskServer, nil
}

// InitSubscriber defines the mapping of which event this domain listen with their handler
func (s *TaskServer) InitSubscriber() {
	s.EventBus.Subscribe(domain.TaskCreatedCode, s.SaveToTaskReadModel)
	s.EventBus.Subscribe(domain.TaskModifiedCode, s.SaveToTaskReadModel)
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
	data := make(map[string][]storage.TaskRead)

	result := <-s.TaskReadQuery.FindAll()
	if result.Error != nil {
		return result.Error
	}

	tasks, ok := result.Result.([]storage.TaskRead)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data["data"] = []storage.TaskRead{}
	for _, v := range tasks {
		data["data"] = append(data["data"], v)
	}

	return c.JSON(http.StatusOK, data)
}

func (s TaskServer) FindFilteredTasks(c echo.Context) error {
	data := make(map[string][]storage.TaskRead)

	queryparams := make(map[string]string)
	queryparams["is_due"] = c.QueryParam("is_due")
	queryparams["priority"] = c.QueryParam("priority")
	queryparams["status"] = c.QueryParam("status")
	queryparams["domain"] = c.QueryParam("domain")
	queryparams["asset_id"] = c.QueryParam("asset_id")

	result := <-s.TaskReadQuery.FindTasksWithFilter(queryparams)

	if result.Error != nil {
		return result.Error
	}

	tasks, ok := result.Result.([]storage.TaskRead)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data["data"] = []storage.TaskRead{}
	for _, v := range tasks {
		data["data"] = append(data["data"], v)
	}

	return c.JSON(http.StatusOK, data)
}

// SaveTask is a TaskServer's handler to save new Task
func (s *TaskServer) SaveTask(c echo.Context) error {

	data := make(map[string]domain.Task)

	form_date := c.FormValue("due_date")
	due_ptr := (*time.Time)(nil)
	if len(form_date) != 0 {
		due_date, err := time.Parse(time.RFC3339, form_date)

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

	data["data"] = *task

	return c.JSON(http.StatusOK, data)
}

func (s *TaskServer) CreateTaskDomainByCode(domaincode string, c echo.Context) (domain.TaskDomain, error) {
	domainvalue := domaincode
	if domainvalue == "" {
		return nil, NewRequestValidationError(REQUIRED, "domain")
	}

	switch domainvalue {
	case domain.TaskDomainAreaCode:
		return domain.CreateTaskDomainArea()
	case domain.TaskDomainCropCode:

		inv_id := c.FormValue("inventory_id")

		inventory_id_ptr := (*uuid.UUID)(nil)
		if len(inv_id) != 0 {
			inv_id, err := uuid.FromString(inv_id)
			if err != nil {
				return nil, Error(c, err)
			}
			inventory_id_ptr = &inv_id
		}
		return domain.CreateTaskDomainCrop(s.TaskService, inventory_id_ptr)
	case domain.TaskDomainFinanceCode:
		return domain.CreateTaskDomainFinance()
	case domain.TaskDomainGeneralCode:
		return domain.CreateTaskDomainGeneral()
	case domain.TaskDomainInventoryCode:
		return domain.CreateTaskDomainInventory()
	case domain.TaskDomainReservoirCode:
		return domain.CreateTaskDomainReservoir()
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

	task, ok := result.Result.(storage.TaskRead)

	if task.UID != uid {
		return Error(c, NewRequestValidationError(NOT_FOUND, "id"))
	}
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data["data"] = task

	return c.JSON(http.StatusOK, data)
}

func (s *TaskServer) UpdateTask(c echo.Context) error {

	data := make(map[string]domain.Task)
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
	task := repository.BuildTaskEventsFromHistory(s.TaskService, events)

	// Construct Task from TaskRead attributes
	updatedTask, err := taskRead.BuildTaskFromTaskRead(*task)

	if err != nil {
		return Error(c, err)
	}

	taskModifiedEvent, err := s.createTaskModifiedEvent(taskRead, c)
	if err != nil {
		return Error(c, err)
	}

	err = updatedTask.TrackChange(s.TaskService, *taskModifiedEvent)
	if err != nil {
		return Error(c, err)
	}

	// Save new TaskEvent
	err = <-s.TaskEventRepo.Save(updatedTask.UID, 0, updatedTask.UncommittedChanges)
	if err != nil {
		return Error(c, err)
	}

	// Trigger Events
	s.publishUncommittedEvents(updatedTask)

	data["data"] = *updatedTask

	return c.JSON(http.StatusOK, data)
}

func (s *TaskServer) createTaskModifiedEvent(taskRead storage.TaskRead, c echo.Context) (*domain.TaskModified, error) {

	// Change Task Title
	title := c.FormValue("title")
	if len(title) == 0 {
		title = taskRead.Title
	}

	// Change Task Description
	description := c.FormValue("description")
	if len(description) == 0 {
		description = taskRead.Description
	}

	// Change Task Due Date
	form_date := c.FormValue("due_date")
	due_ptr := (*time.Time)(nil)
	if len(form_date) == 0 {
		due_ptr = taskRead.DueDate
	}

	// Change Task Priority
	priority := c.FormValue("priority")
	if len(priority) == 0 {
		priority = taskRead.Priority
	}

	// Change Task Asset
	asset_id := c.FormValue("asset_id")
	asset_id_ptr := (*uuid.UUID)(nil)
	if len(asset_id) == 0 {
		asset_id_ptr = taskRead.AssetID
	} else {
		asset_id, err := uuid.FromString(asset_id)
		if err != nil {
			return &domain.TaskModified{}, Error(c, err)
		}
		asset_id_ptr = &asset_id
	}

	// Change Task Category & Domain Details
	var category string
	var details domain.TaskDomain
	var err error
	category = c.FormValue("category")
	if len(category) != 0 {
		inventory_id := c.FormValue("inventory_id")
		if len(inventory_id) != 0 {
			details, err = s.CreateTaskDomainByCode(taskRead.Domain, c)
			if err != nil {
				return &domain.TaskModified{}, Error(c, err)
			}
		} else {
			details = taskRead.DomainDetails
		}
	} else {
		category = taskRead.Category
		details = taskRead.DomainDetails
	}

	event, err := storage.CreateTaskModifiedEvent(taskRead.UID, title, description, due_ptr, priority, details, category, asset_id_ptr)

	if err != nil {
		return &domain.TaskModified{}, Error(c, err)
	}

	return event, nil
}

func (s *TaskServer) createTaskCancelledEvent(task *domain.Task) *domain.TaskCancelled {

	event := storage.CreateTaskCancelledEvent(task.UID, task.Title, task.Description, task.DueDate, task.Priority, task.DomainDetails, task.Category, task.AssetID)

	return event
}

func (s *TaskServer) createTaskCompletedEvent(task *domain.Task) *domain.TaskCompleted {

	event := storage.CreateTaskCompletedEvent(task.UID, task.Title, task.Description, task.DueDate, task.Priority, task.DomainDetails, task.Category, task.AssetID)

	return event
}

func (s *TaskServer) createTaskDueEvent(task *domain.Task) *domain.TaskDue {

	event := storage.CreateTaskDueEvent(task.UID, task.Title, task.Description, task.DueDate, task.Priority, task.DomainDetails, task.Category, task.AssetID)

	return event
}

func (s *TaskServer) CancelTask(c echo.Context) error {

	data := make(map[string]domain.Task)
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
	task := repository.BuildTaskEventsFromHistory(s.TaskService, events)

	// Construct Task from TaskRead attributes
	updatedTask, err := taskRead.BuildTaskFromTaskRead(*task)

	if err != nil {
		return Error(c, err)
	}

	// Create TaskModified event
	taskModifiedEvent, err := s.createTaskModifiedEvent(taskRead, c)
	if err != nil {
		return Error(c, err)
	}

	// Create TaskCancelled event
	taskCancelledEvent := s.createTaskCancelledEvent(updatedTask)

	// Track Changes
	err = updatedTask.TrackChange(s.TaskService, *taskModifiedEvent)
	if err != nil {
		return Error(c, err)
	}

	err = updatedTask.TrackChange(s.TaskService, *taskCancelledEvent)
	if err != nil {
		return Error(c, err)
	}

	// Save new TaskEvent
	err = <-s.TaskEventRepo.Save(updatedTask.UID, 0, updatedTask.UncommittedChanges)
	if err != nil {
		return Error(c, err)
	}

	// Trigger Events
	s.publishUncommittedEvents(updatedTask)

	data["data"] = *updatedTask

	return c.JSON(http.StatusOK, data)
}

func (s *TaskServer) CompleteTask(c echo.Context) error {

	data := make(map[string]domain.Task)
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
	task := repository.BuildTaskEventsFromHistory(s.TaskService, events)

	// Construct Task from TaskRead attributes
	updatedTask, err := taskRead.BuildTaskFromTaskRead(*task)

	if err != nil {
		return Error(c, err)
	}

	// Create TaskModified event
	taskModifiedEvent, err := s.createTaskModifiedEvent(taskRead, c)
	if err != nil {
		return Error(c, err)
	}

	// Create TaskCompleted event
	taskCompletedEvent := s.createTaskCompletedEvent(updatedTask)

	// Track Changes
	err = updatedTask.TrackChange(s.TaskService, *taskModifiedEvent)
	if err != nil {
		return Error(c, err)
	}

	err = updatedTask.TrackChange(s.TaskService, *taskCompletedEvent)
	if err != nil {
		return Error(c, err)
	}

	// Save new TaskEvent
	err = <-s.TaskEventRepo.Save(updatedTask.UID, 0, updatedTask.UncommittedChanges)
	if err != nil {
		return Error(c, err)
	}

	// Trigger Events
	s.publishUncommittedEvents(updatedTask)

	data["data"] = *updatedTask

	return c.JSON(http.StatusOK, data)
}

func (s *TaskServer) SetTaskAsDue(c echo.Context) error {

	data := make(map[string]domain.Task)
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
	task := repository.BuildTaskEventsFromHistory(s.TaskService, events)

	// Construct Task from TaskRead attributes
	updatedTask, err := taskRead.BuildTaskFromTaskRead(*task)

	if err != nil {
		return Error(c, err)
	}

	// Create TaskDue event
	taskDueEvent := s.createTaskDueEvent(updatedTask)

	// Track Changes
	err = updatedTask.TrackChange(s.TaskService, *taskDueEvent)
	if err != nil {
		return err
	}

	// Save new TaskEvent
	err = <-s.TaskEventRepo.Save(updatedTask.UID, 0, updatedTask.UncommittedChanges)
	if err != nil {
		return Error(c, err)
	}

	// Trigger Events
	s.publishUncommittedEvents(updatedTask)

	data["data"] = *updatedTask

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
