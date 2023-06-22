package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

type Person struct {
	ID        uint   `json:id`
	FirstName string `json:first_name`
	LastName  string `json:last_name`
	City      string `json:city`
	Country   string `json:country`
}

func main() {

	db, err = gorm.Open("sqlite3", "./mydatabase.db")

	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	db.AutoMigrate(&Person{})

	r := gin.Default()
	r.GET("/people/", GetPeople)
	r.GET("/people/:id", GetPerson)

	r.GET("/peopleByName/:firstname", GetPersonByName)
	r.GET("/peopleByLastname/:lastname", GetPersonByLastname)

	r.POST("/people", CreatePerson)
	r.PUT("/people/:id", UpdatePerson)
	r.DELETE("/people/:id", DeletePerson)

	r.GET("/peopleNumber/", GetPersonNumber)

	r.GET("/peopleByName2/:firstname/", GetPersonByName2)

	r.Run(":8080")
}

func DeletePerson(c *gin.Context) {
	id := c.Params.ByName("id")
	var person Person
	d := db.Where("id = ?", id).Delete(&person)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}

func UpdatePerson(c *gin.Context) {

	var person Person
	id := c.Params.ByName("id")

	if err := db.Where("id = ?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&person)

	db.Save(&person)
	c.JSON(200, person)

}

func CreatePerson(c *gin.Context) {

	var person Person

	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	db.Create(&person)
	c.JSON(200, person)

}

func GetPerson(c *gin.Context) {
	id := c.Params.ByName("id")
	var person Person
	if err := db.Where("id = ?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, person)
	}
}

func GetPeople(c *gin.Context) {
	var people []Person
	if err := db.Find(&people).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, people)
	}
}

func GetPersonByName(c *gin.Context) {
	firstname := c.Params.ByName("firstname")
	var person []Person
	if err := db.Where("first_name like ?", "%"+firstname+"%").Find(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, person)
	}
}

func GetPersonByLastname(c *gin.Context) {
	lastname := c.Params.ByName("lastname")
	var person []Person
	if err := db.Where("last_name like ?", lastname).Find(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, person)
	}
}

func GetPersonNumber(c *gin.Context) {
	var people []Person
	if err := db.Find(&people).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, len(people))
	}
}

func GetPersonByName2(c *gin.Context) {
	firstname := c.Params.ByName("firstname")
	var people []Person
	db.Raw("select * from people where first_name = ?", firstname).Scan(&people)

	c.JSON(200, people)

}
