# !/usr/bin/make

BIN=protoc-gen-gorm

tidy:
	@go mod tidy

build:
	@go build -o ${BIN}

install:
	@go install