package controllers

import (
	// "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	// //"github.com/gorilla/mux"
	// "io/ioutil"
	"net/http"
	"rest-go-demo/database"
	"rest-go-demo/entity"
	// "strconv"
)

func GetWaterLevelByID(c *gin.Context) {
	key := c.Param("name")
	fmt.Println("GetDFByID: ", key)
	var dataflow []entity.WaterD
	values := [][]string{}

	database.Connector.Where("mn = ? and data_time > ?", key, "2020-08-08 00:12:44").Find(&dataflow)

	for i := 0; i < len(dataflow); i++ {
		var tmp []string
		tmp = append(tmp, dataflow[i].DataTime)
		tmp = append(tmp, dataflow[i].WaterLevel)
		values = append(values, tmp)
	}

	c.JSON(http.StatusOK, values)
}

func GetWaterNode(c *gin.Context) {

	key := c.Param("name")
	fmt.Println("GetDFByID: ", key)
	var dataflow []entity.WaterD

	database.Connector.Where("  data_time > ?", "2020-08-08 00:12:44").Find(&dataflow)

	c.JSON(http.StatusOK, dataflow)

}
