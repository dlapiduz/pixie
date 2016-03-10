package common

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func LoadDB() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", os.Getenv("PSQL_CONN"))

	if err != nil {
		log.Println("Error connecting to DB")
		log.Println(err)
		return nil, err
	}

	log.Println("Connected to DB...")
	db.AutoMigrate(&Action{})
	return db, nil
}
