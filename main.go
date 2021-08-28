package main

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Tanibox/tania-core/src/eventbus"
	"github.com/asaskevich/EventBus"
	"golang.org/x/term"

	"github.com/Tanibox/tania-core/config"
	assetsserver "github.com/Tanibox/tania-core/src/assets/server"
	assetsstorage "github.com/Tanibox/tania-core/src/assets/storage"
	growthserver "github.com/Tanibox/tania-core/src/growth/server"
	growthstorage "github.com/Tanibox/tania-core/src/growth/storage"
	locationserver "github.com/Tanibox/tania-core/src/location/server"
	tasksserver "github.com/Tanibox/tania-core/src/tasks/server"
	taskstorage "github.com/Tanibox/tania-core/src/tasks/storage"
	userserver "github.com/Tanibox/tania-core/src/user/server"
	"github.com/go-sql-driver/mysql"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mattn/go-colorable"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

func init() {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		log.SetFormatter(&log.TextFormatter{ForceColors: true})

		// We need this for Windows to get coloured
		// https://github.com/sirupsen/logrus#formatters
		log.SetOutput(colorable.NewColorableStdout())
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}

	config.InitViperConfig()
}

func main() {
	e := echo.New()

	// Initialize DB.
	log.Print("Using " + *config.Config.TaniaPersistenceEngine + " persistence engine")

	// InMemory DB will always be initialized.
	inMem := initInMemory()

	var db *sql.DB
	switch *config.Config.TaniaPersistenceEngine {
	case config.DB_SQLITE:
		db = initSqlite()
	case config.DB_MYSQL:
		db = initMysql()
	}

	// Initialize Event Bus
	bus := eventbus.NewSimpleEventBus(EventBus.New())

	// Initialize Server
	farmServer, err := assetsserver.NewFarmServer(
		db,
		inMem.farmEventStorage,
		inMem.farmReadStorage,
		inMem.areaEventStorage,
		inMem.areaReadStorage,
		inMem.reservoirEventStorage,
		inMem.reservoirReadStorage,
		inMem.materialEventStorage,
		inMem.materialReadStorage,
		inMem.cropReadStorage,
		bus,
	)
	if err != nil {
		e.Logger.Fatal(err)
	}

	taskServer, err := tasksserver.NewTaskServer(
		db,
		bus,
		inMem.cropReadStorage,
		inMem.areaReadStorage,
		inMem.materialReadStorage,
		inMem.reservoirReadStorage,
		inMem.taskEventStorage,
		inMem.taskReadStorage,
	)
	if err != nil {
		e.Logger.Fatal(err)
	}

	growthServer, err := growthserver.NewGrowthServer(
		db,
		bus,
		inMem.cropEventStorage,
		inMem.cropReadStorage,
		inMem.cropActivityStorage,
		inMem.areaReadStorage,
		inMem.materialReadStorage,
		inMem.farmReadStorage,
		inMem.taskReadStorage,
	)
	if err != nil {
		e.Logger.Fatal(err)
	}

	userServer, err := userserver.NewUserServer(db, bus)
	if err != nil {
		e.Logger.Fatal(err)
	}

	authServer, err := userserver.NewAuthServer(db, bus)
	if err != nil {
		e.Logger.Fatal(err)
	}

	locationServer, err := locationserver.NewLocationServer()
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Initialize user
	err = initUser(authServer)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Initialize Echo Middleware
	e.Use(middleware.Recover())
	e.Use(headerNoCache)
	e.Use(logrusMiddleware())
	e.Use(middleware.RequestID())

	APIMiddlewares := []echo.MiddlewareFunc{}
	if !*config.Config.DemoMode {
		APIMiddlewares = append(APIMiddlewares, tokenValidationWithConfig(db))
	}

	// HTTP routing
	API := e.Group("api")
	API.Use(middleware.CORS())

	// AuthServer is used for endpoint that doesn't need authentication checking
	authGroup := API.Group("/")
	authServer.Mount(authGroup)

	locationGroup := API.Group("/locations", APIMiddlewares...)
	locationServer.Mount(locationGroup)

	farmGroup := API.Group("/farms", APIMiddlewares...)
	farmServer.Mount(farmGroup)
	growthServer.Mount(farmGroup)

	taskGroup := API.Group("/tasks", APIMiddlewares...)
	taskServer.Mount(taskGroup)

	userGroup := API.Group("/user", APIMiddlewares...)
	userServer.Mount(userGroup)

	e.Static("/", "public")

	// Start Server
	e.Logger.Fatal(e.Start(":" + *config.Config.AppPort))
}

