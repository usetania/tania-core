package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/Tanibox/tania-server/config"
	"github.com/Tanibox/tania-server/routing"
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/server"
	assetsstorage "github.com/Tanibox/tania-server/src/assets/storage"
	growthdomain "github.com/Tanibox/tania-server/src/growth/domain"
	growthserver "github.com/Tanibox/tania-server/src/growth/server"
	growthstorage "github.com/Tanibox/tania-server/src/growth/storage"
	taskserver "github.com/Tanibox/tania-server/src/tasks/server"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/paked/configure"
	deadlock "github.com/sasha-s/go-deadlock"
	uuid "github.com/satori/go.uuid"
)

func init() {
	initConfig()
}

func main() {
	e := echo.New()

	// Initialize all In-memory storage, so it can be used in all server
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("DEADLOCK!")
	}
	farmStorage := assetsstorage.FarmStorage{FarmMap: make(map[uuid.UUID]domain.Farm), Lock: &rwMutex}
	areaStorage := assetsstorage.AreaStorage{AreaMap: make(map[uuid.UUID]domain.Area), Lock: &rwMutex}
	reservoirStorage := assetsstorage.ReservoirStorage{ReservoirMap: make(map[uuid.UUID]domain.Reservoir), Lock: &rwMutex}
	inventoryMaterialStorage := assetsstorage.InventoryMaterialStorage{InventoryMaterialMap: make(map[uuid.UUID]domain.InventoryMaterial), Lock: &rwMutex}
	cropStorage := growthstorage.CropStorage{CropMap: make(map[uuid.UUID]growthdomain.Crop), Lock: &rwMutex}

	farmServer, err := server.NewFarmServer(
		&farmStorage,
		&areaStorage,
		&reservoirStorage,
		&inventoryMaterialStorage,
	)
	if err != nil {
		e.Logger.Fatal(err)
	}

	taskServer, err := taskserver.NewTaskServer()
	if err != nil {
		e.Logger.Fatal(err)
	}

	growthServer, err := growthserver.NewGrowthServer(
		&cropStorage,
		&areaStorage,
		&inventoryMaterialStorage,
		&farmStorage,
	)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(headerNoCache)

	// Bootstraping Database
	pwd, _ := os.Getwd()
	db := initDB(pwd + "/resources/storage.db")
	migrate(db)

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

func initDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)

	if err != nil {
		panic(err)
	}

	if db == nil {
		panic("db nill")
	}

	return db
}

func migrate(db *sql.DB) {
	pwd, _ := os.Getwd()
	filerc, err := ioutil.ReadFile(pwd + "/resources/structure.sql")
	if err != nil {
		panic(err)
	}
	sql := string(filerc)

	_, err = db.Exec(sql)

	if err != nil {
		panic(err)
	}
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
		DemoMode:       conf.Bool("DemoMode", true, "Switch for the demo mode"),
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
