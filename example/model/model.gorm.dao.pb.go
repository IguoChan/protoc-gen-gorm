// Code generated by protoc-gen-gorm. DO NOT EDIT.
// versions:
// - protoc-gen-gorm v0.0.0
// - protoc           v4.25.1
// source: model.proto

package model

import (
	context "context"
	gorm "gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) Create(ctx context.Context, model *UserModel, opts ...Option) (*UserModel, error) {
	return model, Create(ctx, dao.db, model, opts...)
}

func (dao *UserDao) Save(ctx context.Context, model *UserModel, opts ...Option) (*UserModel, error) {
	return model, Save(ctx, dao.db, model, opts...)
}

func (dao *UserDao) First(ctx context.Context, opts ...Option) (*UserModel, error) {
	model := &UserModel{}
	return model, First(ctx, dao.db, model, opts...)
}

func (dao *UserDao) FirstForUpdate(ctx context.Context, opts ...Option) (*UserModel, error) {
	model := &UserModel{}
	return model, FirstForUpdate(ctx, dao.db, model, opts...)
}

func (dao *UserDao) Find(ctx context.Context, opts ...Option) ([]*UserModel, error) {
	var models []*UserModel
	return models, Find(ctx, dao.db, &models, opts...)
}

func (dao *UserDao) FindForUpdate(ctx context.Context, opts ...Option) ([]*UserModel, error) {
	var models []*UserModel
	return models, FindForUpdate(ctx, dao.db, &models, opts...)
}

func (dao *UserDao) PartialUpdate(ctx context.Context, model *UserModel, fields map[string]interface{}, opts ...Option) (*UserModel, error) {
	return model, PartialUpdate(ctx, dao.db, model, fields, opts...)
}

func (dao *UserDao) Count(ctx context.Context, opts ...Option) (int64, error) {
	return Count(ctx, dao.db, &UserModel{}, opts...)
}

func (dao *UserDao) Delete(ctx context.Context, opts ...Option) error {
	return Delete(ctx, dao.db, &UserModel{}, opts...)
}

func (dao *UserDao) DeleteUnscoped(ctx context.Context, opts ...Option) error {
	return UnscopedDelete(ctx, dao.db, &UserModel{}, opts...)
}

type StudentDao struct {
	db *gorm.DB
}

func NewStudentDao(db *gorm.DB) *StudentDao {
	return &StudentDao{
		db: db,
	}
}

func (dao *StudentDao) Create(ctx context.Context, model *StudentModel, opts ...Option) (*StudentModel, error) {
	return model, Create(ctx, dao.db, model, opts...)
}

func (dao *StudentDao) Save(ctx context.Context, model *StudentModel, opts ...Option) (*StudentModel, error) {
	return model, Save(ctx, dao.db, model, opts...)
}

func (dao *StudentDao) First(ctx context.Context, opts ...Option) (*StudentModel, error) {
	model := &StudentModel{}
	return model, First(ctx, dao.db, model, opts...)
}

func (dao *StudentDao) FirstForUpdate(ctx context.Context, opts ...Option) (*StudentModel, error) {
	model := &StudentModel{}
	return model, FirstForUpdate(ctx, dao.db, model, opts...)
}

func (dao *StudentDao) Find(ctx context.Context, opts ...Option) ([]*StudentModel, error) {
	var models []*StudentModel
	return models, Find(ctx, dao.db, &models, opts...)
}

func (dao *StudentDao) FindForUpdate(ctx context.Context, opts ...Option) ([]*StudentModel, error) {
	var models []*StudentModel
	return models, FindForUpdate(ctx, dao.db, &models, opts...)
}

func (dao *StudentDao) PartialUpdate(ctx context.Context, model *StudentModel, fields map[string]interface{}, opts ...Option) (*StudentModel, error) {
	return model, PartialUpdate(ctx, dao.db, model, fields, opts...)
}

func (dao *StudentDao) Count(ctx context.Context, opts ...Option) (int64, error) {
	return Count(ctx, dao.db, &StudentModel{}, opts...)
}

func (dao *StudentDao) Delete(ctx context.Context, opts ...Option) error {
	return Delete(ctx, dao.db, &StudentModel{}, opts...)
}

func (dao *StudentDao) DeleteUnscoped(ctx context.Context, opts ...Option) error {
	return UnscopedDelete(ctx, dao.db, &StudentModel{}, opts...)
}