func initUser(authServer *userserver.AuthServer) error {
	defaultUsername := "tania"
	defaultPassword := "tania"

	_, _, err := authServer.RegisterNewUser(defaultUsername, defaultPassword, defaultPassword)
	if err != nil {
		log.Print("User ", defaultUsername, " has already created")
		return err
	}

	log.Print("User created with default username and password")

	return nil
}

// MIDDLEWARES

func headerNoCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		c.Response().Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
		c.Response().Header().Set("Expires", "0")                                         // Proxies.
		return next(c)
	}
}

type InMemory struct {
	farmEventStorage      *assetsstorage.FarmEventStorage
	farmReadStorage       *assetsstorage.FarmReadStorage
	areaEventStorage      *assetsstorage.AreaEventStorage
	areaReadStorage       *assetsstorage.AreaReadStorage
	reservoirEventStorage *assetsstorage.ReservoirEventStorage
	reservoirReadStorage  *assetsstorage.ReservoirReadStorage
	materialEventStorage  *assetsstorage.MaterialEventStorage
	materialReadStorage   *assetsstorage.MaterialReadStorage
	cropEventStorage      *growthstorage.CropEventStorage
	cropReadStorage       *growthstorage.CropReadStorage
	cropActivityStorage   *growthstorage.CropActivityStorage
	taskEventStorage      *taskstorage.TaskEventStorage
	taskReadStorage       *taskstorage.TaskReadStorage
}

func initInMemory() *InMemory {
	return &InMemory{
		farmEventStorage: assetsstorage.CreateFarmEventStorage(),
		farmReadStorage:  assetsstorage.CreateFarmReadStorage(),

		areaEventStorage: assetsstorage.CreateAreaEventStorage(),
		areaReadStorage:  assetsstorage.CreateAreaReadStorage(),

		reservoirEventStorage: assetsstorage.CreateReservoirEventStorage(),
		reservoirReadStorage:  assetsstorage.CreateReservoirReadStorage(),

		materialEventStorage: assetsstorage.CreateMaterialEventStorage(),
		materialReadStorage:  assetsstorage.CreateMaterialReadStorage(),

		cropEventStorage:    growthstorage.CreateCropEventStorage(),
		cropReadStorage:     growthstorage.CreateCropReadStorage(),
		cropActivityStorage: growthstorage.CreateCropActivityStorage(),

		taskEventStorage: taskstorage.CreateTaskEventStorage(),
		taskReadStorage:  taskstorage.CreateTaskReadStorage(),
	}
}

