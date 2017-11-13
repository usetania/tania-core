package main

import (
	"github.com/labstack/echo"

	"./routing"
)

func main() {
	e := echo.New()

	// Group for Farms HTTP routing
	farms := e.Group("/farms")
	farms.GET("/", routing.FarmsIndex)
	farms.POST("/", routing.FarmsCreate)
	farms.GET("/:uid", routing.FarmsShow)
	farms.PUT("/:uid", routing.FarmsUpdate)
	farms.DELETE("/:uid", routing.FarmsDestroy)

	// Group for Areas HTTP routing
	areas := e.Group("/areas")
	areas.GET("/", routing.AreasIndex)
	areas.POST("/", routing.AreasCreate)
	areas.POST("/:uid/harvest", routing.AreasHarvest)
	areas.POST("/:uid/dispose", routing.AreasDispose)
	areas.DELETE("/:uid", routing.AreasDestroy)

	// Group for Reservoirs HTTP routing
	reservoirs := e.Group("/reservoirs")
	reservoirs.GET("/", routing.ReservoirsIndex)

	e.Logger.Fatal(e.Start(":1323"))
}
