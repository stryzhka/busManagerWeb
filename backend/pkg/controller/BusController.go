package controller

import (
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

func (bc BusController) GetById(c *gin.Context) {
	id := c.Param("id")
	data, err := bc.bs.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	//jsonData, err := json.MarshalIndent(data, "", "    ")
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	c.JSON(http.StatusOK, data)
}

//func (bc BusController) GetByNumber(number string) string {
//	if strings.TrimSpace(number) == "" {
//		return responses.NewJsonError(errors.New("Number cant be null"))
//	}
//	data, err := bc.bs.GetByNumber(number)
//	if err != nil {
//		return responses.NewJsonError(err)
//	}
//	jsonData, err := json.MarshalIndent(data, "", "    ")
//	if err != nil {
//		return responses.NewJsonError(err)
//	}
//	return string(jsonData)
//}
//
//func (bc BusController) GetAll() string {
//	data := bc.bs.GetAll()
//	jsonData, err := json.MarshalIndent(data, "", "    ")
//	if err != nil {
//		return responses.NewJsonError(err)
//	}
//	return string(jsonData)
//}
//
//func (bc BusController) Add(busData string) string {
//	byteBus := []byte(busData)
//	var bus models.Bus
//	err := json.Unmarshal(byteBus, &bus)
//	if err != nil {
//		return responses.NewJsonError(err)
//	}
//	err = bc.bs.Add(&bus)
//	if err != nil {
//		return responses.NewJsonError(err)
//	}
//	return busData
//}
//
//func (bc BusController) DeleteById(id string) string {
//	if strings.TrimSpace(id) == "" {
//		return responses.NewJsonError(errors.New("ID cant be null"))
//	}
//	err := bc.bs.DeleteById(id)
//	if err != nil {
//		return responses.NewJsonError(err)
//	}
//	return ""
//}
//
//func (bc BusController) UpdateById(busData string) string {
//	byteBus := []byte(busData)
//	var bus models.Bus
//	err := json.Unmarshal(byteBus, &bus)
//	if err != nil {
//		return responses.NewJsonError(err)
//	}
//	err = bc.bs.UpdateById(&bus)
//	if err != nil {
//		return responses.NewJsonError(err)
//	}
//	return busData
//}
