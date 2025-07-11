package controller

import (
	"backend/pkg/models"
	"backend/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RouteController struct {
	rs service.IRouteService
}

func NewRouteController(rs service.IRouteService) *RouteController {
	return &RouteController{rs}
}

func (rc RouteController) GetById(c *gin.Context) {
	id := c.Param("id")
	data, err := rc.rs.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (rc RouteController) GetByNumber(c *gin.Context) {
	number := c.Param("number")
	data, err := rc.rs.GetByNumber(number)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (rc RouteController) GetAll(c *gin.Context) {
	data, err := rc.rs.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func (rc RouteController) Add(c *gin.Context) {
	var route models.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := rc.rs.Add(&route)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, route)
}

func (rc RouteController) DeleteById(c *gin.Context) {
	id := c.Param("id")
	err := rc.rs.DeleteById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": id})
}

func (rc RouteController) UpdateById(c *gin.Context) {
	var route models.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := rc.rs.UpdateById(&route)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, route)
}

func (rc RouteController) AssignDriver(c *gin.Context) {
	routeId := c.Param("id")
	driverId := c.Param("driverId")
	err := rc.rs.AssignDriver(routeId, driverId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": routeId})
}

func (rc RouteController) AssignBusStop(c *gin.Context) {
	routeId := c.Param("id")
	busStopId := c.Param("busStopId")
	err := rc.rs.AssignBusStop(routeId, busStopId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": routeId})
}

func (rc RouteController) AssignBus(c *gin.Context) {
	routeId := c.Param("id")
	busId := c.Param("busId")
	err := rc.rs.AssignBus(routeId, busId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": routeId})
}

func (rc RouteController) UnassignDriver(c *gin.Context) {
	routeId := c.Param("id")
	driverId := c.Param("driverId")
	err := rc.rs.UnassignDriver(routeId, driverId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": routeId})
}

func (rc RouteController) UnassignBusStop(c *gin.Context) {
	routeId := c.Param("id")
	busStopId := c.Param("busStopId")
	err := rc.rs.UnassignBusStop(routeId, busStopId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": routeId})
}

func (rc RouteController) UnassignBus(c *gin.Context) {
	routeId := c.Param("id")
	busId := c.Param("busId")
	err := rc.rs.UnassignBus(routeId, busId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": routeId})
}

func (rc RouteController) GetAllDriversById(c *gin.Context) {
	id := c.Param("id")
	data, err := rc.rs.GetAllDriversById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, data)
}

func (rc RouteController) GetAllBusesById(c *gin.Context) {
	id := c.Param("id")
	data, err := rc.rs.GetAllBusesById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, data)
}

func (rc RouteController) GetAllBusStopsById(c *gin.Context) {
	id := c.Param("id")
	data, err := rc.rs.GetAllBusStopsById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, data)
}
