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
	//fmt.Printf(c.FormValue("description"))
	//fmt.Printf(c.FormValue("duedate"))
	//fmt.Printf(c.FormValue("priority"))
	//fmt.Printf(c.FormValue("status"))
	//fmt.Printf(c.FormValue("type"))
	//fmt.Printf(c.FormValue("assetid"))

	duedate, err := time.Parse(time.RFC3339, c.FormValue("duedate"))

	if err != nil {
		return Error(c, err)
	}


	task, err := domain.CreateTask(
		c.FormValue("description"),
		duedate,
		c.FormValue("priority"),
		c.FormValue("status"),
		c.FormValue("type"),
		c.FormValue("assetid"))

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
