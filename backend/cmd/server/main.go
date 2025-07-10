package main

import (
	"backend/pkg/controller"
	"backend/pkg/database"
	"backend/pkg/repository"
	"backend/pkg/service"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.NewPostgresDatabase()
	if err != nil {
		panic(err)
	}
	busRepo, err := repository.NewPostgresBusRepository(db)
	if err != nil {
		panic(err)
	}
	busService := service.NewBusService(busRepo)
	busController := controller.NewBusController(*busService)

	router := gin.Default()
	router.GET("/api/buses/:id", busController.GetById)
	router.Run("localhost:8080")

}
