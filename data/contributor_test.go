package data

import (
	"os"
	"testing"

	// Pull the sqlite drivers
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
)

func delete() {
	database.Close()
	os.Remove("test.db")
}

func TestSave(t *testing.T) {
	Init()
	defer delete()
	Save(User{"", "Tony Stark", "iron.man@avengers.com", "IronMan", "10", "CEO and Avenger"})
	savedUser := User{}
	database.Where("nick_name = ?", "IronMan").First(&savedUser)
	assert.Equal(t, "IronMan", savedUser.NickName)
}

func TestShouldNotSaveUserWithTheSameIDTwice(t *testing.T) {
	Init()
	defer delete()
	Save(User{"", "Tony Stark", "iron.man@avengers.com", "IronMan", "10", "CEO and Avenger"})
	Save(User{"", "Tony Stark", "iron.man@avengers.com", "IronMan", "10", "CEO and Avenger"})
	var count int
	database.Model(&User{}).Where("nick_name = ?", "IronMan").Count(&count)
	assert.Equal(t, 1, count)
}

func TestIfUserSigned(t *testing.T) {
	Init()
	defer delete()
	//TODO: Randomize name
	Save(User{"", "Tony Stark", "iron.man@avengers.com", "IronMan", "10", "CEO and Avenger"})
	assert.True(t, Signed("IronMan"))
}

func TestIfUserSignedIgnoreCase(t *testing.T) {
	Init()
	defer delete()
	//TODO: Randomize name
	Save(User{"", "Tony Stark", "iron.man@avengers.com", "IronMan", "10", "CEO and Avenger"})
	assert.True(t, Signed("ironman"))
}

func TestIfUserNotSigned(t *testing.T) {
	Init()
	defer delete()
	assert.False(t, Signed("IronMan"))
}

func TestMain(m *testing.M) {
	os.Setenv("DB_DIALECT", "sqlite3")
	os.Setenv("DATABASE_URL", "test.db")
	os.Exit(m.Run())
}
