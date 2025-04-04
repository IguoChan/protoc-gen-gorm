package protoc_gen_gorm

import (
	"google.golang.org/protobuf/compiler/protogen"
)

type GormDaoGenerator struct {
	gen  *protogen.Plugin
	file *protogen.File
	g    *protogen.GeneratedFile
}

func NewGormDaoGenerator(gen *protogen.Plugin, file *protogen.File, opts ...Option) *GormDaoGenerator {
	for _, opt := range opts {
		opt(defaultOptions)
	}

	return &GormDaoGenerator{
		gen:  gen,
		file: file,
	}
}

func (gg *GormDaoGenerator) GenerateFile() *protogen.GeneratedFile {
	if len(gg.file.Messages) == 0 || !existExtension(gg.file.Messages) {
		return nil
	}

	filename := gg.file.GeneratedFilenamePrefix + ".gorm.dao.pb.go"
	gg.g = gg.gen.NewGeneratedFile(filename, gg.file.GoImportPath)
	gg.g.P("// Code generated by protoc-gen-gorm. DO NOT EDIT.")
	gg.g.P("// versions:")
	gg.g.P("// - protoc-gen-gorm v", defaultOptions.version)
	gg.g.P("// - protoc           ", protocVersion(gg.gen))
	if gg.file.Proto.GetOptions().GetDeprecated() {
		gg.g.P("// ", gg.file.Desc.Path(), " is a deprecated gg.file.")
	} else {
		gg.g.P("// source: ", gg.file.Desc.Path())
	}
	gg.g.P()

	gg.g.P("package ", gg.file.GoPackageName)
	gg.g.P()

	for _, m := range gg.file.Messages {
		if !existExtension([]*protogen.Message{m}) {
			continue
		}
		gg.genMessageDao(m)
	}

	return gg.g
}

func (gg *GormDaoGenerator) genMessageDao(message *protogen.Message) {
	gg.genMessageDaoType(message)
	gg.genCreate(message)
	gg.genSave(message)
	gg.genFirst(message)
	gg.genFirstForUpdate(message)
	gg.genFind(message)
	gg.genFindForUpdate(message)
	gg.genPartialUpdate(message)
	gg.genCount(message)
	gg.genDelete(message)
	gg.genUnscopedDelete(message)
}

func (gg *GormDaoGenerator) genMessageDaoType(message *protogen.Message) {
	gg.g.P("type ", message.GoIdent.GoName, "Dao struct {")
	gg.g.P("    db *", gormPackage.Ident("DB"))
	gg.g.P("}")
	gg.g.P()
	gg.g.P("func New", message.GoIdent.GoName, "Dao(db *", gormPackage.Ident("DB"), ") *", message.GoIdent.GoName, "Dao {")
	gg.g.P("    return &", message.GoIdent.GoName, "Dao{")
	gg.g.P("        db: db,")
	gg.g.P("    }")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormDaoGenerator) genCreate(message *protogen.Message) {
	gg.g.P("func (dao *", message.GoIdent.GoName, "Dao) Create(ctx ", contextPackage.Ident("Context"), ", model *", message.GoIdent.GoName, structSuffix, ", opts ...Option) (*", message.GoIdent.GoName, structSuffix, ", error) {")
	gg.g.P("    return model, Create(ctx, dao.db, model, opts...)")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormDaoGenerator) genSave(message *protogen.Message) {
	gg.g.P("func (dao *", message.GoIdent.GoName, "Dao) Save(ctx ", contextPackage.Ident("Context"), ", model *", message.GoIdent.GoName, structSuffix, ", opts ...Option)  (*", message.GoIdent.GoName, structSuffix, ", error) {")
	gg.g.P("    return model, Save(ctx, dao.db, model, opts...)")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormDaoGenerator) genFirst(message *protogen.Message) {
	gg.g.P("func (dao *", message.GoIdent.GoName, "Dao) First(ctx ", contextPackage.Ident("Context"), ", opts ...Option) (*", message.GoIdent.GoName, structSuffix, ", error) {")
	gg.g.P("    model := &", message.GoIdent.GoName, structSuffix, "{}")
	gg.g.P("    return model, First(ctx, dao.db, model, opts...)")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormDaoGenerator) genFirstForUpdate(message *protogen.Message) {
	gg.g.P("func (dao *", message.GoIdent.GoName, "Dao) FirstForUpdate(ctx ", contextPackage.Ident("Context"), ", opts ...Option) (*", message.GoIdent.GoName, structSuffix, ", error) {")
	gg.g.P("    model := &", message.GoIdent.GoName, structSuffix, "{}")
	gg.g.P("    return model, FirstForUpdate(ctx, dao.db, model, opts...)")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormDaoGenerator) genFind(message *protogen.Message) {
	gg.g.P("func (dao *", message.GoIdent.GoName, "Dao) Find(ctx ", contextPackage.Ident("Context"), ", opts ...Option) ([]*", message.GoIdent.GoName, structSuffix, ", error) {")
	gg.g.P("    var models []*", message.GoIdent.GoName, structSuffix)
	gg.g.P("    return models, Find(ctx, dao.db, &models, opts...)")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormDaoGenerator) genFindForUpdate(message *protogen.Message) {
	gg.g.P("func (dao *", message.GoIdent.GoName, "Dao) FindForUpdate(ctx ", contextPackage.Ident("Context"), ", opts ...Option) ([]*", message.GoIdent.GoName, structSuffix, ", error) {")
	gg.g.P("    var models []*", message.GoIdent.GoName, structSuffix)
	gg.g.P("    return models, FindForUpdate(ctx, dao.db, &models, opts...)")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormDaoGenerator) genPartialUpdate(message *protogen.Message) {
	gg.g.P("func (dao *", message.GoIdent.GoName, "Dao) PartialUpdate(ctx ", contextPackage.Ident("Context"), ", model *", message.GoIdent.GoName, structSuffix, ", fields map[string]interface{}, opts ...Option) (*", message.GoIdent.GoName, structSuffix, ", error) {")
	gg.g.P("    return model, PartialUpdate(ctx, dao.db, model, fields, opts...)")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormDaoGenerator) genCount(message *protogen.Message) {
	gg.g.P("func (dao *", message.GoIdent.GoName, "Dao) Count(ctx ", contextPackage.Ident("Context"), ", opts ...Option) (int64, error) {")
	gg.g.P("    return Count(ctx, dao.db, &", message.GoIdent.GoName, structSuffix, "{}, opts...)")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormDaoGenerator) genDelete(message *protogen.Message) {
	gg.g.P("func (dao *", message.GoIdent.GoName, "Dao) Delete(ctx ", contextPackage.Ident("Context"), ", opts ...Option) error {")
	gg.g.P("    return Delete(ctx, dao.db, &", message.GoIdent.GoName, structSuffix, "{}, opts...)")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormDaoGenerator) genUnscopedDelete(message *protogen.Message) {
	gg.g.P("func (dao *", message.GoIdent.GoName, "Dao) DeleteUnscoped(ctx ", contextPackage.Ident("Context"), ", opts ...Option) error {")
	gg.g.P("    return UnscopedDelete(ctx, dao.db, &", message.GoIdent.GoName, structSuffix, "{}, opts...)")
	gg.g.P("}")
	gg.g.P()
}
