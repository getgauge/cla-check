package data

import (
	"fmt"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/jinzhu/gorm"
	// Pull the postgres drivers
	"github.com/getgauge/cla-check/configuration"
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

// Seed the data
func Seed() {
	clientsFile, err := os.OpenFile("cla.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer clientsFile.Close()

	clients := []*User{}

	if err := gocsv.UnmarshalFile(clientsFile, &clients); err != nil { // Load clients from file
		panic(err)
	}
	for _, client := range clients {
		Save(*client)
	}
}
