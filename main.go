package main

import (
	"database/sql"
	"io/ioutil"

	"github.com/Tanibox/tania-server/routing"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Bootstraping Database
	db := initDB("resources/storage.db")
	migrate(db)

	// HTTP routing
	API := e.Group("api")
	API.Use(middleware.CORS())
	routing.AreasRouter(API.Group("/farms"))
	routing.FarmsRouter(API.Group("/areas"))
	routing.ReservoirsRouter(API.Group("/reservoirs"))
	routing.LocationsRouter(API.Group("/locations"))

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
	filerc, err := ioutil.ReadFile("resources/structure.sql")
	if err != nil {
		panic(err)
	}
	sql := string(filerc)

	_, err = db.Exec(sql)

	if err != nil {
		panic(err)
	}
}
