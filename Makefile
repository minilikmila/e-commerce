APP_NAME=ecommerce
BIN=./cmd/server


.PHONY: run build test test-handlers fmt docker-build up down swagger

run:
	go run $(BIN)

build:
	go build -o bin/$(APP_NAME) $(BIN)

test:
	go test ./internal/adapter/handler/... -v

fmt:
	go fmt ./...

docker-build:
	docker build -t bin/$(APP_NAME):latest .

up:
	docker-compose up --build

up-detach-logs:
	docker-compose up -d --build

down:
	docker-compose down
 
swagger:
	@echo "Installing swag..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Generating Swagger docs..."
	@swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal --parseDepth 10
	@echo "Swagger docs generated in docs/"
	@echo "Check docs/swagger.json for generated paths"

