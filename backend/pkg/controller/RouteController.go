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

// @Summary      Get route
// @Description  Get route by ID
// @Tags         routes
// @Produce      json
// @Param        id   path      string  true  "Route ID"
// @Success      200  {object}  models.Route
// @Failure      404  {object}  string
// @Router       /routes/id/{id} [get]
func (rc RouteController) GetById(c *gin.Context) {
	id := c.Param("id")
	data, err := rc.rs.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary      Get route
// @Description  Get route by number
// @Tags         routes
// @Produce      json
// @Param        number   path      string  true  "Route number"
// @Success      200  {object}  models.Route
// @Failure      404  {object}  string
// @Router       /routes/id/{id} [get]
func (rc RouteController) GetByNumber(c *gin.Context) {
	number := c.Param("number")
	data, err := rc.rs.GetByNumber(number)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary      Get routes list
// @Description  Get routes list
// @Tags         routes
// @Produce      json
// @Success      200  {array}  models.Route
// @Failure      500  {object}  string
// @Router       /routes/ [get]
func (rc RouteController) GetAll(c *gin.Context) {
	data, err := rc.rs.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary      Add route
// @Description  Add route
// @Tags         routes
// @Produce      json
// @Param route body models.Route required "route model"
// @Success      200  {object}  models.Route
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /routes/ [post]
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

// @Summary      Delete route
// @Description  Delete route by ID
// @Tags         routes
// @Produce      json
// @Param        id   path      string  true  "Route ID"
// @Success      200  {object}  string
// @Failure      500  {object}  string
// @Router       /routes/{id} [delete]
func (rc RouteController) DeleteById(c *gin.Context) {
	id := c.Param("id")
	err := rc.rs.DeleteById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": id})
}

// @Summary      Update route
// @Description  Update route by ID
// @Tags         routes
// @Produce      json
// @Param        id   path      string  true  "Route ID"
// @Param route body models.Route required "route model"
// @Success      200  {object}  models.Route
// @Failure      500  {object}  string
// @Router       /routes/{id} [put]
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

// @Summary      Assign driver to route
// @Description  Assign driver to route
// @Tags         routes
// @Produce      json
// @Param        id   path      string  true  "Route ID"
// @Param        driverId   path      string  true  "Driver ID"
// @Success      200  {object}  string
// @Failure      500  {object}  string
// @Router       /routes/{id}/drivers/{driverId} [post]
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

// @Summary      Assign bus stop to route
// @Description  Assign bus stop to route
// @Tags         routes
// @Produce      json
// @Param        id   path      string  true  "Route ID"
// @Param        busStopId   path      string  true  "Bus stop ID"
// @Success      200  {object}  string
// @Failure      500  {object}  string
// @Router       /routes/{id}/stops/{busStopId} [post]
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

// @Summary      Assign bus to route
// @Description  Assign bus to route
// @Tags         routes
// @Produce      json
// @Param        id   path      string  true  "Route ID"
// @Param        busId   path      string  true  "Bus ID"
// @Success      200  {object}  string
// @Failure      500  {object}  string
// @Router       /routes/{id}/buses/{busId} [post]
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

// @Summary      Unassign driver from route
// @Description  Unassign driver from route
// @Tags         routes
// @Produce      json
// @Param        id   path      string  true  "Route ID"
// @Param        driverId   path      string  true  "Driver ID"
// @Success      200  {object}  string
// @Failure      500  {object}  string
// @Router       /routes/{id}/drivers/{driverId} [delete]
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

// @Summary      Unassign bus stop from route
// @Description  Unassign bus stop from route
// @Tags         routes
// @Produce      json
// @Param        id   path      string  true  "Route ID"
// @Param        busStopId   path      string  true  "Bus stop ID"
// @Success      200  {object}  string
// @Failure      500  {object}  string
// @Router       /routes/{id}/stops/{busStopId} [delete]
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

// @Summary      Unassign bus from route
// @Description  Unassign bus from route
// @Tags         routes
// @Produce      json
// @Param        id   path      string  true  "Route ID"
// @Param        busId   path      string  true  "Bus ID"
// @Success      200  {object}  string
// @Failure      500  {object}  string
// @Router       /routes/{id}/buses/{busId} [delete]
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

// @Summary      Get all drivers on route
// @Description  Get all drivers on route by route ID
// @Tags         routes
// @Produce      json
// @Param        id   path      string  true  "Route ID"
// @Success      200  {object}  models.Driver
// @Failure      500  {object}  string
// @Router       /routes/{id}/drivers/ [get]
func (rc RouteController) GetAllDriversById(c *gin.Context) {
	id := c.Param("id")
	data, err := rc.rs.GetAllDriversById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, data)
}

// @Summary      Get all buses on route
// @Description  Get all buses on route by route ID
// @Tags         routes
// @Produce      json
// @Param        id   path      string  true  "Route ID"
// @Success      200  {object}  models.Bus
// @Failure      500  {object}  string
// @Router       /routes/{id}/buses/ [get]
func (rc RouteController) GetAllBusesById(c *gin.Context) {
	id := c.Param("id")
	data, err := rc.rs.GetAllBusesById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, data)
}

// @Summary      Get all bus stops on route
// @Description  Get all bus stops on route by route ID
// @Tags         routes
// @Produce      json
// @Param        id   path      string  true  "Route ID"
// @Success      200  {object}  models.BusStop
// @Failure      500  {object}  string
// @Router       /routes/{id}/stops/ [get]
func (rc RouteController) GetAllBusStopsById(c *gin.Context) {
	id := c.Param("id")
	data, err := rc.rs.GetAllBusStopsById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, data)
}
