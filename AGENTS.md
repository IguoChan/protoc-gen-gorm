# Repository Guidelines

## Project Structure & Module Organization

This repository is a Go `protoc` plugin that generates GORM model and DAO code from protobuf definitions.

- `main.go` contains the plugin entrypoint and command-line flags.
- `internal/protoc-gen-gorm/` contains generator implementation and option handling.
- `option/` defines the custom protobuf options in `gorm.proto` and the generated `gorm.pb.go`.
- `example/` contains sample proto input, generation script, generated outputs, and example tests under `example/model/`.
- `third_party/proto/` vendors protobuf compiler descriptors used during generation.

## Build, Test, and Development Commands

- `make build` builds the local `protoc-gen-gorm` binary at the repository root.
- `make install` installs the plugin with `go install`.
- `make tidy` runs `go mod tidy` to normalize module dependencies.
- `go test ./...` runs all Go tests. Note that `example/model/model.gorm.pb_test.go` expects a local MySQL server and matching credentials.
- `sh option/genpb.sh` regenerates Go code for `option/gorm.proto`.
- `sh example/genpb.sh` regenerates example protobuf, GORM, option, and DAO outputs. Build or install the plugin first so `protoc` can find `protoc-gen-gorm`.

## Coding Style & Naming Conventions

Use standard Go formatting: run `gofmt` on edited `.go` files and keep imports organized by `goimports` or the Go toolchain. Follow existing package naming such as `protoc_gen_gorm` for the internal generator package. Generated files use suffixes like `.gorm.pb.go`, `.gorm.option.pb.go`, and `.gorm.dao.pb.go`; do not hand-edit generated output unless updating expected generated examples.

## Testing Guidelines

Add focused Go tests next to the package or generated example they exercise, using `TestXxx` names. For generator changes, prefer updating `example/model.proto`, regenerating outputs, and adding assertions around generated behavior. Treat database-backed tests as integration tests and document required services or credentials when adding new ones.

## Commit & Pull Request Guidelines

The current history is minimal (`init`, `Initial commit`), so keep commits short, imperative, and scoped, for example `add dao option generation`. Pull requests should include a concise summary, the generator behavior affected, commands run, and any regeneration steps. Link related issues when available and call out required local services such as MySQL for tests.

## Agent-Specific Notes

Avoid unrelated rewrites of generated files. When changing protobuf options or generator logic, regenerate the affected artifacts and review the generated diff carefully.
