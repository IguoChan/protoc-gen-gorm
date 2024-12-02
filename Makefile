# !/usr/bin/make

BIN=protoc-gen-gorm

tidy:
	@go mod tidy

build:
	@go build -o ${BIN} cmd/main.go

install:
	@go install cmd/protoc-gen-gorm.go