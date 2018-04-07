package server

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Tanibox/tania-core/src/eventbus"
	"github.com/Tanibox/tania-core/src/helper/structhelper"
	"github.com/Tanibox/tania-core/src/user/domain"
	"github.com/Tanibox/tania-core/src/user/domain/service"
	"github.com/Tanibox/tania-core/src/user/query"
	querySqlite "github.com/Tanibox/tania-core/src/user/query/sqlite"
	"github.com/Tanibox/tania-core/src/user/repository"
	repoSqlite "github.com/Tanibox/tania-core/src/user/repository/sqlite"
	"github.com/Tanibox/tania-core/src/user/storage"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// UserServer ties the routes and handlers with injected dependencies
type UserServer struct {
	UserEventRepo  repository.UserEventRepository
	UserReadRepo   repository.UserReadRepository
	UserEventQuery query.UserEventQuery
	UserReadQuery  query.UserReadQuery
	UserAuthRepo   repository.UserAuthRepository
	UserAuthQuery  query.UserAuthQuery
	UserService    domain.UserService
	EventBus       eventbus.TaniaEventBus
}

// NewUserServer initializes UserServer's dependencies and create new UserServer struct
func NewUserServer(
	db *sql.DB,
	eventBus eventbus.TaniaEventBus,
) (*UserServer, error) {
	userEventRepo := repoSqlite.NewUserEventRepositorySqlite(db)
	userReadRepo := repoSqlite.NewUserReadRepositorySqlite(db)
	userEventQuery := querySqlite.NewUserEventQuerySqlite(db)
	userReadQuery := querySqlite.NewUserReadQuerySqlite(db)

	userAuthRepo := repoSqlite.NewUserAuthRepositorySqlite(db)
	userAuthQuery := querySqlite.NewUserAuthQuerySqlite(db)

	userService := service.UserServiceImpl{UserReadQuery: userReadQuery}

	userServer := UserServer{
		UserEventRepo:  userEventRepo,
		UserReadRepo:   userReadRepo,
		UserEventQuery: userEventQuery,
		UserReadQuery:  userReadQuery,
		UserAuthRepo:   userAuthRepo,
		UserAuthQuery:  userAuthQuery,
		UserService:    userService,
		EventBus:       eventBus,
	}

	userServer.InitSubscriber()

	return &userServer, nil
}

// InitSubscriber defines the mapping of which event this domain listen with their handler
func (s *UserServer) InitSubscriber() {
	s.EventBus.Subscribe("PasswordChanged", s.SaveToUserReadModel)
}

// Mount defines the UserServer's endpoints with its handlers
func (s *UserServer) Mount(g *echo.Group) {
	g.POST("/change_password", s.ChangePassword)
}

func (s *UserServer) ChangePassword(c echo.Context) error {
	oldPassword := c.FormValue("old_password")
	newPassword := c.FormValue("new_password")
	confirmNewPassword := c.FormValue("confirm_new_password")

	// We only use one default user, which is `tania`, so it's okay to hardcoded it here
	queryResult := <-s.UserReadQuery.FindByUsername("tania")
	if queryResult.Error != nil {
		return queryResult.Error
	}

	userRead, ok := queryResult.Result.(storage.UserRead)
	if !ok {
		return Error(c, errors.New("Error type assertion"))
	}

	if userRead.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NOT_FOUND, "id"))
	}

	if newPassword != confirmNewPassword {
		return Error(c, NewRequestValidationError(NOT_MATCH, "password"))
	}

	// Process
	eventQueryResult := <-s.UserEventQuery.FindAllByID(userRead.UID)
	if eventQueryResult.Error != nil {
		return Error(c, eventQueryResult.Error)
	}

	events := eventQueryResult.Result.([]storage.UserEvent)
	user := repository.NewUserFromHistory(events)

	isValid, err := user.IsPasswordValid(oldPassword)
	if err != nil {
		return Error(c, err)
	}

	if !isValid {
		return Error(c, NewRequestValidationError(NOT_FOUND, "password"))
	}

	err = user.ChangePassword(oldPassword, newPassword, confirmNewPassword)
	if err != nil {
		return Error(c, err)
	}

	// Persists //
	resultSave := <-s.UserEventRepo.Save(user.UID, user.Version, user.UncommittedChanges)
	if resultSave != nil {
		return Error(c, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error"))
	}

	// Publish //
	s.publishUncommittedEvents(user)

	data := make(map[string]storage.UserRead)
	data["data"] = MapToUserRead(user)

	return c.JSON(http.StatusOK, data)

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
