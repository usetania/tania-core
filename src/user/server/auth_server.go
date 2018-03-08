package server

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/Tanibox/tania-server/config"
	"github.com/Tanibox/tania-server/src/helper/structhelper"
	"github.com/Tanibox/tania-server/src/user/domain"
	"github.com/Tanibox/tania-server/src/user/domain/service"
	"github.com/Tanibox/tania-server/src/user/query"
	querySqlite "github.com/Tanibox/tania-server/src/user/query/sqlite"
	"github.com/Tanibox/tania-server/src/user/repository"
	repoSqlite "github.com/Tanibox/tania-server/src/user/repository/sqlite"
	"github.com/Tanibox/tania-server/src/user/storage"
	"github.com/asaskevich/EventBus"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// AuthServer ties the routes and handlers with injected dependencies
type AuthServer struct {
	UserEventRepo  repository.UserEventRepository
	UserReadRepo   repository.UserReadRepository
	UserEventQuery query.UserEventQuery
	UserReadQuery  query.UserReadQuery
	UserAuthRepo   repository.UserAuthRepository
	UserAuthQuery  query.UserAuthQuery
	UserService    domain.UserService
	EventBus       EventBus.Bus
}

// NewAuthServer initializes AuthServer's dependencies and create new AuthServer struct
func NewAuthServer(
	db *sql.DB,
	eventBus EventBus.Bus,
) (*AuthServer, error) {
	userEventRepo := repoSqlite.NewUserEventRepositorySqlite(db)
	userReadRepo := repoSqlite.NewUserReadRepositorySqlite(db)
	userEventQuery := querySqlite.NewUserEventQuerySqlite(db)
	userReadQuery := querySqlite.NewUserReadQuerySqlite(db)

	userAuthRepo := repoSqlite.NewUserAuthRepositorySqlite(db)
	userAuthQuery := querySqlite.NewUserAuthQuerySqlite(db)

	userService := service.UserServiceImpl{UserReadQuery: userReadQuery}

	authServer := AuthServer{
		UserEventRepo:  userEventRepo,
		UserReadRepo:   userReadRepo,
		UserEventQuery: userEventQuery,
		UserReadQuery:  userReadQuery,
		UserAuthRepo:   userAuthRepo,
		UserAuthQuery:  userAuthQuery,
		UserService:    userService,
		EventBus:       eventBus,
	}

	authServer.InitSubscriber()

	return &authServer, nil
}

// InitSubscriber defines the mapping of which event this domain listen with their handler
func (s *AuthServer) InitSubscriber() {
	s.EventBus.Subscribe("UserCreated", s.SaveToUserReadModel)
}

// Mount defines the AuthServer's endpoints with its handlers
func (s *AuthServer) Mount(g *echo.Group) {
	g.POST("authorize", s.Authorize)
	g.POST("register", s.Register)
}

func (s *AuthServer) Authorize(c echo.Context) error {
	responseType := "token"
	redirectURI := *config.Config.RedirectURI

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
		return Error(c, errors.New("Error type assertion"))
	}

	if userRead.UID == (uuid.UUID{}) {
		return Error(c, errors.New("Invalid username or password"))
	}

	queryResult = <-s.UserAuthQuery.FindByUserID(userRead.UID)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	userAuth, ok := queryResult.Result.(storage.UserAuth)
	if !ok {
		return Error(c, errors.New("Error type assertion"))
	}

	if reqClientID != userAuth.ClientID {
		return Error(c, NewRequestValidationError(INVALID, "client_id"))
	}

	if reqRedirectURI != redirectURI {
		return Error(c, NewRequestValidationError(INVALID, "redirect_uri"))
	}

	if reqResponseType != responseType {
		return Error(c, NewRequestValidationError(INVALID, "response_type"))
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

	c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+accessToken)
	redirectURI += "#" + "access_token=" + accessToken + "&state=" + reqState + "&expires_in=" + strconv.Itoa(expiresIn)

	return c.Redirect(302, redirectURI)
}

func (s *AuthServer) Register(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")

	if password != confirmPassword {
		return Error(c, errors.New("Confirm password didn't match"))
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
// It is used by the register handler and in the initial user creation
func (s *AuthServer) RegisterNewUser(username, password, confirmPassword string) (*domain.User, *storage.UserAuth, error) {
	user, err := domain.CreateUser(s.UserService, username, password, confirmPassword)
	if err != nil {
		return nil, nil, err
	}

	err = <-s.UserEventRepo.Save(user.UID, user.Version, user.UncommittedChanges)
	if err != nil {
		return nil, nil, err
	}

	// Generate client ID
	clientID, err := uuid.NewV4()
	if err != nil {
		return nil, nil, err
	}

	userAuth := storage.UserAuth{
		UserUID:     user.UID,
		ClientID:    clientID.String(),
		CreatedDate: user.CreatedDate,
		LastUpdated: user.LastUpdated,
	}

	err = <-s.UserAuthRepo.Save(&userAuth)

	s.publishUncommittedEvents(user)

	return user, &userAuth, nil
}

func (s *AuthServer) publishUncommittedEvents(entity interface{}) error {
	switch e := entity.(type) {
	case *domain.User:
		for _, v := range e.UncommittedChanges {
			name := structhelper.GetName(v)
			s.EventBus.Publish(name, v)
		}
	}

	return nil
}
