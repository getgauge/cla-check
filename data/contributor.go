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
	db, err := gorm.Open(dialect(), os.Getenv("DB_CONNECTION"))

	if err != nil {
		fmt.Println(err)

	}
	db.AutoMigrate(&User{})
	database = db

	return db
}

// Save user to a postgres db
func Save(user User) {
	database.Save(user)
}

// IsCommitter check if a github user is registered with the DB
func IsCommitter(nickname string) bool {
	result := User{}
	database.Where("nick_name = ?", "IronMan").First(&result)
	return result.NickName == nickname
}

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
