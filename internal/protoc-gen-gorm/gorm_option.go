package protoc_gen_gorm

import "google.golang.org/protobuf/compiler/protogen"

type GormOptionGenerator struct {
	gen  *protogen.Plugin
	file *protogen.File
	g    *protogen.GeneratedFile
}

func NewGormOptionGenerator(gen *protogen.Plugin, file *protogen.File, opts ...Option) *GormOptionGenerator {
	for _, opt := range opts {
		opt(defaultOptions)
	}

	return &GormOptionGenerator{
		gen:  gen,
		file: file,
	}
}

func (gg *GormOptionGenerator) GenerateFile() *protogen.GeneratedFile {
	if len(gg.file.Messages) == 0 || !existExtension(gg.file.Messages) {
		return nil
	}

	filename := gg.file.GeneratedFilenamePrefix + ".gorm.option.pb.go"
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

	gg.generateOption()

	return gg.g
}

func (gg *GormOptionGenerator) generateOption() {
	gg.g.P("type Option func(*", gormPackage.Ident("DB"), ")  *", gormPackage.Ident("DB"))
	gg.g.P()

	gg.autoMigrate()
	gg.generateNormalOption()
	gg.genFinisherFunction()
}

func (gg *GormOptionGenerator) generateNormalOption() {
	gg.dbOption()
	gg.selectOption()
	gg.tableNameOption()
	gg.limitOption()
	gg.offsetOption()
	gg.orderOption()
}

func (gg *GormOptionGenerator) autoMigrate() {
	gg.g.P("func AutoMigrate(db *", gormPackage.Ident("DB"), ", dst ...interface{}) ", " error {")
	gg.g.P("if db == nil {")
	gg.g.P("panic(\"db is nil\")")
	gg.g.P("}")
	gg.g.P("return db.AutoMigrate(dst...)")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) dbOption() {
	gg.g.P("func DB(db *", gormPackage.Ident("DB"), ") ", " Option {")
	gg.g.P("if db == nil {")
	gg.g.P("panic(\"db is nil\")")
	gg.g.P("}")
	gg.g.P("return func(*", gormPackage.Ident("DB"), ") *", gormPackage.Ident("DB"), " {")
	gg.g.P("return db")
	gg.g.P("}")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) selectOption() {
	gg.g.P("func Select(fields ...string) ", "Option {")
	gg.g.P("return func(db *", gormPackage.Ident("DB"), ") *", gormPackage.Ident("DB"), " {")
	gg.g.P("return db.Select(fields)")
	gg.g.P("}")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) tableNameOption() {
	gg.g.P("func TableName(name string) ", "Option {")
	gg.g.P("return func(db *", gormPackage.Ident("DB"), ") *", gormPackage.Ident("DB"), " {")
	gg.g.P("return db.Table(name)")
	gg.g.P("}")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) limitOption() {
	gg.g.P("func Limit(limit int) ", "Option {")
	gg.g.P("return func(db *", gormPackage.Ident("DB"), ") *", gormPackage.Ident("DB"), " {")
	gg.g.P("return db.Limit(limit)")
	gg.g.P("}")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) offsetOption() {
	gg.g.P("func Offset(offset int) ", "Option {")
	gg.g.P("return func(db *", gormPackage.Ident("DB"), ") *", gormPackage.Ident("DB"), " {")
	gg.g.P("return db.Offset(offset)")
	gg.g.P("}")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) orderOption() {
	gg.g.P("func Order(value interface{}) ", "Option {")
	gg.g.P("return func(db *", gormPackage.Ident("DB"), ") *", gormPackage.Ident("DB"), " {")
	gg.g.P("return db.Order(value)")
	gg.g.P("}")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) genFinisherFunction() {
	gg.genCreateFunction()
	gg.genSaveFunction()
	gg.genFirstFunction()
	gg.genFirstForUpdateFunction()
	gg.genFindFunction()
	gg.genFindForUpdateFunction()
	gg.genPartialUpdateFunction()
	gg.genCountFunction()
	gg.genDeleteFunction()
	gg.genUnscopedDeleteFunction()

	gg.genIsDbClosure()
}

