package main

import (
	_ "backend/docs"
	"backend/pkg"
	"backend/pkg/controller"
	"backend/pkg/database"
	"backend/pkg/repository"
	"backend/pkg/service"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"os"
	"strings"
)

// @title           Bus manager API
// @version         1.0
// @description     API for bus manager backend

// @host      localhost:8080
// @BasePath  /api

// --@securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	err, connStr, serverConn := initConfig()
	if err != nil {
		panic(err)
	}
	db, err := database.NewPostgresDatabase(connStr)
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
	userRepo, err := repository.NewPostgresUserRepository(db)
	if err != nil {
		panic(err)
	}
	busService := service.NewBusService(busRepo)
	driverService := service.NewDriverService(driverRepo)
	busStopService := service.NewBusStopService(busStopRepo)
	routeService := service.NewRouteService(routeRepo, driverRepo, busRepo, busStopRepo)
	userService := service.NewUserService(userRepo)
	busController := controller.NewBusController(*busService)
	driverController := controller.NewDriverController(*driverService)
	busStopController := controller.NewBusStopController(*busStopService)
	routeController := controller.NewRouteController(routeService)
	userController := controller.NewUserController(*userService)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5173", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,

		MaxAge: 12 * 60 * 60, // 12 часов
	}))
	// Группировка API роутов под /api
	api := router.Group("/api")
	{
		// Группа для автобусов
		buses := api.Group("/buses")
		buses.Use(func(c *gin.Context) {
			pkg.UserIdentity(c, *userService)
		})
		{
			buses.GET("/id/:id", busController.GetById)
			buses.GET("/number/:number", busController.GetByNumber)
			buses.GET("/", busController.GetAll)
			buses.POST("/", busController.Add)
			buses.DELETE("/:id", busController.DeleteById)
			buses.PUT("/:id", busController.UpdateById)
		}

		// Группа для водителей
		drivers := api.Group("/drivers")
		drivers.Use(func(c *gin.Context) {
			pkg.UserIdentity(c, *userService)
		})
		{
			drivers.GET("/id/:id", driverController.GetById)
			drivers.GET("/series/:series", driverController.GetByPassportSeries)
			drivers.GET("/", driverController.GetAll)
			drivers.POST("/", driverController.Add)
			drivers.DELETE("/:id", driverController.DeleteById)
			drivers.PUT("/:id", driverController.UpdateById)
		}

		// Группа для остановок
		stops := api.Group("/stops")
		stops.Use(func(c *gin.Context) {
			pkg.UserIdentity(c, *userService)
		})
		{
			stops.GET("/id/:id", busStopController.GetById)
			stops.GET("/name/:name", busStopController.GetByName)
			stops.GET("/", busStopController.GetAll)
			stops.POST("/", busStopController.Add)
			stops.DELETE("/:id", busStopController.DeleteById)
			stops.PUT("/:id", busStopController.UpdateById)
		}

		// Группа для маршрутов
		routes := api.Group("/routes")
		routes.Use(func(c *gin.Context) {
			pkg.UserIdentity(c, *userService)
		})
		{
			routes.GET("/:id", routeController.GetById)
			routes.GET("/number/:number", routeController.GetByNumber)
			routes.GET("/", routeController.GetAll)
			routes.POST("/", routeController.Add)
			routes.DELETE("/:id", routeController.UnassignBus)
			routes.PUT("/:id", routeController.UpdateById)
			routes.POST("/:id/drivers/:driverId", routeController.AssignDriver)
			routes.POST("/:id/stops/:busStopId", routeController.AssignBusStop)
			routes.POST("/:id/buses/:busId", routeController.AssignBus)
			routes.GET("/:id/drivers", routeController.GetAllDriversById)
			routes.GET("/:id/stops", routeController.GetAllBusStopsById)
			routes.GET("/:id/buses", routeController.GetAllBusesById)
			routes.DELETE("/:id/drivers/:driverId", routeController.UnassignDriver)
			routes.DELETE("/:id/stops/:busStopId", routeController.UnassignBusStop)
			routes.DELETE("/:id/buses/:busId", routeController.UnassignBus)
		}

		// Группа для пользователей
		users := api.Group("/auth")
		{
			users.POST("/sign-in", userController.Signin)
			users.POST("/sign-up", userController.Signup)
		}
	}

	// Swagger документация
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(serverConn)
}

func initConfig() (error, string, string) {
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.SetConfigType("toml")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("configs/") // path to look for the config file in
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			return err, "", ""
		}
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	serverConn := strings.TrimSpace(os.Getenv("SERVER"))
	fmt.Println(serverConn)
	fmt.Println(connStr)
	return nil, connStr, serverConn
}
