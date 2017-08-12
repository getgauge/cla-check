package data

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// Pull the postgres drivers
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/zabil/cla-check/configuration"
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
	db, err := gorm.Open(configuration.Dialect(), configuration.Connection())

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
