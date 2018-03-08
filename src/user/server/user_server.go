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

// UserServer ties the routes and handlers with injected dependencies
type UserServer struct {
	UserEventRepo  repository.UserEventRepository
	UserReadRepo   repository.UserReadRepository
	UserEventQuery query.UserEventQuery
	UserReadQuery  query.UserReadQuery
	UserAuthRepo   repository.UserAuthRepository
	UserAuthQuery  query.UserAuthQuery
	UserService    domain.UserService
	EventBus       EventBus.Bus
}

// NewUserServer initializes UserServer's dependencies and create new UserServer struct
func NewUserServer(
	db *sql.DB,
	eventBus EventBus.Bus,
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
	g.POST("/:id/change_password", s.ChangePassword)
}

func (s *UserServer) Authorize(c echo.Context) error {
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

func (s *UserServer) ChangePassword(c echo.Context) error {
	id := c.Param("id")
	oldPassword := c.FormValue("old_password")
	newPassword := c.FormValue("new_password")
	confirmNewPassword := c.FormValue("confirm_new_password")

	userUID, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	queryResult := <-s.UserReadQuery.FindByID(userUID)
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
