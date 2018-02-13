package main

import (
	"database/sql"
	"io/ioutil"
	"os"
	"time"

	"github.com/Tanibox/tania-server/config"
	"github.com/Tanibox/tania-server/routing"
	assetsdomain "github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/server"
	assetsserver "github.com/Tanibox/tania-server/src/assets/server"
	assetsstorage "github.com/Tanibox/tania-server/src/assets/storage"
	growthserver "github.com/Tanibox/tania-server/src/growth/server"
	growthstorage "github.com/Tanibox/tania-server/src/growth/storage"
	taskserver "github.com/Tanibox/tania-server/src/tasks/server"
	"github.com/asaskevich/EventBus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/paked/configure"
	uuid "github.com/satori/go.uuid"
)

func init() {
	initConfig()
}

func main() {
	e := echo.New()

	// Initialize all In-memory storage, so it can be used in all server
	farmStorage := assetsstorage.CreateFarmStorage()
	farmEventStorage := assetsstorage.CreateFarmEventStorage()
	farmReadStorage := assetsstorage.CreateFarmReadStorage()

	areaStorage := assetsstorage.CreateAreaStorage()
	reservoirStorage := assetsstorage.CreateReservoirStorage()
	materialStorage := assetsstorage.CreateMaterialStorage()

	cropStorage := growthstorage.CreateCropStorage()
	cropEventStorage := growthstorage.CreateCropEventStorage()
	cropReadStorage := growthstorage.CreateCropReadStorage()
	cropActivityStorage := growthstorage.CreateCropActivityStorage()

	// Initialize Event Bus
	bus := EventBus.New()

	farmServer, err := server.NewFarmServer(
		farmStorage,
		farmEventStorage,
		farmReadStorage,
		areaStorage,
		reservoirStorage,
		materialStorage,
		cropStorage,
		bus,
	)
	if err != nil {
		e.Logger.Fatal(err)
	}

	taskServer, err := taskserver.NewTaskServer(
		cropReadStorage,
		areaStorage,
		materialStorage,
		reservoirStorage)
	if err != nil {
		e.Logger.Fatal(err)
	}

	growthServer, err := growthserver.NewGrowthServer(
		bus,
		cropStorage,
		cropEventStorage,
		cropReadStorage,
		cropActivityStorage,
		areaStorage,
		materialStorage,
		farmStorage,
	)
	if err != nil {
		e.Logger.Fatal(err)
	}

	if *config.Config.DemoMode {
		initDataDemo(
			farmServer, growthServer,
			farmStorage, areaStorage, reservoirStorage, materialStorage, cropStorage,
		)
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
		UploadPathCrop: conf.String("UploadPathCrop", "/home/tania/uploads", "Upload path for the Crop photo"),
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

func initDataDemo(
	farmServer *assetsserver.FarmServer,
	growthServer *growthserver.GrowthServer,
	farmStorage *assetsstorage.FarmStorage,
	areaStorage *assetsstorage.AreaStorage,
	reservoirStorage *assetsstorage.ReservoirStorage,
	materialStorage *assetsstorage.MaterialStorage,
	cropStorage *growthstorage.CropStorage,
) {
	log.Info("==== DEMO DATA SEEDED ====")

	farmUID, _ := uuid.NewV4()
	farm1 := assetsdomain.Farm{
		UID:         farmUID,
		Name:        "MyFarm",
		Type:        "organic",
		Latitude:    "10.00",
		Longitude:   "11.00",
		CountryCode: "ID",
		CityCode:    "JK",
		IsActive:    true,
	}

	farmStorage.FarmMap[farmUID] = farm1

	uid, _ := uuid.NewV4()

	noteUID, _ := uuid.NewV4()
	reservoirNotes := make(map[uuid.UUID]assetsdomain.ReservoirNote, 0)
	reservoirNotes[noteUID] = assetsdomain.ReservoirNote{
		UID:         noteUID,
		Content:     "Don't forget to close the bucket after using",
		CreatedDate: time.Now(),
	}

	reservoir1 := assetsdomain.Reservoir{
		UID:         uid,
		Name:        "MyBucketReservoir",
		PH:          8,
		EC:          12.5,
		Temperature: 29,
		WaterSource: assetsdomain.Bucket{Capacity: 100, Volume: 10},
		Farm:        farm1,
		Notes:       reservoirNotes,
		CreatedDate: time.Now(),
	}

	farm1.AddReservoir(&reservoir1)
	farmStorage.FarmMap[farmUID] = farm1
	reservoirStorage.ReservoirMap[uid] = reservoir1

	uid, _ = uuid.NewV4()
	reservoir2 := assetsdomain.Reservoir{
		UID:         uid,
		Name:        "MyTapReservoir",
		PH:          8,
		EC:          12.5,
		Temperature: 29,
		WaterSource: assetsdomain.Tap{},
		Farm:        farm1,
		Notes:       make(map[uuid.UUID]assetsdomain.ReservoirNote),
		CreatedDate: time.Now(),
	}

	farm1.AddReservoir(&reservoir2)
	farmStorage.FarmMap[farmUID] = farm1
	reservoirStorage.ReservoirMap[uid] = reservoir2

	uid, _ = uuid.NewV4()

	noteUID, _ = uuid.NewV4()
	areaNotes := make(map[uuid.UUID]assetsdomain.AreaNote, 0)
	areaNotes[noteUID] = assetsdomain.AreaNote{
		UID:         noteUID,
		Content:     "This area should only be used for seeding.",
		CreatedDate: time.Now(),
	}

	areaSeeding := assetsdomain.Area{
		UID:       uid,
		Name:      "MySeedingArea",
		Size:      assetsdomain.AreaSize{Value: 10, Unit: assetsdomain.AreaUnit{Symbol: assetsdomain.SquareMeter}},
		Type:      assetsdomain.GetAreaType(assetsdomain.AreaTypeSeeding),
		Location:  assetsdomain.GetAreaLocation(assetsdomain.AreaLocationIndoor),
		Photo:     assetsdomain.AreaPhoto{},
		Notes:     areaNotes,
		Reservoir: reservoir2,
		Farm:      farm1,
	}

	farm1.AddArea(&areaSeeding)
	farmStorage.FarmMap[farmUID] = farm1
	areaStorage.AreaMap[uid] = areaSeeding

	uid, _ = uuid.NewV4()
	areaGrowing := assetsdomain.Area{
		UID:       uid,
		Name:      "MyGrowingArea",
		Size:      assetsdomain.AreaSize{Value: 50, Unit: assetsdomain.AreaUnit{Symbol: assetsdomain.Hectare}},
		Type:      assetsdomain.GetAreaType(assetsdomain.AreaTypeGrowing),
		Location:  assetsdomain.GetAreaLocation(assetsdomain.AreaLocationOutdoor),
		Photo:     assetsdomain.AreaPhoto{},
		Notes:     make(map[uuid.UUID]assetsdomain.AreaNote),
		Reservoir: reservoir1,
		Farm:      farm1,
	}

	farm1.AddArea(&areaGrowing)
	farmStorage.FarmMap[farmUID] = farm1
	areaStorage.AreaMap[uid] = areaGrowing

	// uid, _ = uuid.NewV4()
	// inventory1 := assetsdomain.Material{
	// 	UID:       uid,
	// 	PlantType: assetsdomain.Vegetable{},
	// 	Variety:   "Bayam Lu Hsieh",
	// }

	// inventoryMaterialStorage.InventoryMaterialMap[uid] = inventory1

	// uid, _ = uuid.NewV4()
	// inventory2 := assetsdomain.InventoryMaterial{
	// 	UID:       uid,
	// 	PlantType: assetsdomain.Vegetable{},
	// 	Variety:   "Tomat Super One",
	// }

	// inventoryMaterialStorage.InventoryMaterialMap[uid] = inventory2

	// uid, _ = uuid.NewV4()
	// inventory3 := assetsdomain.InventoryMaterial{
	// 	UID:       uid,
	// 	PlantType: assetsdomain.Fruit{},
	// 	Variety:   "Apple Rome Beauty",
	// }

	// inventoryMaterialStorage.InventoryMaterialMap[uid] = inventory3

	// uid, _ = uuid.NewV4()
	// inventory4 := assetsdomain.InventoryMaterial{
	// 	UID:       uid,
	// 	PlantType: assetsdomain.Fruit{},
	// 	Variety:   "Orange Sweet Mandarin",
	// }

	// inventoryMaterialStorage.InventoryMaterialMap[uid] = inventory4

	// /******************************
	// CROP
	// ******************************/
	// now := strings.ToLower(time.Now().Format("2Jan"))

	// uid, _ = uuid.NewV4()

	// noteUID, _ = uuid.NewV4()
	// cropNotes := make(map[uuid.UUID]growthdomain.CropNote, 0)
	// cropNotes[noteUID] = growthdomain.CropNote{
	// 	UID:         noteUID,
	// 	Content:     "This crop must be intensely watched because its expensive",
	// 	CreatedDate: time.Now(),
	// }

	// crop1 := growthdomain.Crop{
	// 	UID:     uid,
	// 	BatchID: fmt.Sprintf("%s%s", "bay-lu-hsi-", now),
	// 	Status:  growthdomain.CropStatus{Code: growthdomain.CropActive},
	// 	Type:    growthdomain.CropType{Code: growthdomain.CropTypeSeeding},
	// 	Container: growthdomain.CropContainer{
	// 		Quantity: 10,
	// 		Type:     growthdomain.Tray{Cell: 15},
	// 	},
	// 	InventoryUID: inventory1.UID,
	// 	FarmUID:      farmUID,
	// 	CreatedDate:  time.Now(),
	// 	InitialArea: growthdomain.InitialArea{
	// 		AreaUID:         areaSeeding.UID,
	// 		InitialQuantity: 10,
	// 		CurrentQuantity: 10,
	// 	},
	// 	Notes: cropNotes,
	// }

	// cropStorage.CropMap[uid] = crop1

	// uid, _ = uuid.NewV4()

	// noteUID, _ = uuid.NewV4()
	// cropNotes = make(map[uuid.UUID]growthdomain.CropNote, 0)
	// cropNotes[noteUID] = growthdomain.CropNote{
	// 	UID:         noteUID,
	// 	Content:     "This crop must be intensely watched because its expensive",
	// 	CreatedDate: time.Now(),
	// }

	// crop2 := growthdomain.Crop{
	// 	UID:     uid,
	// 	BatchID: fmt.Sprintf("%s%s", "tom-sup-one-", now),
	// 	Status:  growthdomain.CropStatus{Code: growthdomain.CropActive},
	// 	Type:    growthdomain.CropType{Code: growthdomain.CropTypeSeeding},
	// 	Container: growthdomain.CropContainer{
	// 		Quantity: 50,
	// 		Type:     growthdomain.Pot{},
	// 	},
	// 	InventoryUID: inventory2.UID,
	// 	FarmUID:      farmUID,
	// 	CreatedDate:  time.Now(),
	// 	InitialArea: growthdomain.InitialArea{
	// 		AreaUID:         areaSeeding.UID,
	// 		InitialQuantity: 50,
	// 		CurrentQuantity: 50,
	// 	},
	// 	Notes: cropNotes,
	// }

	// cropStorage.CropMap[uid] = crop2
}
