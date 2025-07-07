package main

import "github.com/gin-gonic/gin"
import "net/http"

func main() {
	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, "hello pidor")
	})
	router.Run("localhost:8080")
}
