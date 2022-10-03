package server

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/usetania/tania-core/config"
	"github.com/usetania/tania-core/src/eventbus"
	"github.com/usetania/tania-core/src/helper/structhelper"
	"github.com/usetania/tania-core/src/user/domain"
	"github.com/usetania/tania-core/src/user/domain/service"
	"github.com/usetania/tania-core/src/user/query"
	queryMysql "github.com/usetania/tania-core/src/user/query/mysql"
	querySqlite "github.com/usetania/tania-core/src/user/query/sqlite"
	"github.com/usetania/tania-core/src/user/repository"
	repoMysql "github.com/usetania/tania-core/src/user/repository/mysql"
	repoSqlite "github.com/usetania/tania-core/src/user/repository/sqlite"
	"github.com/usetania/tania-core/src/user/storage"
)

// UserServer ties the routes and handlers with injected dependencies.
type UserServer struct {
	UserEventRepo  repository.UserEvent
	UserReadRepo   repository.UserRead
	UserEventQuery query.UserEvent
	UserReadQuery  query.UserRead
	UserAuthRepo   repository.UserAuth
	UserAuthQuery  query.UserAuth
	UserService    domain.UserService
	EventBus       eventbus.TaniaEventBus
}

// NewUserServer initializes UserServer's dependencies and create new UserServer struct.
func NewUserServer(
	db *sql.DB,
	eventBus eventbus.TaniaEventBus,
) (*UserServer, error) {
	userServer := &UserServer{
		EventBus: eventBus,
	}

	switch *config.Config.TaniaPersistenceEngine {
	case config.DBSqlite:
		userServer.UserEventRepo = repoSqlite.NewUserEventRepositorySqlite(db)
		userServer.UserReadRepo = repoSqlite.NewUserReadRepositorySqlite(db)
		userServer.UserEventQuery = querySqlite.NewUserEventQuerySqlite(db)
		userServer.UserReadQuery = querySqlite.NewUserReadQuerySqlite(db)

		userServer.UserAuthRepo = repoSqlite.NewUserAuthRepositorySqlite(db)
		userServer.UserAuthQuery = querySqlite.NewUserAuthQuerySqlite(db)

		userServer.UserService = service.UserServiceImpl{UserReadQuery: userServer.UserReadQuery}

	case config.DBMysql:
		userServer.UserEventRepo = repoMysql.NewUserEventRepositoryMysql(db)
		userServer.UserReadRepo = repoMysql.NewUserReadRepositoryMysql(db)
		userServer.UserEventQuery = queryMysql.NewUserEventQueryMysql(db)
		userServer.UserReadQuery = queryMysql.NewUserReadQueryMysql(db)

		userServer.UserAuthRepo = repoMysql.NewUserAuthRepositoryMysql(db)
		userServer.UserAuthQuery = queryMysql.NewUserAuthQueryMysql(db)

		userServer.UserService = service.UserServiceImpl{UserReadQuery: userServer.UserReadQuery}
	}

	userServer.InitSubscriber()

	return userServer, nil
}

// InitSubscriber defines the mapping of which event this domain listen with their handler.
func (s *UserServer) InitSubscriber() {
	s.EventBus.Subscribe("PasswordChanged", s.SaveToUserReadModel)
}

// Mount defines the UserServer's endpoints with its handlers.
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
		return Error(c, errors.New("error type assertion"))
	}

	if userRead.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(NotFound, "id"))
	}

	if newPassword != confirmNewPassword {
		return Error(c, NewRequestValidationError(NorMatch, "password"))
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
		return Error(c, NewRequestValidationError(NotFound, "password"))
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

func (s *UserServer) publishUncommittedEvents(entity interface{}) {
	switch e := entity.(type) {
	case *domain.User:
		for _, v := range e.UncommittedChanges {
			name := structhelper.GetName(v)
			s.EventBus.Publish(name, v)
		}
	}
}
