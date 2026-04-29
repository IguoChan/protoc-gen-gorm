package model

import (
	"context"
	"strings"
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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

func TestOptionFinisherCreateSaveFirstFindDelete(t *testing.T) {
	ctx := context.Background()
	u := testUser("option_crud")

	if err := Create(ctx, c.DB, u); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = c.DB.Unscoped().Where("id = ?", u.Id).Delete(&UserModel{}).Error
	})

	u.Phone = "10010"
	if err := Save(ctx, c.DB, u); err != nil {
		t.Fatal(err)
	}

	first := &UserModel{}
	if err := First(ctx, c.DB, first, UserModelId(u.Id)); err != nil {
		t.Fatal(err)
	}
	if first.Phone != "10010" {
		t.Fatalf("phone = %q, want %q", first.Phone, "10010")
	}

	var found []*UserModel
	if err := Find(ctx, c.DB, &found, Where("id = ?", u.Id)); err != nil {
		t.Fatal(err)
	}
	if len(found) != 1 {
		t.Fatalf("found %d users, want 1", len(found))
	}

	if err := Delete(ctx, c.DB, &UserModel{}, UserModelId(u.Id)); err != nil {
		t.Fatal(err)
	}
}

func TestOptionAdvancedQueryFinishers(t *testing.T) {
	ctx := context.Background()
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
}
