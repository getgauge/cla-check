package data

import (
	"fmt"
	"strings"

	"github.com/getgauge/cla-check/configuration"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User data from Github after signing the CLA
type User struct {
	Time        string `csv:"time"`
	Name        string `csv:"name"`
	Email       string `csv:"email"`
	NickName    string `csv:"nick_name"`
	UserID      string `csv:"-"`
	Description string `csv:"Something"`
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
	database.Model(&User{}).Where("nick_name = ?", user.NickName).Count(&count)
	if count == 0 {
		database.Save(user)
	}
}

// Signed check if a github user is registered with the DB
func Signed(nickName string) bool {
	result := User{}
	database.Where("UPPER(nick_name) = ?", strings.ToUpper(nickName)).First(&result)
	return strings.EqualFold(result.NickName, nickName)
}
