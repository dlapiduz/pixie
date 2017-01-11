package common

import (
	"log"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func LoadDB() (*gorm.DB, error) {
	var err error
	var db *gorm.DB

	if strings.Contains(os.Args[0], "test") {
		// We are doing testing!
		log.Println("Connecting to TEST DB")
		db, err = gorm.Open("sqlite3", ":memory:")
	} else {
		log.Println("Connecting to DB")
		db, err = gorm.Open("postgres", os.Getenv("PSQL_CONN"))
		// DB.LogMode(true)
		db.DB().SetMaxOpenConns(10)
	}

	if err != nil {
		log.Println("Error connecting to DB")
		log.Println(err)
		return nil, err
	}

	log.Println("Connected to DB...")
	db.AutoMigrate(&Action{})
	return db, nil
}
