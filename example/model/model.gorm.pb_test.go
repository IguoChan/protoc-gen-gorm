package model

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/IguoChan/go-pkg/mysqlx"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		Status:  Status_PENDING.String(),
		Status1: Status_UNKNOWN.String(),
		Status2: Status_PENDING,
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

func TestAutoMigrateUser(t *testing.T) {
	ensureUserTable(t)
}

func TestOptionBuilders(t *testing.T) {
	ensureUserTable(t)

	options := []Option{
		DB(c.DB),
		Clauses(clause.Locking{Strength: "UPDATE"}),
		Distinct("name"),
		Select("id", "name"),
		Omit("phone"),
		MapColumns(map[string]string{"name_alias": "name"}),
		TableName((&UserModel{}).TableName()),
		Where("id > ?", 0),
		Not("name = ?", "none"),
		Or("email = ?", "none@example.com"),
		Joins("LEFT JOIN `user` option_join_user ON option_join_user.id = `user`.id"),
		InnerJoins("JOIN `user` option_inner_join_user ON option_inner_join_user.id = `user`.id"),
		Group("name"),
		Having("count(*) >= ?", 0),
		Limit(1),
		Offset(0),
		Order("id desc"),
		Scopes(func(db *gorm.DB) *gorm.DB { return db.Where("id >= ?", 0) }),
		Preload(clause.Associations),
		Attrs(map[string]interface{}{"phone": "attr_phone"}),
		Assign(map[string]interface{}{"phone": "assign_phone"}),
		Unscoped(),
		Raw("SELECT 1"),
	}

	for _, opt := range options {
		if db := opt.Apply(c.DB); db == nil {
			t.Fatal("option returned nil db")
		}
	}
	if !IsDBClosure(DB(c.DB)) {
		t.Fatal("DB option should be detected as DB closure")
	}
	if IsDBClosure(Option{apply: func(*gorm.DB) *gorm.DB { return c.DB }}) {
		t.Fatal("plain closure should not be detected as DB closure")
	}
}

func TestApplyOptionsMovesDBClosureFirst(t *testing.T) {
	baseDB := &gorm.DB{}
	optionDB := &gorm.DB{}

	var sawOptionDB bool
	opt := Option{
		apply: func(db *gorm.DB) *gorm.DB {
			sawOptionDB = db == optionDB
			return db
		},
	}

	applyOptions(baseDB, opt, DB(optionDB))
	if !sawOptionDB {
		t.Fatal("DB option should be applied before earlier non-DB options")
	}
}

func TestCreateSaveFirstFindDelete(t *testing.T) {
	ctx := context.Background()
	dao := userDao(t)
	u := testUser("crud")

	created, err := dao.Create(ctx, u)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = c.DB.Unscoped().Where("id = ?", created.Id).Delete(&UserModel{}).Error
	})

	created.Phone = "10010"
	if _, err = dao.Save(ctx, created); err != nil {
		t.Fatal(err)
	}
	first, err := dao.First(ctx, UserModelId(created.Id))
	if err != nil {
		t.Fatal(err)
	}
	if first.Phone != "10010" {
		t.Fatalf("phone = %q, want %q", first.Phone, "10010")
	}
	found, err := dao.Find(ctx, Where("id = ?", created.Id))
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 1 {
		t.Fatalf("found %d users, want 1", len(found))
	}
	if err = dao.Delete(ctx, UserModelId(created.Id)); err != nil {
		t.Fatal(err)
	}
	if _, err = dao.First(ctx, UserModelId(created.Id)); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("first after delete err = %v, want ErrRecordNotFound", err)
	}
}

func TestCreateInBatches(t *testing.T) {
	ctx := context.Background()
	dao := userDao(t)
	users := []*UserModel{testUser("batch_a"), testUser("batch_b")}

	created, err := dao.CreateInBatches(ctx, users, 2)
	if err != nil {
		t.Fatal(err)
	}
	for _, u := range created {
		id := u.Id
		t.Cleanup(func() {
			_ = c.DB.Unscoped().Where("id = ?", id).Delete(&UserModel{}).Error
		})
	}
	if len(created) != 2 || created[0].Id == 0 || created[1].Id == 0 {
		t.Fatalf("unexpected batch create result: %+v", created)
	}
}

