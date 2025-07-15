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

// @Summary      Get bus stop
// @Description  Get bus stop by ID
// @Tags         stops
// @Produce      json
// @Param        id   path      string  true  "Bus stop ID"
// @Success      200  {object}  models.BusStop
// @Failure      400  {object}  string
// @Router       /stops/id/{id}/ [get]
func (bsc BusStopController) GetById(c *gin.Context) {
	id := c.Param("id")
	data, err := bsc.bss.GetById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary      Get bus stop
// @Description  Get bus stop by name
// @Tags         stops
// @Produce      json
// @Param        name   path      string  true  "Bus stop name"
// @Success      200  {object}  models.BusStop
// @Failure      400  {object}  string
// @Router       /stops/name/{name}/ [get]
func (bsc BusStopController) GetByName(c *gin.Context) {
	name := c.Param("name")
	data, err := bsc.bss.GetByName(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary      Get bus stop list
// @Description  Get bus stop list
// @Tags         stops
// @Produce      json
// @Success      200  {array}  models.BusStop
// @Failure      400  {object}  string
// @Router       /stops/ [get]
func (bsc BusStopController) GetAll(c *gin.Context) {
	data, err := bsc.bss.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary      Add bus stop
// @Description  Add bus stop
// @Tags         stops
// @Produce      json
// @Param bus body models.BusStop required "bus stop model"
// @Success      200  {object}  models.BusStop
// @Failure      400  {object}  string
// @Router       /stops/ [post]
func (bsc BusStopController) Add(c *gin.Context) {
	var busStop models.BusStop
	if err := c.ShouldBindJSON(&busStop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := bsc.bss.Add(&busStop)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, busStop)
}

// @Summary      Delete bus stop
// @Description  Delete bus stop by ID
// @Tags         stops
// @Produce      json
// @Param        id   path      string  true  "Bus stop ID"
// @Success      200  {object}  string
// @Failure      400  {object}  string
// @Router       /stops/{id}/ [delete]
func (bsc BusStopController) DeleteById(c *gin.Context) {
	id := c.Param("id")
	err := bsc.bss.DeleteById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": id})
}

// @Summary      Update bus stop
// @Description  Update bus stop by ID
// @Tags         stops
// @Produce      json
// @Param        id   path      string  true  "Bus stop ID"
// @Param bus body models.BusStop required "bus stop model"
// @Success      200  {object}  models.BusStop
// @Failure      400  {object}  string
// @Router       /stops/{id}/ [put]
func (bsc BusStopController) UpdateById(c *gin.Context) {
	var busStop models.BusStop
	if err := c.ShouldBindJSON(&busStop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := bsc.bss.UpdateById(&busStop)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, busStop)
}
