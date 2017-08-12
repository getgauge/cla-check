package data

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	// Pull the postgres drivers
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User data from Github after signing the CLA
type User struct {
	Name        string
	Email       string
	NickName    string
	UserID      string
	Description string
}

var database *gorm.DB

// Init the database
func Init() *gorm.DB {
	db, err := gorm.Open(dialect(), connection())

	if err != nil {
		fmt.Println(err)

	}
	db.AutoMigrate(&User{})
	database = db

	return db
}

// Save user to a postgres db
func Save(user User) {
	var count int
	database.Model(&User{}).Where("user_id = ?", user.UserID).Count(&count)
	if count == 0 {
		database.Save(user)
	}
}

// Signed check if a github user is registered with the DB
func Signed(nickName string) bool {
	result := User{}
	database.Where("nick_name = ?", nickName).First(&result)
	return result.NickName == nickName
}

//TODO: Move environment variables to a configuration package
func dialect() string {
	return environmentVariable("DB_DIALECT")
}

func connection() string {
	return environmentVariable("DB_CONNECTION")
}

func environmentVariable(variable string) string {
	value := os.Getenv(variable)
	if value == "" {
		log.Fatal(fmt.Sprintf("$%s must be set", variable))
	}
	return value
}