func TestTakeLastAndLockingFinders(t *testing.T) {
	ctx := context.Background()
	dao := userDao(t)
	u := seedUser(t, "finder")

	if _, err := dao.Take(ctx, UserModelId(u.Id)); err != nil {
		t.Fatal(err)
	}
	if _, err := dao.Last(ctx, UserModelId(u.Id)); err != nil {
		t.Fatal(err)
	}
	if err := dao.Transaction(ctx, func(tx *UserDao) error {
		if _, err := tx.FirstForUpdate(ctx, UserModelId(u.Id)); err != nil {
			return err
		}
		if _, err := tx.TakeForUpdate(ctx, UserModelId(u.Id)); err != nil {
			return err
		}
		if _, err := tx.LastForUpdate(ctx, UserModelId(u.Id)); err != nil {
			return err
		}
		_, err := tx.FindForUpdate(ctx, UserModelId(u.Id))
		return err
	}); err != nil {
		t.Fatal(err)
	}
}

func TestFindInBatches(t *testing.T) {
	ctx := context.Background()
	dao := userDao(t)
	u1 := seedUser(t, "find_batch_1")
	u2 := seedUser(t, "find_batch_2")

	var batches int
	err := dao.FindInBatches(ctx, 1, func(tx *gorm.DB, batch int) error {
		batches++
		return tx.Error
	}, Where("id IN ?", []int32{u1.Id, u2.Id}), Order("id asc"))
	if err != nil {
		t.Fatal(err)
	}
	if batches == 0 {
		t.Fatal("FindInBatches callback was not called")
	}
}

func TestFirstOrInit(t *testing.T) {
	ctx := context.Background()
	dao := userDao(t)
	u := testUser("first_or_init")

	got, err := dao.FirstOrInit(ctx, u, Where("name = ?", u.Name), Attrs(map[string]interface{}{"phone": "init_phone"}))
	if err != nil {
		t.Fatal(err)
	}
	if got.Id != 0 {
		t.Fatalf("FirstOrInit inserted record id %d", got.Id)
	}
}

func TestFirstOrCreate(t *testing.T) {
	ctx := context.Background()
	dao := userDao(t)
	u := testUser("first_or_create")

	got, err := dao.FirstOrCreate(ctx, u, Where("name = ? AND address = ?", u.Name, u.Address))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = c.DB.Unscoped().Where("id = ?", got.Id).Delete(&UserModel{}).Error
	})
	if got.Id == 0 {
		t.Fatal("FirstOrCreate did not create a record")
	}
}

func TestUpdateFunctions(t *testing.T) {
	ctx := context.Background()
	dao := userDao(t)
	u := seedUser(t, "updates")

	if _, err := dao.Update(ctx, u, "phone", "20001"); err != nil {
		t.Fatal(err)
	}
	if _, err := dao.Updates(ctx, u, map[string]interface{}{"phone": "20002", "balance": 2.2}); err != nil {
		t.Fatal(err)
	}
	if _, err := dao.UpdateColumn(ctx, u, "phone", "20003"); err != nil {
		t.Fatal(err)
	}
	if _, err := dao.UpdateColumns(ctx, u, map[string]interface{}{"phone": "20004", "balance": 4.4}); err != nil {
		t.Fatal(err)
	}
	if _, err := dao.PartialUpdate(ctx, u, map[string]interface{}{"phone": "20005"}); err != nil {
		t.Fatal(err)
	}

	got, err := dao.First(ctx, UserModelId(u.Id))
	if err != nil {
		t.Fatal(err)
	}
	if got.Phone != "20005" {
		t.Fatalf("phone = %q, want %q", got.Phone, "20005")
	}
}

func TestCountPluckScanQuery(t *testing.T) {
	ctx := context.Background()
	dao := userDao(t)
	u := seedUser(t, "query")

	count, err := dao.Count(ctx, UserModelId(u.Id))
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("count = %d, want 1", count)
	}

	var names []string
	if err = dao.Pluck(ctx, "name", &names, UserModelId(u.Id)); err != nil {
		t.Fatal(err)
	}
	if len(names) != 1 || names[0] != u.Name {
		t.Fatalf("names = %#v, want %q", names, u.Name)
	}

	var scanned []struct {
		Name string
	}
	if err = dao.Scan(ctx, &scanned, Raw("SELECT name FROM `user` WHERE id = ?", u.Id)); err != nil {
		t.Fatal(err)
	}
	if len(scanned) != 1 || scanned[0].Name != u.Name {
		t.Fatalf("scanned = %#v, want %q", scanned, u.Name)
	}
}

