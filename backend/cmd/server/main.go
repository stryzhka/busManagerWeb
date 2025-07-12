package main

import (
	_ "backend/docs"
	"backend/pkg/controller"
	"backend/pkg/database"
	"backend/pkg/repository"
	"backend/pkg/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// @title           Bus manager API
// @version         1.0
// @description     API for bus manager backend

// @host      localhost:8080
// @BasePath  /api

// --@securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	err := initConfig()
	if err != nil {
		fmt.Println(err)
	}
	db, err := database.NewPostgresDatabase()
	if err != nil {
		panic(err)
	}
	busRepo, err := repository.NewPostgresBusRepository(db)
	if err != nil {
		panic(err)
	}
	driverRepo, err := repository.NewPostgresDriverRepository(db)
	if err != nil {
		panic(err)
	}
	busStopRepo, err := repository.NewPostgresBusStopRepository(db)
	if err != nil {
		panic(err)
	}
	routeRepo, err := repository.NewPostgresRouteRepository(db)
	if err != nil {
		panic(err)
	}
	busService := service.NewBusService(busRepo)
	driverService := service.NewDriverService(driverRepo)
	busStopService := service.NewBusStopService(busStopRepo)
	routeService := service.NewRouteService(routeRepo, driverRepo, busRepo, busStopRepo)
	busController := controller.NewBusController(*busService)
	driverController := controller.NewDriverController(*driverService)
	busStopController := controller.NewBusStopController(*busStopService)
	routeController := controller.NewRouteController(routeService)

	router := gin.Default()
	router.GET("/api/buses/id/:id", busController.GetById)
	router.GET("/api/buses/number/:number", busController.GetByNumber)
	router.GET("/api/buses/", busController.GetAll)
	router.POST("/api/buses/", busController.Add)
	router.DELETE("/api/buses/:id", busController.DeleteById)
	router.PUT("/api/buses/:id", busController.UpdateById)

	router.GET("/api/drivers/id/:id", driverController.GetById)
	router.GET("/api/drivers/series/:series", driverController.GetByPassportSeries)
	router.GET("/api/drivers/", driverController.GetAll)
	router.POST("/api/drivers/", driverController.Add)
	router.DELETE("/api/drivers/:id", driverController.DeleteById)
	router.PUT("/api/drivers/:id", driverController.UpdateById)

	router.GET("/api/stops/id/:id", busStopController.GetById)
	router.GET("/api/stops/name/:name", busStopController.GetByName)
	router.GET("/api/stops/", busStopController.GetAll)
	router.POST("/api/stops/", busStopController.Add)
	router.DELETE("/api/stops/:id", busStopController.DeleteById)
	router.PUT("/api/stops/:id", busStopController.UpdateById)

	router.GET("/api/routes/:id", routeController.GetById)
	router.GET("/api/routes/number/:number", routeController.GetByNumber)
	router.GET("/api/routes/", routeController.GetAll)
	router.POST("/api/routes/", routeController.Add)
	router.DELETE("/api/routes/:id", routeController.DeleteById)
	router.PUT("/api/routes/:id", routeController.UpdateById)
	router.POST("/api/routes/:id/drivers/:driverId", routeController.AssignDriver)
	router.POST("/api/routes/:id/stops/:busStopId", routeController.AssignBusStop)
	router.POST("/api/routes/:id/buses/:busId", routeController.AssignBus)
	router.GET("/api/routes/:id/drivers/", routeController.GetAllDriversById)
	router.GET("/api/routes/:id/stops/", routeController.GetAllBusStopsById)
	router.GET("/api/routes/:id/buses/", routeController.GetAllBusesById)
	router.DELETE("/api/routes/:id/drivers/:driverId", routeController.UnassignDriver)
	router.DELETE("/api/routes/:id/stops/:busStopId", routeController.UnassignBusStop)
	router.DELETE("/api/routes/:id/buses/:busId", routeController.UnassignBus)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(viper.GetString("server_host_port"))
}

func initConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("configs")
	return viper.ReadInConfig()

}
