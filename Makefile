.PHONY: dev build test migrate

dev:
	go run cmd/server/main.go

build:
	go build -o tmp/bin/brayat cmd/server/main.go

test:
	go test -v ./... -cover

migrate:
	@echo "Migrations are handled by GORM AutoMigrate on startup."
