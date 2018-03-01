package main

import (
	"database/sql"
	"io/ioutil"

	"github.com/Tanibox/tania-server/config"
	"github.com/Tanibox/tania-server/routing"
	"github.com/Tanibox/tania-server/src/assets/server"
	assetsstorage "github.com/Tanibox/tania-server/src/assets/storage"
	growthserver "github.com/Tanibox/tania-server/src/growth/server"
	growthstorage "github.com/Tanibox/tania-server/src/growth/storage"
	taskserver "github.com/Tanibox/tania-server/src/tasks/server"
	taskstorage "github.com/Tanibox/tania-server/src/tasks/storage"
	"github.com/asaskevich/EventBus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/paked/configure"
)

func init() {
	initConfig()
}

func main() {
	e := echo.New()

	// Initialize all In-memory storage, so it can be used in all server
	farmReadStorage := assetsstorage.CreateFarmReadStorage()

	areaEventStorage := assetsstorage.CreateAreaEventStorage()
	areaReadStorage := assetsstorage.CreateAreaReadStorage()

	reservoirEventStorage := assetsstorage.CreateReservoirEventStorage()
	reservoirReadStorage := assetsstorage.CreateReservoirReadStorage()

	materialEventStorage := assetsstorage.CreateMaterialEventStorage()
	materialReadStorage := assetsstorage.CreateMaterialReadStorage()

	cropEventStorage := growthstorage.CreateCropEventStorage()
	cropReadStorage := growthstorage.CreateCropReadStorage()
	cropActivityStorage := growthstorage.CreateCropActivityStorage()

	taskEventStorage := taskstorage.CreateTaskEventStorage()
	taskReadStorage := taskstorage.CreateTaskReadStorage()

	// Initialize SQLite3
	db := initSqlite()

	// Initialize Event Bus
	bus := EventBus.New()

	farmServer, err := server.NewFarmServer(
		db,
		areaEventStorage,
		areaReadStorage,
		reservoirEventStorage,
		reservoirReadStorage,
		materialEventStorage,
		materialReadStorage,
		cropReadStorage,
		bus,
	)
	if err != nil {
		e.Logger.Fatal(err)
	}

	taskServer, err := taskserver.NewTaskServer(
		db,
		bus,
		cropReadStorage,
		areaReadStorage,
		materialReadStorage,
		reservoirReadStorage,
		taskEventStorage,
		taskReadStorage)
	if err != nil {
		e.Logger.Fatal(err)
	}

	growthServer, err := growthserver.NewGrowthServer(
		db,
		bus,
		cropEventStorage,
		cropReadStorage,
		cropActivityStorage,
		areaReadStorage,
		materialReadStorage,
		farmReadStorage,
	)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(headerNoCache)

	// HTTP routing
	API := e.Group("api")
	API.Use(middleware.CORS())

	routing.LocationsRouter(API.Group("/locations"))

	farmGroup := API.Group("/farms")
	farmServer.Mount(farmGroup)
	growthServer.Mount(farmGroup)

	taskGroup := API.Group("/tasks")
	taskServer.Mount(taskGroup)

	e.Static("/", "public")
	e.Logger.Fatal(e.Start(":8080"))
}

func initSqlite() *sql.DB {
	db, err := sql.Open("sqlite3", *config.Config.SqlitePath)
	if err != nil {
		panic(err)
	}

	ddl, err := ioutil.ReadFile("db/sqlite/ddl.sql")
	if err != nil {
		panic(err)
	}
	sql := string(ddl)

	_, err = db.Exec(sql)

	if err != nil {
		panic(err)
	}

	return db
}

/*
	Example setting and usage of configure package:

	// main.initConfig()
	configuration := config.Configuration{
		// this will be filled from environment variables
		DBPassword: conf.String("TANIA_DB_PASSWORD", "123456", "Description"),

		// this will be filled from flags (ie ./tania-server --port=9000)
		Port: conf.String("port", "3000", "Description"),

		// this will be filled from conf.json
		UploadPath: conf.String("UploadPath", "/home/tania/uploads", "Description"),
	}

	// config.Configuration struct
	type Configuration struct {
		DBPassword 		*string
		Port 			*string
		UploadPath 		*string
	}

	// Usage. config.Config can be called globally
	fmt.Println(*config.Config.DBPassword)
	fmt.Println(*config.Config.Port)
	fmt.Println(*config.Config.UploadPath)

*/
func initConfig() {
	conf := configure.New()

	configuration := config.Configuration{
		UploadPathArea: conf.String("UploadPathArea", "/home/tania/uploads", "Upload path for the Area photo"),
		UploadPathCrop: conf.String("UploadPathCrop", "/home/tania/uploads", "Upload path for the Crop photo"),
		DemoMode:       conf.Bool("DemoMode", true, "Switch for the demo mode"),
		SqlitePath:     conf.String("SqlitePath", "db/sqlite/tania.db", "Path of sqlite file db"),
	}

	// This config will read the first configuration.
	// If it doesn't find the key, then it go to the next configuration.
	conf.Use(configure.NewEnvironment())
	conf.Use(configure.NewFlag())
	conf.Use(configure.NewJSONFromFile("conf.json"))

	conf.Parse()

	config.Config = configuration
}

func headerNoCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		c.Response().Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
		c.Response().Header().Set("Expires", "0")                                         // Proxies.
		return next(c)
	}
}
