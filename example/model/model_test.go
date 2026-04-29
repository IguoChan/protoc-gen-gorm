package model

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/IguoChan/go-pkg/mysqlx"
	"gorm.io/datatypes"
)

var (
	c           *mysqlx.Client
	migrateOnce sync.Once
	migrateErr  error
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

func ensureUserTable(t *testing.T) {
	t.Helper()
	migrateOnce.Do(func() {
		migrateErr = AutoMigrate(c.DB, &UserModel{})
	})
	if migrateErr != nil {
		t.Fatal(migrateErr)
	}
}

func testUser(prefix string) *UserModel {
	now := time.Now().Truncate(time.Second)
	token := fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
	return &UserModel{
		Name:    token,
		Email:   token + "@example.com",
		Address: "addr_" + token,
		Phone:   "10086",
		Score:   1.23,
		Balance: 9.87,
		Date:    datatypes.Date(now),
		Extra:   []byte("abc"),
		Author: &User_AuthorModel{
			Name:  "author_" + token,
			Email: "author_" + token + "@example.com",
		},
		Status:  Status_PENDING.String(),
		Status1: Status_UNKNOWN.String(),
		Status2: Status_PENDING,
		T1:      now,
		T2:      now,
		T3:      now,
		D1:      datatypes.Date(now),
		Ids:     []int32{1, 2, 3},
		M1:      map[string]string{"k": "v"},
	}
}

func seedUser(t *testing.T, prefix string) *UserModel {
	t.Helper()
	ensureUserTable(t)

	u := testUser(prefix)
	if err := Create(context.Background(), c.DB, u); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = c.DB.Unscoped().Where("id = ?", u.Id).Delete(&UserModel{}).Error
	})
	return u
}

func userDao(t *testing.T) *UserDao {
	t.Helper()
	ensureUserTable(t)
	return NewUserDao(c.DB)
}
