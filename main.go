package main

import (
	"database/sql"
	"io/ioutil"
	"os"

	"github.com/Tanibox/tania-server/farm/server"
	"github.com/Tanibox/tania-server/routing"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	e := echo.New()

	farmServer, err := server.NewFarmServer()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

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
