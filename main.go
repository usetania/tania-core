package main

import (
	"database/sql"
	"io/ioutil"
	"os"

	"github.com/Tanibox/tania-server/config"
	"github.com/Tanibox/tania-server/routing"
	"github.com/Tanibox/tania-server/src/assets/server"
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

	farmServer, err := server.NewFarmServer()
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
