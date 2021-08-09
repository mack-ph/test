package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"

	//"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"rest-go-demo/database"
	"rest-go-demo/entity"
	"strconv"
)

//GetAllPerson get all person data
func GetAllPerson(c *gin.Context) {
	fmt.Println("CGetAllPerson")

	var persons []entity.Person
	database.Connector.Find(&persons)
	// w.Header().Set("Content-Type", "application/json")

	c.JSON(http.StatusOK, persons)
	//json.NewEncoder(w).Encode(persons)
}

//GetPersonByID returns person with specific ID
func GetPersonByID(c *gin.Context) {
	key := c.Param("name")
	fmt.Println(key)
	var person entity.Person
	database.Connector.First(&person, key)
	c.JSON(http.StatusOK, person)
}

//CreatePerson creates person
func CreatePerson(c *gin.Context) {

	requestBody, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {

		c.String(http.StatusOK, "err is %s", err)
	}
	var person entity.Person
	json.Unmarshal(requestBody, &person)
	database.Connector.Create(person)
	c.JSON(http.StatusCreated, person)

}

//UpdatePersonByID updates person with respective ID
func UpdatePersonByID(c *gin.Context) {
	key := c.Param("name")

	requestBody, _ := ioutil.ReadAll(c.Request.Body)
	var person entity.Person
	json.Unmarshal(requestBody, &person)
	person.ID, _ = strconv.Atoi(key)
	database.Connector.Save(&person)
	c.JSON(http.StatusOK, person)
}

//DeletPersonByID delete's person with specific ID
func DeletPersonByID(c *gin.Context) {
	key := c.Param("name")
	var person entity.Person
	id, _ := strconv.ParseInt(key, 10, 64)
	database.Connector.Where("id = ?", id).Delete(&person)

	c.String(http.StatusOK, "sucesses")
}

func GetWaterDByID(c *gin.Context) {
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