func initMysql() *sql.DB {
	host := *config.Config.MysqlHost
	port := *config.Config.MysqlPort
	dbname := *config.Config.MysqlDbname
	user := *config.Config.MysqlUsername
	pwd := *config.Config.MysqlPassword

	dsn := user + ":" + pwd + "@(" + host + ":" + port + ")/" + dbname + "?parseTime=true&clientFoundRows=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	log.Print("Using MySQL at ", host, ":", port, "/", dbname)

	ddl, err := ioutil.ReadFile("db/mysql/ddl.sql")
	if err != nil {
		panic(err)
	}
	sqls := string(ddl)

	// We need to split the DDL query by `;` and execute it one by one.
	// Because sql.DB.Exec() from mysql driver cannot executes multiple query at once
	// and it will give weird syntax error messages.
	splitted := strings.Split(sqls, ";")

	for _, v := range splitted {
		trimmed := strings.TrimSpace(v)

		if len(trimmed) > 0 {
			_, err = db.Exec(v)

			if err != nil {
				me, ok := err.(*mysql.MySQLError)
				if !ok {
					panic("Error executing DDL query")
				}

				// http://dev.mysql.com/doc/refman/5.7/en/error-messages-server.html
				// We will skip error duplicate key name in database (code: 1061),
				// because CREATE INDEX doesn't have IF NOT EXISTS clause,
				// otherwise we will stop the loop and print the error
				if me.Number == 1061 {

				} else {
					log.Print(err)
					return db
				}
			}
		}
	}

	log.Print("DDL file executed")

	return db
}

func initSqlite() *sql.DB {
	if _, err := os.Stat(*config.Config.SqlitePath); os.IsNotExist(err) {
		log.Print("Creating database file ", *config.Config.SqlitePath)
	}

	db, err := sql.Open("sqlite3", *config.Config.SqlitePath)
	if err != nil {
		panic(err)
	}

	log.Print("Using SQLite at ", *config.Config.SqlitePath)

	log.Print("Executing DDL file for ", *config.Config.SqlitePath)

	ddl, err := ioutil.ReadFile("db/sqlite/ddl.sql")
	if err != nil {
		panic(err)
	}
	sql := string(ddl)

	_, err = db.Exec(sql)
	if err != nil {
		panic(err)
	}

	log.Print("DDL file executed")

	return db
}

func tokenValidationWithConfig(db *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorization := c.Request().Header.Get("Authorization")

			if authorization == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"data": "Unauthorized"})
			}

			splitted := strings.Split(authorization, " ")
			if len(splitted) <= 1 {
				return c.JSON(http.StatusUnauthorized, map[string]string{"data": "Unauthorized"})
			}

			var uid interface{}
			err := db.QueryRow(`SELECT USER_UID
				FROM USER_AUTH WHERE ACCESS_TOKEN = ?`, splitted[1]).Scan(&uid)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"data": "Unauthorized"})
			}

			ubyte, ok := uid.([]byte)
			if !ok {
				return c.JSON(http.StatusInternalServerError, map[string]string{"data": "Error user UID type assertion"})
			}

			var userUID uuid.UUID
			if *config.Config.TaniaPersistenceEngine == config.DB_SQLITE {
				userUID, err = uuid.FromString(string(ubyte))
			} else if *config.Config.TaniaPersistenceEngine == config.DB_MYSQL {
				ubyte := uid.([]byte)
				userUID, err = uuid.FromBytes(ubyte)
			}

			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]error{"data": err})
			}

			c.Set("USER_UID", userUID)

			return next(c)
		}
	}
}

func logrusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			start := time.Now()

			if err := next(c); err != nil {
				c.Error(err)
			}

			res := c.Response()
			stop := time.Now()

			fields := map[string]interface{}{
				"request_id":      res.Header().Get(echo.HeaderXRequestID),
				"ip":              c.RealIP(),
				"host":            req.Host,
				"uri":             req.RequestURI,
				"method":          req.Method,
				"user_agent":      req.UserAgent(),
				"status":          res.Status,
				"roundtrip":       strconv.FormatInt(stop.Sub(start).Nanoseconds()/1000, 10),
				"roundtrip_human": stop.Sub(start).String(),
			}

			// We will add log Form Values and Query String if...
			if res.Status == http.StatusInternalServerError {
				if !strings.HasPrefix(req.Header.Get(echo.HeaderContentType), echo.MIMEMultipartForm) {
					qs := c.QueryString()

					forms, err := c.FormParams()
					if err != nil {
						c.Error(err)
					}

					fields["query_string"] = qs
					fields["form_values"] = forms
				}
			}

			log.WithFields(fields).Info("Response Info")

			return nil
		}
	}
}
