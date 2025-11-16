.PHONY: build loadtest e2e-full dep fmt lint swag new-mig mig-up mig-down-by-one mig-down-full clean env-up env-down

build:
	go mod tidy
	go build -o bin/app cmd/app/main.go

loadtest:
	k6 run loadtest/load_test.js

e2e-full:
	go test -v ./test/e2e/... -timeout=10m

swag: build dep
	swag init --pd --parseInternal --parseVendor --parseDepth 2 -g http.go -o etc/api -d internal/app/delivery/http/impl -ot yml

lint:
	golangci-lint cache clean
	golangci-lint run -c .golangci.yml

fmt:
	golangci-lint fmt
	swag fmt

dep:
	go install github.com/swaggo/swag/cmd/swag@latest

clean:
	@rm -rf bin*

env-up:
	docker compose -f docker-compose.yml up -d --build

env-down:
	docker compose -f docker-compose.yml down

new-mig:
	goose -dir migrations/sql create $(NAME) sql

mig-up:

mig-down-by-one:

mig-down-full:
