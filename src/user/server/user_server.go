package server

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Tanibox/tania-server/src/helper/structhelper"
	"github.com/Tanibox/tania-server/src/user/domain"
	"github.com/Tanibox/tania-server/src/user/domain/service"
	"github.com/Tanibox/tania-server/src/user/query"
	querySqlite "github.com/Tanibox/tania-server/src/user/query/sqlite"
	"github.com/Tanibox/tania-server/src/user/repository"
	repoSqlite "github.com/Tanibox/tania-server/src/user/repository/sqlite"
	"github.com/asaskevich/EventBus"
	"github.com/labstack/echo"
)

// UserServer ties the routes and handlers with injected dependencies
type UserServer struct {
	UserEventRepo repository.UserEventRepository
	UserReadRepo  repository.UserReadRepository
	UserReadQuery query.UserReadQuery
	UserService   domain.UserService
	EventBus      EventBus.Bus
}

// NewUserServer initializes UserServer's dependencies and create new UserServer struct
func NewUserServer(
	db *sql.DB,
	eventBus EventBus.Bus,
) (*UserServer, error) {
	userEventRepo := repoSqlite.NewUserEventRepositorySqlite(db)
	userReadRepo := repoSqlite.NewUserReadRepositorySqlite(db)
	userReadQuery := querySqlite.NewUserReadQuerySqlite(db)

	userService := service.UserServiceImpl{UserReadQuery: userReadQuery}

	userServer := UserServer{
		UserEventRepo: userEventRepo,
		UserReadRepo:  userReadRepo,
		UserReadQuery: userReadQuery,
		UserService:   userService,
		EventBus:      eventBus,
	}

	userServer.InitSubscriber()

	return &userServer, nil
}

// InitSubscriber defines the mapping of which event this domain listen with their handler
func (s *UserServer) InitSubscriber() {
	s.EventBus.Subscribe("UserCreated", s.SaveToUserReadModel)
	s.EventBus.Subscribe("PasswordChanged", s.SaveToUserReadModel)
}

// Mount defines the UserServer's endpoints with its handlers
func (s *UserServer) Mount(g *echo.Group) {
	g.GET("login", s.Login)
	g.POST("register", s.Register)
}

func (s *UserServer) Login(c echo.Context) error {
	return c.JSON(http.StatusOK, "SIP LOGIN")
}

func (s *UserServer) Register(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")

	if password != confirmPassword {
		return errors.New("Confirm password didn't match")
	}

	user, err := s.RegisterNewUser(username, password, confirmPassword)
	if err != nil {
		return err
	}

	data := make(map[string]domain.User)
	data["data"] = *user

	return c.JSON(http.StatusOK, data)
}

func (s *UserServer) RegisterNewUser(username, password, confirmPassword string) (*domain.User, error) {
	user, err := domain.CreateUser(s.UserService, username, password, confirmPassword)
	if err != nil {
		return nil, err
	}

	err = <-s.UserEventRepo.Save(user.UID, user.Version, user.UncommittedChanges)
	if err != nil {
		return nil, err
	}

	s.publishUncommittedEvents(user)

	return user, nil
}

func (s *UserServer) publishUncommittedEvents(entity interface{}) error {
	switch e := entity.(type) {
	case *domain.User:
		for _, v := range e.UncommittedChanges {
			name := structhelper.GetName(v)
			s.EventBus.Publish(name, v)
		}
	}

	return nil
}
