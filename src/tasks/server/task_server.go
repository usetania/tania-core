package server

import (
	"net/http"
	"github.com/labstack/echo"
	"time"
	domain "github.com/Tanibox/tania-server/src/tasks/domain"
	repository "github.com/Tanibox/tania-server/src/tasks/repository"
	storage "github.com/Tanibox/tania-server/src/tasks/storage"
	uuid "github.com/satori/go.uuid"

)

// TaskServer ties the routes and handlers with injected dependencies
type TaskServer struct {
	TaskRepo      repository.TaskRepository
}

// NewTaskServer initializes TaskServer's dependencies and create new TaskServer struct
func NewTaskServer() (*TaskServer, error) {

	taskStorage := storage.TaskStorage{TaskMap: make(map[uuid.UUID]domain.Task)}
	taskRepo := repository.NewTaskRepositoryInMemory(&taskStorage)
	return &TaskServer{
		TaskRepo:	taskRepo,
	}, nil
}

// Mount defines the TaskServer's endpoints with its handlers
func (s *TaskServer) Mount(g *echo.Group) {
	g.POST("", s.SaveTask)
	g.GET("", s.FindAllTask)
	g.GET("/:id", s.FindTaskByID)
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

// SaveTask is a TaskServer's handler to save new Task
func (s *TaskServer) SaveTask(c echo.Context) error {
	data := make(map[string]domain.Task)

	due_date, err := time.Parse(time.RFC3339, c.FormValue("due_date"))

	if err != nil {
		return Error(c, err)
	}


	task, err := domain.CreateTask(
		c.FormValue("description"),
		due_date,
		c.FormValue("priority"),
		c.FormValue("type"),
		c.FormValue("asset_id"))

	if err != nil {
		return Error(c, err)
	}

	err = <-s.TaskRepo.Save(&task)
	if err != nil {
		return Error(c, err)
	}

	data["data"] = task

	return c.JSON(http.StatusOK, data)
}

func (s *TaskServer) FindTaskByID(c echo.Context) error {
	data := make(map[string]domain.Task)

	result := <-s.TaskRepo.FindByID(c.Param("id"))
	if result.Error != nil {
		return result.Error
	}

	task, ok := result.Result.(domain.Task)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Internal server error")
	}

	data["data"] = task

	return c.JSON(http.StatusOK, data)
}
