package controller

import (
	_ "backend/docs"
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

// @Summary      Get driver
// @Description  Get driver by ID
// @Tags         drivers
// @Security ApiKeyAuth
// @Produce      json
// @Param        id   path      string  true  "Driver ID"
// @Success      200  {object}  models.Driver
// @Failure      400  {object}  string
// @Router       /drivers/id/{id}/ [get]
func (dc DriverController) GetById(c *gin.Context) {
	id := c.Param("id")
	data, err := dc.ds.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary      Get driver
// @Description  Get driver by passport series
// @Tags         drivers
// @Security ApiKeyAuth
// @Produce      json
// @Param        series   path      string  true  "Driver passport series"
// @Success      200  {object}  models.Driver
// @Failure      400  {object}  string
// @Router       /drivers/series/{series}/ [get]
func (dc DriverController) GetByPassportSeries(c *gin.Context) {
	series := c.Param("series")
	data, err := dc.ds.GetByPassportSeries(series)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary      Get driver list
// @Description  Get driver list
// @Tags         drivers
// @Security ApiKeyAuth
// @Produce      json
// @Success      200  {array}  models.Driver
// @Failure      400  {object}  string
// @Router       /drivers/ [get]
func (dc DriverController) GetAll(c *gin.Context) {
	data := dc.ds.GetAll()
	c.JSON(http.StatusOK, data)
}

// @Summary      Add driver
// @Description  Add driver
// @Tags         drivers
// @Security ApiKeyAuth
// @Produce      json
// @Param driver body models.Driver required "driver model"
// @Success      200  {object}  models.Driver
// @Failure      400  {object}  string

// @Router       /drivers/ [post]
func (dc DriverController) Add(c *gin.Context) {
	var driver models.Driver
	if err := c.ShouldBindJSON(&driver); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := dc.ds.Add(&driver)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, driver)
}

// @Summary      Delete driver
// @Description  Delete driver by ID
// @Tags         drivers
// @Security ApiKeyAuth
// @Produce      json
// @Param        id   path      string  true  "Driver ID"
// @Success      200  {object}  string
// @Failure      400  {object}  string
// @Router       /drivers/{id}/ [delete]
func (dc DriverController) DeleteById(c *gin.Context) {
	id := c.Param("id")
	err := dc.ds.DeleteById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": id})
}

// @Summary      Update driver
// @Description  Update driver by ID
// @Tags         drivers
// @Security ApiKeyAuth
// @Produce      json
// @Param        id   path      string  true  "Driver ID"
// @Param driver body models.Driver required "driver model"
// @Success      200  {object}  models.Driver
// @Failure      400  {object}  string
// @Router       /drivers/{id}/ [put]
func (dc DriverController) UpdateById(c *gin.Context) {
	var driver models.Driver
	if err := c.ShouldBindJSON(&driver); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := dc.ds.UpdateById(&driver)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, driver)
}
