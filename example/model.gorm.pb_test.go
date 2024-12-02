package model

import (
	"testing"
	"time"

	"github.com/IguoChan/go-pkg/mysqlx"
	"gorm.io/datatypes"
)

var (
	c *mysqlx.Client
)

func init() {
	var err error
	c, err = mysqlx.NewClient(&mysqlx.Config{
		Addr:         "localhost:3306",
		Username:     "root",
		Password:     "cyg711024",
		DBName:       "ailurus",
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	})
	if err != nil {
		panic(err)
	}
}

func TestAutoMigrateUser(t *testing.T) {
	err := c.AutoMigrate(&UserModel{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestInsertUser(t *testing.T) {
	u := &UserModel{
		Name: "fff",
		Date: datatypes.Date{},
	}

	// create
	err := c.Create(u).Error
	if err != nil {
		t.Fatal(err)
	}

	// update
	time.Sleep(3 * time.Second)
	u.Name = "eee"
	err = c.Save(u).Error
	if err != nil {
		t.Fatal(err)
	}

	// get
	u1 := &UserModel{}
	err = c.First(u1, "id = ?", u.Id).Error
	if err != nil {
		t.Fatal(err)
	}

	// delete
	time.Sleep(3 * time.Second)
	err = c.Delete(u).Error
	if err != nil {
		t.Fatal(err)
	}
}
