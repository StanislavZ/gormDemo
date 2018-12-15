package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//Person struct represents Person in db
type Person struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var db *gorm.DB
var err error
var dbConnectionCounter int

const dbConnectionAttempts int = 30

func main() {

	db := initDatabaseConnection()

	db.AutoMigrate(&Person{})

	defer db.Close()

	initAPI()
}

func initDatabaseConnection() *gorm.DB {
	if db == nil { // during startup - if it does not exist, create it
		db = connectToDB()
	}

	connected := db != nil

	for connected != true && dbConnectionCounter < dbConnectionAttempts { // reconnect if we lost connection
		fmt.Println("Connection to Postgres was lost. Waiting for 1s...")
		time.Sleep(1 * time.Second)
		fmt.Println("Reconnecting...")
		db = connectToDB()
		connected = db != nil
	}

	if dbConnectionCounter >= dbConnectionAttempts {
		fmt.Println("Database is not available")
		os.Exit(2)
	}

	fmt.Println("Connection to Database established")

	return db
}

func connectToDB() *gorm.DB {
	db, err := gorm.Open("postgres", "host=postgreshost port=5432 user=admin password=admin dbname=demo sslmode=disable")

	dbConnectionCounter++

	if err != nil {
		fmt.Println(err)
		fmt.Println("Could not connect to Postgres.")
		return nil
	}
	return db
}

//GetPeople returns all users
func GetPeople(c *gin.Context) {
	var people []Person
	if err := db.Find(&people).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println("Error happend during DB data retreival: ", err)
	} else {
		c.JSON(200, people)
	}
}

//GetPerson returns one user
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

//CreatePerson creates user
func CreatePerson(c *gin.Context) {
	var person Person

	c.BindJSON(&person)

	db.Create(&person)
	c.JSON(200, person)
}

func initAPI() {
	r := gin.Default()

	r.GET("/", GetPeople)
	r.GET("/people/:id", GetPerson)
	r.POST("/people", CreatePerson)

	r.Run(":8090")
}