func TestRowRowsAndExec(t *testing.T) {
	ctx := context.Background()
	dao := userDao(t)
	u := seedUser(t, "rows")

	var id int32
	if err := dao.Row(ctx, Raw("SELECT id FROM `user` WHERE id = ?", u.Id)).Scan(&id); err != nil {
		t.Fatal(err)
	}
	if id != u.Id {
		t.Fatalf("row id = %d, want %d", id, u.Id)
	}

	rows, err := dao.Rows(ctx, Raw("SELECT id FROM `user` WHERE id = ?", u.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Fatal("Rows returned no records")
	}

	if err = dao.Exec(ctx, "UPDATE `user` SET phone = ? WHERE id = ?", "30001", u.Id); err != nil {
		t.Fatal(err)
	}
	got, err := dao.First(ctx, UserModelId(u.Id))
	if err != nil {
		t.Fatal(err)
	}
	if got.Phone != "30001" {
		t.Fatalf("phone = %q, want %q", got.Phone, "30001")
	}
}

func TestDeleteModelAndUnscopedDelete(t *testing.T) {
	ctx := context.Background()
	dao := userDao(t)

	soft := seedUser(t, "delete_model")
	if err := dao.DeleteModel(ctx, soft); err != nil {
		t.Fatal(err)
	}
	if _, err := dao.First(ctx, UserModelId(soft.Id)); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("soft-deleted first err = %v, want ErrRecordNotFound", err)
	}

	hardByCondition := seedUser(t, "delete_unscoped")
	if err := dao.DeleteUnscoped(ctx, UserModelId(hardByCondition.Id)); err != nil {
		t.Fatal(err)
	}
	var count int64
	if err := c.DB.Unscoped().Model(&UserModel{}).Where("id = ?", hardByCondition.Id).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("hard delete by condition count = %d, want 0", count)
	}

	hardByModel := seedUser(t, "delete_model_unscoped")
	if err := dao.DeleteModelUnscoped(ctx, hardByModel); err != nil {
		t.Fatal(err)
	}
	if err := c.DB.Unscoped().Model(&UserModel{}).Where("id = ?", hardByModel.Id).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("hard delete by model count = %d, want 0", count)
	}
}

func TestTransaction(t *testing.T) {
	ctx := context.Background()
	dao := userDao(t)
	u := seedUser(t, "transaction")

	if err := dao.Transaction(ctx, func(tx *UserDao) error {
		_, err := tx.Update(ctx, u, "phone", "40001")
		return err
	}); err != nil {
		t.Fatal(err)
	}
	got, err := dao.First(ctx, UserModelId(u.Id))
	if err != nil {
		t.Fatal(err)
	}
	if got.Phone != "40001" {
		t.Fatalf("phone = %q, want %q", got.Phone, "40001")
	}
}

func TestAdvancedQueryOptions(t *testing.T) {
	ctx := context.Background()
	dao := userDao(t)
	u := seedUser(t, "advanced")

	var found []*UserModel
	err := Find(ctx, c.DB, &found,
		Select([]string{"`user`.id", "`user`.name"}),
		Joins("LEFT JOIN `user` joined_user ON joined_user.id = `user`.id"),
		Where("`user`.id = ?", u.Id),
		Not("`user`.name = ?", "missing"),
		Or("`user`.email = ?", u.Email),
		Order("`user`.id desc"),
		Limit(1),
		Offset(0),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 1 {
		t.Fatalf("advanced query found %d users, want 1", len(found))
	}

	count, err := Count(ctx, c.DB, &UserModel{}, Group("name"), Having("count(*) >= ?", 1), Where("id = ?", u.Id))
	if err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("group/having count returned 0")
	}

	var distinctNames []string
	if err = Pluck(ctx, c.DB, &UserModel{}, "name", &distinctNames, Distinct("name"), Where("id = ?", u.Id)); err != nil {
		t.Fatal(err)
	}
	if len(distinctNames) != 1 || !strings.HasPrefix(distinctNames[0], "advanced_") {
		t.Fatalf("distinct names = %#v", distinctNames)
	}

	if _, err = dao.Find(ctx,
		InnerJoins("JOIN `user` inner_join_user ON inner_join_user.id = `user`.id"),
		Where("`user`.id = ?", u.Id),
	); err != nil {
		t.Fatal(err)
	}
}
