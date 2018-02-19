package server

import (
	"net/http"
	"time"

	assetsstorage "github.com/Tanibox/tania-server/src/assets/storage"
	cropstorage "github.com/Tanibox/tania-server/src/growth/storage"
	"github.com/Tanibox/tania-server/src/helper/structhelper"
	"github.com/Tanibox/tania-server/src/tasks/domain"
	service "github.com/Tanibox/tania-server/src/tasks/domain/service"
	"github.com/Tanibox/tania-server/src/tasks/query/inmemory"
	"github.com/Tanibox/tania-server/src/tasks/repository"
	"github.com/Tanibox/tania-server/src/tasks/storage"
	"github.com/asaskevich/EventBus"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// TaskServer ties the routes and handlers with injected dependencies
type TaskServer struct {
	TaskRepo      repository.TaskRepository
	TaskEventRepo repository.TaskEventRepository
	TaskReadRepo  repository.TaskReadRepository
	TaskService   domain.TaskService
	EventBus      EventBus.Bus
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
	return &TaskServer{
		TaskRepo:      taskRepo,
		TaskEventRepo: taskEventRepo,
		TaskReadRepo:  taskReadRepo,
		TaskService:   taskService,
		EventBus:      bus,
	}, nil
}

// InitSubscriber defines the mapping of which event this domain listen with their handler
func (s *TaskServer) InitSubscriber() {
	s.EventBus.Subscribe(domain.TaskCreatedCode, s.SaveToTaskReadModel)
}

// Mount defines the TaskServer's endpoints with its handlers
func (s *TaskServer) Mount(g *echo.Group) {
	g.POST("", s.SaveTask)

	g.GET("", s.FindAllTask)
	g.GET("/search", s.FindFilteredTasks)
	g.GET("/:id", s.FindTaskByID)
	g.PUT("/:id", s.UpdateTask)
	g.PUT("/:id/cancel", s.CancelTask)
	g.PUT("/:id/complete", s.CompleteTask)
	// As we don't have an async task right now to check for Due state,
	// I'm adding a rest call to be able to manually do that. We can remove it in the future
	g.PUT("/:id/due", s.SetTaskAsDue)
}

func (s TaskServer) FindAllTask(c echo.Context) error {
	data := make(map[string][]SimpleTask)

	result := <-s.TaskRepo.FindAll()
	if result.Error != nil {
		return result.Error
	}

	Tasks, ok := result.Result.([]domain.Task)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data["data"] = MapToSimpleTask(Tasks)

	return c.JSON(http.StatusOK, data)
}

func (s TaskServer) FindFilteredTasks(c echo.Context) error {
	data := make(map[string][]SimpleTask)

	queryparams := make(map[string]string)
	queryparams["is_due"] = c.QueryParam("is_due")
	queryparams["priority"] = c.QueryParam("priority")
	queryparams["status"] = c.QueryParam("status")
	queryparams["domain"] = c.QueryParam("domain")
	queryparams["asset_id"] = c.QueryParam("asset_id")

	result := <-s.TaskRepo.FindTasksWithFilter(queryparams)

	if result.Error != nil {
		return result.Error
	}

	Tasks, ok := result.Result.([]domain.Task)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data["data"] = MapToSimpleTask(Tasks)

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

	/*err = <-s.TaskRepo.Save(task)
	if err != nil {
		return Error(c, err)
	}
	*/
	err = <-s.TaskEventRepo.Save(task.UID, 0, task.UncommittedChanges)
	if err != nil {
		return Error(c, err)
	}

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
	data := make(map[string]domain.Task)
	uid, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	result := <-s.TaskRepo.FindByID(c.Param("id"))
	if result.Error != nil {
		return result.Error
	}

	task, ok := result.Result.(domain.Task)

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

	result := <-s.TaskRepo.FindByID(c.Param("id"))
	if result.Error != nil {
		return result.Error
	}

	task, ok := result.Result.(domain.Task)

	if task.UID != uid {
		return Error(c, NewRequestValidationError(NOT_FOUND, "id"))
	}
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	updated_task, err := s.updateTaskAttributes(task, c)
	if err != nil {
		return Error(c, err)
	}

	err = <-s.TaskRepo.Save(&updated_task)
	if err != nil {
		return Error(c, err)
	}

	data["data"] = updated_task

	return c.JSON(http.StatusOK, data)
}

func (s *TaskServer) updateTaskAttributes(task domain.Task, c echo.Context) (domain.Task, error) {

	// Change Task Title
	title := c.FormValue("title")
	if len(title) != 0 {
		err := task.ChangeTaskTitle(title)
		if err != nil {
			return task, Error(c, err)
		}
	}

	// Change Task Description

	description := c.FormValue("description")
	if len(description) != 0 {
		err := task.ChangeTaskDescription(description)
		if err != nil {
			return task, Error(c, err)
		}
	}

	// Change Task Due Date

	form_date := c.FormValue("due_date")
	due_ptr := (*time.Time)(nil)
	if len(form_date) != 0 {
		due_date, err := time.Parse(time.RFC3339, form_date)

		if err != nil {
			return task, Error(c, err)
		}
		due_ptr = &due_date
		task.ChangeTaskDueDate(due_ptr)
	}

	// Change Task Priority

	priority := c.FormValue("priority")
	if len(priority) != 0 {
		err := task.ChangeTaskPriority(priority)
		if err != nil {
			return task, Error(c, err)
		}
	}

	// Change Task Asset

	asset_id := c.FormValue("asset_id")
	asset_id_ptr := (*uuid.UUID)(nil)
	if len(asset_id) != 0 {
		asset_id, err := uuid.FromString(asset_id)
		if err != nil {
			return task, Error(c, err)
		}
		asset_id_ptr = &asset_id
		task.ChangeTaskAssetID(s.TaskService, asset_id_ptr)
	}

	// Change Task Category & Domain Details

	category := c.FormValue("category")
	if len(category) != 0 {
		err := task.ChangeTaskCategory(category)
		if err != nil {
			return task, Error(c, err)
		}
		// Change Domain Details

		inventory_id := c.FormValue("inventory_id")
		if len(inventory_id) != 0 {
			details, err := s.CreateTaskDomainByCode(task.Domain, c)
			if err != nil {
				return task, Error(c, err)
			}
			task.ChangeTaskDomainDetails(details)
		}
	}
	return task, nil
}

func (s *TaskServer) CancelTask(c echo.Context) error {

	data := make(map[string]domain.Task)
	uid, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	result := <-s.TaskRepo.FindByID(c.Param("id"))
	if result.Error != nil {
		return result.Error
	}

	task, ok := result.Result.(domain.Task)

	if task.UID != uid {
		return Error(c, NewRequestValidationError(NOT_FOUND, "id"))
	}
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	updated_task, err := s.updateTaskAttributes(task, c)
	if err != nil {
		return Error(c, err)
	}

	updated_task.ChangeTaskStatus(domain.TaskStatusCancelled)

	err = <-s.TaskRepo.Save(&updated_task)
	if err != nil {
		return Error(c, err)
	}

	data["data"] = updated_task

	return c.JSON(http.StatusOK, data)
}

func (s *TaskServer) CompleteTask(c echo.Context) error {

	data := make(map[string]domain.Task)
	uid, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	result := <-s.TaskRepo.FindByID(c.Param("id"))
	if result.Error != nil {
		return result.Error
	}

	task, ok := result.Result.(domain.Task)

	if task.UID != uid {
		return Error(c, NewRequestValidationError(NOT_FOUND, "id"))
	}
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	updated_task, err := s.updateTaskAttributes(task, c)
	if err != nil {
		return Error(c, err)
	}

	updated_task.ChangeTaskStatus(domain.TaskStatusComplete)
	updated_task.SetTaskCompletedDate()

	err = <-s.TaskRepo.Save(&updated_task)
	if err != nil {
		return Error(c, err)
	}

	data["data"] = updated_task

	return c.JSON(http.StatusOK, data)
}

func (s *TaskServer) SetTaskAsDue(c echo.Context) error {

	data := make(map[string]domain.Task)
	uid, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return Error(c, err)
	}

	result := <-s.TaskRepo.FindByID(c.Param("id"))
	if result.Error != nil {
		return result.Error
	}

	task, ok := result.Result.(domain.Task)

	if task.UID != uid {
		return Error(c, NewRequestValidationError(NOT_FOUND, "id"))
	}
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	updated_task, err := s.updateTaskAttributes(task, c)
	if err != nil {
		return Error(c, err)
	}

	updated_task.SetTaskAsDue()

	err = <-s.TaskRepo.Save(&updated_task)
	if err != nil {
		return Error(c, err)
	}

	data["data"] = updated_task

	return c.JSON(http.StatusOK, data)
}

func (s *TaskServer) publishUncommittedEvents(entity interface{}) error {
	switch e := entity.(type) {
	case *domain.Task:
		for _, v := range e.UncommittedChanges {
			name := structhelper.GetName(v)
			s.EventBus.Publish(name, v)
		}
	}

	return nil
}
