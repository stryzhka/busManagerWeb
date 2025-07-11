package controller

import (
	"backend/pkg/models"
	"backend/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BusStopController struct {
	bss service.IBusStopService
}

func NewBusStopController(bss service.IBusStopService) *BusStopController {
	return &BusStopController{bss}
}

func (bsc BusStopController) GetById(c *gin.Context) {
	id := c.Param("id")
	data, err := bsc.bss.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (bsc BusStopController) GetByName(c *gin.Context) {
	name := c.Param("name")
	data, err := bsc.bss.GetByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (bsc BusStopController) GetAll(c *gin.Context) {
	data, err := bsc.bss.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func (bsc BusStopController) Add(c *gin.Context) {
	var busStop models.BusStop
	if err := c.ShouldBindJSON(&busStop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := bsc.bss.Add(&busStop)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, busStop)
}

func (bsc BusStopController) DeleteById(c *gin.Context) {
	id := c.Param("id")
	err := bsc.bss.DeleteById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": id})
}

func (bsc BusStopController) UpdateById(c *gin.Context) {
	var busStop models.BusStop
	if err := c.ShouldBindJSON(&busStop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := bsc.bss.UpdateById(&busStop)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, busStop)
}
