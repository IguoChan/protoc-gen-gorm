// Code generated by protoc-gen-gorm. DO NOT EDIT.
// versions:
// - protoc-gen-gorm v0.0.0
// - protoc           v4.25.1
// source: model.proto

package model

import (
	context "context"
	gorm "gorm.io/gorm"
	clause "gorm.io/gorm/clause"
	reflect "reflect"
	runtime "runtime"
	strings "strings"
)

type Option func(*gorm.DB) *gorm.DB

func AutoMigrate(db *gorm.DB, dst ...interface{}) error {
	if db == nil {
		panic("db is nil")
	}
	return db.AutoMigrate(dst...)
}

func DB(db *gorm.DB) Option {
	if db == nil {
		panic("db is nil")
	}
	return func(*gorm.DB) *gorm.DB {
		return db
	}
}

func Select(fields ...string) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(fields)
	}
}

func TableName(name string) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table(name)
	}
}

func Limit(limit int) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

func Offset(offset int) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	}
}

func Order(value interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(value)
	}
}

func Create(ctx context.Context, db *gorm.DB, model interface{}, opts ...Option) error {
	for i, opt := range opts {
		if IsDBClosure(opt) {
			opts[i], opts[0] = opts[0], opts[i]
			break
		}
	}

	for _, opt := range opts {
		db = opt(db)
	}

	return db.WithContext(ctx).Create(model).Error
}

func Save(ctx context.Context, db *gorm.DB, model interface{}, opts ...Option) error {
	for i, opt := range opts {
		if IsDBClosure(opt) {
			opts[i], opts[0] = opts[0], opts[i]
			break
		}
	}

	for _, opt := range opts {
		db = opt(db)
	}

	return db.WithContext(ctx).Save(model).Error
}

func First(ctx context.Context, db *gorm.DB, model interface{}, opts ...Option) error {
	for i, opt := range opts {
		if IsDBClosure(opt) {
			opts[i], opts[0] = opts[0], opts[i]
			break
		}
	}

	for _, opt := range opts {
		db = opt(db)
	}

	return db.WithContext(ctx).First(model).Error
}

func FirstForUpdate(ctx context.Context, db *gorm.DB, model interface{}, opts ...Option) error {
	for i, opt := range opts {
		if IsDBClosure(opt) {
			opts[i], opts[0] = opts[0], opts[i]
			break
		}
	}

	for _, opt := range opts {
		db = opt(db)
	}

	return db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(model, "").Error
}

func Find(ctx context.Context, db *gorm.DB, model interface{}, opts ...Option) error {
	for i, opt := range opts {
		if IsDBClosure(opt) {
			opts[i], opts[0] = opts[0], opts[i]
			break
		}
	}

	for _, opt := range opts {
		db = opt(db)
	}

	return db.WithContext(ctx).Find(model).Error
}

func FindForUpdate(ctx context.Context, db *gorm.DB, model interface{}, opts ...Option) error {
	for i, opt := range opts {
		if IsDBClosure(opt) {
			opts[i], opts[0] = opts[0], opts[i]
			break
		}
	}

	for _, opt := range opts {
		db = opt(db)
	}

	return db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Find(model, "").Error
}

func PartialUpdate(ctx context.Context, db *gorm.DB, model interface{}, fields map[string]interface{}, opts ...Option) error {
	for i, opt := range opts {
		if IsDBClosure(opt) {
			opts[i], opts[0] = opts[0], opts[i]
			break
		}
	}

	for _, opt := range opts {
		db = opt(db)
	}

	return db.WithContext(ctx).Model(model).Updates(fields).Error
}

func Count(ctx context.Context, db *gorm.DB, model interface{}, opts ...Option) (int64, error) {
	for i, opt := range opts {
		if IsDBClosure(opt) {
			opts[i], opts[0] = opts[0], opts[i]
			break
		}
	}

	for _, opt := range opts {
		db = opt(db)
	}

	var count int64
	return count, db.WithContext(ctx).Model(model).Count(&count).Error
}

func Delete(ctx context.Context, db *gorm.DB, model interface{}, opts ...Option) error {
	for i, opt := range opts {
		if IsDBClosure(opt) {
			opts[i], opts[0] = opts[0], opts[i]
			break
		}
	}

	for _, opt := range opts {
		db = opt(db)
	}

	return db.WithContext(ctx).Delete(model).Error
}

func UnscopedDelete(ctx context.Context, db *gorm.DB, model interface{}, opts ...Option) error {
	for i, opt := range opts {
		if IsDBClosure(opt) {
			opts[i], opts[0] = opts[0], opts[i]
			break
		}
	}

	for _, opt := range opts {
		db = opt(db)
	}

	return db.WithContext(ctx).Unscoped().Delete(model).Error
}

func IsDBClosure(opt Option) bool {
	fn := runtime.FuncForPC(reflect.ValueOf(opt).Pointer()).Name()
	seps := []rune{'/', '.'}
	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})
	if size := len(fields); size > 2 {
		// return fields[size-1] == "func1" && fields[size-2] == "DB" && fields[size-3] == "model"
		return fields[size-1] == "func1" && fields[size-2] == "DB" // imprecise judgment
	}
	return false
}
