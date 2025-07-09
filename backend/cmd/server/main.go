package main

import (
	"backend/pkg/database"
	"backend/pkg/repository"
	"fmt"
)

func main() {
	//router := gin.Default()
	//router.GET("/hello", func(c *gin.Context) {
	//	c.IndentedJSON(http.StatusOK, "hello pidor")
	//})
	//router.Run("localhost:8080")
	db, err := database.NewPostgresDatabase()
	if err != nil {
		fmt.Println(err)
	}
	p, err := repository.NewPostgresBusRepository(db)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p.GetAll())
}
