# protoc-gen-gorm

[中文文档](README.zh-CN.md)

`protoc-gen-gorm` is a Go `protoc` plugin that generates GORM-friendly model code from protobuf messages. It reads custom protobuf options from `option/gorm.proto` and emits Go structs, GORM tags, table names, protobuf/model conversion helpers, optional query helpers, and optional DAO wrappers.

## Features

- Generates `MessageModel` structs from protobuf messages marked with `(gorm.set_table) = true`.
- Emits GORM tags for columns, primary keys, indexes, unique indexes, defaults, `not null`, soft deletes, and SQL column types.
- Supports custom table names with `(gorm.table_name)`.
- Converts protobuf messages to generated model structs and back.
- Maps protobuf timestamps, integer/string time fields, enums, repeated fields, and maps to GORM-compatible Go types.
- Optionally generates reusable GORM helper functions with `with_gorm_option=true`.
- Optionally generates per-message DAO types with `with_gorm_dao=true`.

## Repository Layout

```text
.
├── main.go                         # protoc plugin entrypoint
├── internal/protoc-gen-gorm/       # generator implementation
├── option/
│   ├── gorm.proto                  # custom protobuf options
│   ├── gorm.pb.go                  # generated Go option definitions
│   └── genpb.sh                    # regenerates option/gorm.pb.go
├── example/
│   ├── model.proto                 # example protobuf input
│   ├── genpb.sh                    # regenerates example outputs
│   └── model/                      # generated example code and tests
└── third_party/proto/              # vendored protobuf descriptor files
```

## Requirements

- Go 1.20 or newer
- `protoc`
- `protoc-gen-go`
- GORM dependencies from `go.mod`

Install the standard Go protobuf generator if needed:

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

## Build and Install

Build a local binary:

```sh
make build
```

Install the plugin into your Go binary path:

```sh
make install
```

Keep dependencies tidy:

```sh
make tidy
```

## Usage

Import the GORM option definitions in your proto file:

```proto
import "option/gorm.proto";
```

Mark a message for model generation:

```proto
message User {
  option (gorm.set_table) = true;
  option (gorm.table_name) = "users";

  int32 id = 1 [(gorm.rules) = {int: {type: bigint}, primary_key: true}];
  string name = 2 [(gorm.rules) = {index: {name: "idx_name"}}];
  google.protobuf.Timestamp created_at = 3 [
    (gorm.rules) = {time: {type: datetime, auto_create_time: true}, not_null: true}
  ];
}
```

Run `protoc` with both the Go generator and this plugin:

```sh
protoc -I. \
  -I./third_party/proto \
  --go_out=. --go_opt paths=source_relative \
  --gorm_out=. --gorm_opt paths=source_relative \
  your_file.proto
```

To generate the extra helper and DAO files:

```sh
protoc -I. \
  -I./third_party/proto \
  --go_out=. --go_opt paths=source_relative \
  --gorm_out=. --gorm_opt paths=source_relative \
  --gorm_opt with_gorm_option=true \
  --gorm_opt with_gorm_dao=true \
  your_file.proto
```

`with_gorm_dao=true` is intended to be used together with `with_gorm_option=true`.

## Generated Files

For an input such as `model.proto`, the plugin can generate:

- `model.gorm.pb.go`: GORM model structs, tags, table names, conversion methods, and JSON scanner/valuer helpers for repeated and map fields.
- `model.gorm.option.pb.go`: shared GORM options and helper functions such as `Create`, `Save`, `First`, `Find`, `Count`, `Delete`, `Limit`, `Offset`, `Order`, and `TableName`.
- `model.gorm.dao.pb.go`: per-message DAO wrappers such as `UserDao` with typed CRUD methods.

Messages without `(gorm.set_table) = true` are skipped unless they are embedded by generated model types.

## Supported Options

Message options:

- `(gorm.set_table) = true`: enables generation for the message.
- `(gorm.table_name) = "name"`: overrides the default snake_case table name.
- `(gorm.disable_snake_case)`: defined in the proto options and reserved for naming behavior.

Field options:

- `(gorm.ignore_gorm_column) = true`: emits `gorm:"-"`.
- `(gorm.rules).column_name`: overrides the generated column name.
- `(gorm.rules).primary_key`: emits a primary key tag.
- `(gorm.rules).index` and `(gorm.rules).uniqueIndex`: emits index tags.
- `(gorm.rules).not_null`: emits `not null`.
- Type-specific rules are available for bool, int, float, time, binary, string, and enum fields.

Repeated fields and maps are stored as `longtext` and generated as custom Go types that implement `driver.Valuer` and `sql.Scanner` through JSON marshal/unmarshal helpers.

## Type Mapping Notes

- `google.protobuf.Timestamp` with time rules maps to `time.Time`, `datatypes.Date`, or `gorm.DeletedAt`.
- `int64`, `uint64`, and `string` fields can also be treated as time fields when time rules are provided.
- Enum fields can be stored as SQL `enum`, `string`, or `int`.
- Non-timestamp message fields are embedded with `embeddedPrefix`.
- Only one primary key and one soft-delete field are allowed per generated model hierarchy.

## Example

Regenerate the bundled example:

```sh
make build
sh example/genpb.sh
```

The example uses `example/model.proto` and emits generated code under `example/model/`.

## Testing

Run all tests:

```sh
go test ./...
```

The current example tests in `example/model/model.gorm.pb_test.go` connect to a local MySQL instance using hard-coded development credentials. If that database is unavailable, build and generation checks may still pass while the example integration tests fail.

## Development Notes

- Run `gofmt` on edited Go files.
- Regenerate `option/gorm.pb.go` after changing `option/gorm.proto`:

```sh
sh option/genpb.sh
```

- Regenerate example outputs after changing generator behavior:

```sh
make build
sh example/genpb.sh
```

- Generated files contain `// Code generated by protoc-gen-gorm. DO NOT EDIT.` and should normally be changed only through proto or generator updates.