func (gg *GormOptionGenerator) genCreateFunction() {
	gg.g.P("func Create(ctx ", contextPackage.Ident("Context"), ", db *", gormPackage.Ident("DB"), ", model interface{}, opts ...Option) error {")
	gg.genSwapDBFirst()
	gg.genApplyOpts()
	gg.g.P("return db.WithContext(ctx).Create(model).Error")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) genSaveFunction() {
	gg.g.P("func Save(ctx ", contextPackage.Ident("Context"), ", db *", gormPackage.Ident("DB"), ", model interface{}, opts ...Option) error {")
	gg.genSwapDBFirst()
	gg.genApplyOpts()
	gg.g.P("return db.WithContext(ctx).Save(model).Error")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) genFirstFunction() {
	gg.g.P("func First(ctx ", contextPackage.Ident("Context"), ", db *", gormPackage.Ident("DB"), ", model interface{}, opts ...Option) error {")
	gg.genSwapDBFirst()
	gg.genApplyOpts()
	gg.g.P("return db.WithContext(ctx).First(model).Error")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) genFirstForUpdateFunction() {
	gg.g.P("func FirstForUpdate(ctx ", contextPackage.Ident("Context"), ", db *", gormPackage.Ident("DB"), ", model interface{}, opts ...Option) error {")
	gg.genSwapDBFirst()
	gg.genApplyOpts()
	gg.g.P("return db.WithContext(ctx).Clauses(", clausePackage.Ident("Locking"), "{Strength: \"UPDATE\"}).First(model, \"\").Error")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) genFindFunction() {
	gg.g.P("func Find(ctx ", contextPackage.Ident("Context"), ", db *", gormPackage.Ident("DB"), ", model interface{}, opts ...Option) error {")
	gg.genSwapDBFirst()
	gg.genApplyOpts()
	gg.g.P("return db.WithContext(ctx).Find(model).Error")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) genFindForUpdateFunction() {
	gg.g.P("func FindForUpdate(ctx ", contextPackage.Ident("Context"), ", db *", gormPackage.Ident("DB"), ", model interface{}, opts ...Option) error {")
	gg.genSwapDBFirst()
	gg.genApplyOpts()
	gg.g.P("return db.WithContext(ctx).Clauses(", clausePackage.Ident("Locking"), "{Strength: \"UPDATE\"}).Find(model, \"\").Error")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) genPartialUpdateFunction() {
	gg.g.P("func PartialUpdate(ctx ", contextPackage.Ident("Context"), ", db *", gormPackage.Ident("DB"), ", model interface{}, fields map[string]interface{}, opts ...Option) error {")
	gg.genSwapDBFirst()
	gg.genApplyOpts()
	gg.g.P("return db.WithContext(ctx).Model(model).Updates(fields).Error")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) genCountFunction() {
	gg.g.P("func Count(ctx ", contextPackage.Ident("Context"), ", db *", gormPackage.Ident("DB"), ", model interface{}, opts ...Option) (int64, error) {")
	gg.genSwapDBFirst()
	gg.genApplyOpts()
	gg.g.P("var count int64")
	gg.g.P("return count, db.WithContext(ctx).Model(model).Count(&count).Error")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) genDeleteFunction() {
	gg.g.P("func Delete(ctx ", contextPackage.Ident("Context"), ", db *", gormPackage.Ident("DB"), ", model interface{}, opts ...Option) error {")
	gg.genSwapDBFirst()
	gg.genApplyOpts()
	gg.g.P("return db.WithContext(ctx).Delete(model).Error")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) genUnscopedDeleteFunction() {
	gg.g.P("func UnscopedDelete(ctx ", contextPackage.Ident("Context"), ", db *", gormPackage.Ident("DB"), ", model interface{}, opts ...Option) error {")
	gg.genSwapDBFirst()
	gg.genApplyOpts()
	gg.g.P("return db.WithContext(ctx).Unscoped().Delete(model).Error")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) genSwapDBFirst() {
	gg.g.P("for i, opt := range opts {")
	gg.g.P("if IsDBClosure(opt) {")
	gg.g.P("opts[i], opts[0] = opts[0], opts[i]")
	gg.g.P("break")
	gg.g.P("}")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) genApplyOpts() {
	gg.g.P("for _, opt := range opts {")
	gg.g.P("db = opt(db)")
	gg.g.P("}")
	gg.g.P()
}

func (gg *GormOptionGenerator) genIsDbClosure() {
	gg.g.P("func IsDBClosure(opt Option) bool {")
	gg.g.P("fn := ", runtimePackage.Ident("FuncForPC"), "(", reflectPackage.Ident("ValueOf"), "(opt).Pointer()).Name()")
	gg.g.P("seps := []rune{'/', '.'}")
	gg.g.P("fields := ", stringsPackage.Ident("FieldsFunc"), "(fn, func(sep rune) bool {")
	gg.g.P("for _, s := range seps {")
	gg.g.P("if sep == s {")
	gg.g.P("return true")
	gg.g.P("}")
	gg.g.P("}")
	gg.g.P("return false")
	gg.g.P("})")

	gg.g.P("if size := len(fields); size > 2 {")
	gg.g.P("// return fields[size-1] == \"func1\" && fields[size-2] == \"DB\" && fields[size-3] == \"", gg.file.GoPackageName, "\"")
	gg.g.P("return fields[size-1] == \"func1\" && fields[size-2] == \"DB\" // imprecise judgment")
	gg.g.P("}")
	gg.g.P("return false")
	gg.g.P("}")
	gg.g.P()
}
