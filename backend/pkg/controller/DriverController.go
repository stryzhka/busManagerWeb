package controller

import (
	"backend/pkg/models"
	"backend/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DriverController struct {
	ds service.IDriverService
}

func NewDriverController(ds service.DriverService) *DriverController {
	return &DriverController{ds}
}

func (dc DriverController) GetById(c *gin.Context) {
	id := c.Param("id")
	data, err := dc.ds.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (dc DriverController) GetByPassportSeries(c *gin.Context) {
	series := c.Param("series")
	data, err := dc.ds.GetByPassportSeries(series)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (dc DriverController) GetAll(c *gin.Context) {
	data := dc.ds.GetAll()
	c.JSON(http.StatusOK, data)
}

func (dc DriverController) Add(c *gin.Context) {
	var driver models.Driver
	if err := c.ShouldBindJSON(&driver); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := dc.ds.Add(&driver)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, driver)
}

func (dc DriverController) DeleteById(c *gin.Context) {
	id := c.Param("id")
	err := dc.ds.DeleteById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": id})
}

func (dc DriverController) UpdateById(c *gin.Context) {
	var driver models.Driver
	if err := c.ShouldBindJSON(&driver); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := dc.ds.UpdateById(&driver)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, driver)
}
