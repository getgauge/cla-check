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
	Save(User{"Tony Stark", "iron.man@avengers.com", "IronMan", "10", "CEO and Avenger"})
	savedUser := User{}
	database.Where("nick_name = ?", "IronMan").First(&savedUser)
	assert.Equal(t, "IronMan", savedUser.NickName)
}

func TestIfSavedUserIsCommitter(t *testing.T) {
	Init()
	defer delete()
	Save(User{"Tony Stark", "iron.man@avengers.com", "IronMan", "10", "CEO and Avenger"})
	assert.True(t, IsCommitter("IronMan"))
}

func TestIfUnSavedUserIsNotCommitter(t *testing.T) {
	Init()
	defer delete()
	assert.False(t, IsCommitter("IronMan"))
}

func TestMain(m *testing.M) {
	os.Setenv("DB_DIALECT", "sqlite3")
	os.Setenv("DB_CONNECTION", "test.db")
	os.Exit(m.Run())
}
