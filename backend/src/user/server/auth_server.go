package server

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strconv"

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

// AuthServer ties the routes and handlers with injected dependencies.
type AuthServer struct {
	UserEventRepo  repository.UserEvent
	UserReadRepo   repository.UserRead
	UserEventQuery query.UserEvent
	UserReadQuery  query.UserRead
	UserAuthRepo   repository.UserAuth
	UserAuthQuery  query.UserAuth
	UserService    domain.UserService
	EventBus       eventbus.TaniaEventBus
}

// NewAuthServer initializes AuthServer's dependencies and create new AuthServer struct.
func NewAuthServer(
	db *sql.DB,
	eventBus eventbus.TaniaEventBus,
) (*AuthServer, error) {
	authServer := &AuthServer{
		EventBus: eventBus,
	}

	switch *config.Config.TaniaPersistenceEngine {
	case config.DBSqlite:
		authServer.UserEventRepo = repoSqlite.NewUserEventRepositorySqlite(db)
		authServer.UserReadRepo = repoSqlite.NewUserReadRepositorySqlite(db)
		authServer.UserEventQuery = querySqlite.NewUserEventQuerySqlite(db)
		authServer.UserReadQuery = querySqlite.NewUserReadQuerySqlite(db)

		authServer.UserAuthRepo = repoSqlite.NewUserAuthRepositorySqlite(db)
		authServer.UserAuthQuery = querySqlite.NewUserAuthQuerySqlite(db)

		authServer.UserService = service.UserServiceImpl{UserReadQuery: authServer.UserReadQuery}

	case config.DBMysql:
		authServer.UserEventRepo = repoMysql.NewUserEventRepositoryMysql(db)
		authServer.UserReadRepo = repoMysql.NewUserReadRepositoryMysql(db)
		authServer.UserEventQuery = queryMysql.NewUserEventQueryMysql(db)
		authServer.UserReadQuery = queryMysql.NewUserReadQueryMysql(db)

		authServer.UserAuthRepo = repoMysql.NewUserAuthRepositoryMysql(db)
		authServer.UserAuthQuery = queryMysql.NewUserAuthQueryMysql(db)

		authServer.UserService = service.UserServiceImpl{UserReadQuery: authServer.UserReadQuery}
	}

	authServer.InitSubscriber()

	return authServer, nil
}

// InitSubscriber defines the mapping of which event this domain listen with their handler.
func (s *AuthServer) InitSubscriber() {
	s.EventBus.Subscribe("UserCreated", s.SaveToUserReadModel)
}

// Mount defines the AuthServer's endpoints with its handlers.
func (s *AuthServer) Mount(g *echo.Group) {
	g.POST("authorize", s.Authorize)
	g.POST("register", s.Register)
}

func (s *AuthServer) Authorize(c echo.Context) error {
	responseType := "token"
	redirectURI := config.Config.RedirectURI
	clientID := *config.Config.ClientID

	reqUsername := c.FormValue("username")
	reqPassword := c.FormValue("password")
	reqClientID := c.FormValue("client_id")
	reqResponseType := c.FormValue("response_type")
	reqRedirectURI := c.FormValue("redirect_uri")
	reqState := c.FormValue("state")

	queryResult := <-s.UserReadQuery.FindByUsernameAndPassword(reqUsername, reqPassword)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	userRead, ok := queryResult.Result.(storage.UserRead)
	if !ok {
		return Error(c, errors.New("error type assertion"))
	}

	queryResult = <-s.UserAuthQuery.FindByUserID(userRead.UID)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	userAuth, ok := queryResult.Result.(storage.UserAuth)
	if !ok {
		return Error(c, errors.New("error type assertion"))
	}

	if userRead.UID == (uuid.UUID{}) {
		return Error(c, NewRequestValidationError(Invalid, "username"))
	}

	if reqClientID != clientID {
		return Error(c, NewRequestValidationError(Invalid, "client_id"))
	}

	if reqRedirectURI == "" {
		return Error(c, NewRequestValidationError(Required, "redirect_uri"))
	}

	var err error

	reqRedirectURI, err = url.PathUnescape(reqRedirectURI)
	if err != nil {
		return Error(c, err)
	}

	selectedRedirectURI := ""

	for _, v := range redirectURI {
		if reqRedirectURI == *v {
			selectedRedirectURI = *v

			break
		}
	}

	if selectedRedirectURI == "" {
		return Error(c, NewRequestValidationError(Invalid, "redirect_uri"))
	}

	if reqResponseType != responseType {
		return Error(c, NewRequestValidationError(Invalid, "response_type"))
	}

	// Generate access token here
	// We use uuid method temporarily until we find better method
	uidAccessToken, err := uuid.NewV4()
	if err != nil {
		return Error(c, err)
	}

	accessToken := uidAccessToken.String()

	// We don't expire token because it's complicating things
	// Also Google recommend it. https://developers.google.com/actions/identity/oauth2-implicit-flow
	expiresIn := 0

	userAuth.AccessToken = accessToken
	userAuth.TokenExpires = expiresIn

	err = <-s.UserAuthRepo.Save(&userAuth)
	if err != nil {
		return Error(c, err)
	}

	selectedRedirectURI += "?" + "access_token=" + accessToken + "&state=" + reqState + "&expires_in=" + strconv.Itoa(expiresIn) //nolint:lll

	c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+accessToken)

	return c.Redirect(302, selectedRedirectURI)
}

func (s *AuthServer) Register(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")

	if password != confirmPassword {
		return Error(c, errors.New("confirm password didn't match"))
	}

	user, _, err := s.RegisterNewUser(username, password, confirmPassword)
	if err != nil {
		return Error(c, err)
	}

	data := make(map[string]storage.UserRead)
	data["data"] = MapToUserRead(user)

	return c.JSON(http.StatusOK, data)
}

// RegisterNewUser is used to call the behaviour and persist it
// It is used by the register handler and in the initial user creation.
func (s *AuthServer) RegisterNewUser(
	username, password, confirmPassword string,
) (*domain.User, *storage.UserAuth, error) {
	user, err := domain.CreateUser(s.UserService, username, password, confirmPassword)
	if err != nil {
		return nil, nil, err
	}

	err = <-s.UserEventRepo.Save(user.UID, user.Version, user.UncommittedChanges)
	if err != nil {
		return nil, nil, err
	}

	userAuth := storage.UserAuth{
		UserUID:     user.UID,
		CreatedDate: user.CreatedDate,
		LastUpdated: user.LastUpdated,
	}

	err = <-s.UserAuthRepo.Save(&userAuth)
	if err != nil {
		return nil, nil, err
	}

	s.publishUncommittedEvents(user)

	return user, &userAuth, nil
}

func (s *AuthServer) publishUncommittedEvents(entity interface{}) {
	switch e := entity.(type) {
	case *domain.User:
		for _, v := range e.UncommittedChanges {
			name := structhelper.GetName(v)
			s.EventBus.Publish(name, v)
		}
	}
}
