package database

import (
	//"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"rest-go-demo/entity"
)

//Connector variable used for CRUD operation's
var Connector *gorm.DB

//Connect creates MySQL connection
func Connect(connectionString string) error {
	var err error
	log.Println("gorm.Open  mysql")
	Connector, err = gorm.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	log.Println("Connection was successful!!")
	return nil
}

//Migrate create/updates database table
func Migrate(table *entity.Person) {
	Connector.AutoMigrate(&table)
	log.Println("Table migrated")
}

func Migrat2(table *entity.Account) {
	Connector.AutoMigrate(&table)
	log.Println("Table migrated")
}

func Migrat3(table *entity.WaterD) {
	Connector.AutoMigrate(&table)
	log.Println("Table migrated water_d")
}
