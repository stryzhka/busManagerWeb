package controller

import (
	_ "backend/docs"
	"backend/pkg/models"
	"backend/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BusController struct {
	bs service.IBusService
}

func NewBusController(bs service.BusService) *BusController {
	return &BusController{bs}
}

// @Summary      Get bus
// @Description  Get bus by ID
// @Tags         buses
// @Produce      json
// @Param        id   path      string  true  "Bus ID"
// @Success      200  {object}  models.Bus
// @Failure      404  {object}  string
// @Router       /buses/id/{id} [get]
func (bc BusController) GetById(c *gin.Context) {
	id := c.Param("id")
	data, err := bc.bs.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary      Get bus
// @Description  Get bus by number
// @Tags         buses
// @Produce      json
// @Param        number   path      string  true  "Bus register number"
// @Success      200  {object}  models.Bus
// @Failure      404  {object}  string
// @Router       /buses/number/{number} [get]
func (bc BusController) GetByNumber(c *gin.Context) {
	number := c.Param("number")
	data, err := bc.bs.GetByNumber(number)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary      Get bus list
// @Description  Get bus list
// @Tags         buses
// @Produce      json
// @Success      200  {array}  models.Bus
// @Failure      404  {object}  string
// @Router       /buses/ [get]
func (bc BusController) GetAll(c *gin.Context) {
	data := bc.bs.GetAll()
	c.JSON(http.StatusOK, data)
}

// @Summary      Add bus
// @Description  Add bus
// @Tags         buses
// @Produce      json
// @Param bus body models.Bus required "bus model"
// @Success      200  {object}  models.Bus
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /buses/ [post]
func (bc BusController) Add(c *gin.Context) {
	var bus models.Bus
	if err := c.ShouldBindJSON(&bus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := bc.bs.Add(&bus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bus)
}

// @Summary      Delete bus
// @Description  Delete bus by ID
// @Tags         buses
// @Produce      json
// @Param        id   path      string  true  "Bus ID"
// @Success      200  {object}  string
// @Failure      500  {object}  string
// @Router       /buses/{id} [delete]
func (bc BusController) DeleteById(c *gin.Context) {
	id := c.Param("id")
	err := bc.bs.DeleteById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": id})
}

// @Summary      Update bus
// @Description  Update bus by ID
// @Tags         buses
// @Produce      json
// @Param        id   path      string  true  "Bus ID"
// @Param bus body models.Bus required "bus model"
// @Success      200  {object}  models.Bus
// @Failure      500  {object}  string
// @Router       /buses/{id} [put]
func (bc BusController) UpdateById(c *gin.Context) {
	var bus models.Bus
	if err := c.ShouldBindJSON(&bus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err := bc.bs.UpdateById(&bus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bus)
}
