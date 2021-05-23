package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MYSQL_USERNAME = os.Getenv("MYSQL_USERNAME")
var MYSQL_PASSWORD = os.Getenv("MYSQL_PASSWORD")
var MYSQL_DB = os.Getenv("MYSQL_DB")
var MYSQL_HOST = os.Getenv("MYSQL_HOST")
var MYSQL_PORT = os.Getenv("MYSQL_PORT")

var Db *gorm.DB = nil

func getDb() *gorm.DB {
	if Db == nil {
		Db = connect()
	}
	return Db
}

func connect() *gorm.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?authenticationDatabase=admin", MYSQL_USERNAME, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_PORT, MYSQL_DB)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting database: ", err)
		return nil
	}

		log.Println("Db is ready")
	return db
}
