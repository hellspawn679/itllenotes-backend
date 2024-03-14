package Database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var db *gorm.DB
var notes_db *gorm.DB

type User struct {
	gorm.Model
	Name          string `json:"name" gorm:"uniqueIndex"`
	Email         string `json:"email" gorm:"uniqueIndex"`
	Password      string `json:"password"`
	User_type     string `json:"user_type"`
	Token         string `json:"token"`
	Refresh_token string `json:"refresh_token"`
}
type User_Notes struct {
	gorm.Model
	Name  string `json:"name" gorm:"uniqueIndex"`
	Title string
	Notes string
}

func ConnecttoDB() (db *gorm.DB, notes_db *gorm.DB) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	db, err = gorm.Open("postgres", "host="+os.Getenv("DB_HOST")+" port="+os.Getenv("DB_PORT")+" user="+os.Getenv("DB_USER_NAME")+" dbname="+os.Getenv("DB_NAME")+" sslmode=disable password="+os.Getenv("DB_PASSWORD"))
	if err != nil {
		log.Fatal(err)
	}
	notes_db, err = gorm.Open("postgres", "host="+os.Getenv("DB_HOST")+" port="+os.Getenv("DB_PORT")+" user="+os.Getenv("DB_USER_NAME")+" dbname="+os.Getenv("DB_NAME")+" sslmode=disable password="+os.Getenv("DB_PASSWORD"))
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&User{})
	notes_db.AutoMigrate(&User_Notes{})
	return db, notes_db
}
