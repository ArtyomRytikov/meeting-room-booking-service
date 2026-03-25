up:
	docker compose up --build

down:
	docker compose down

down-v:
	docker compose down -v

test:
	go test ./... -count=1

coverage:
	go test ./... -coverprofile=coverage.out -count=1
	go tool cover -func=coverage.out

lint:
	golangci-lint run

seed:
	@echo "Step 1: get admin token:"
	@echo "curl -X POST http://localhost:8080/dummyLogin -H \"Content-Type: application/json\" -d \"{\\\"role\\\":\\\"admin\\\"}\""
	@echo ""
	@echo "Step 2: use token in commands below:"
	@echo "curl -X POST http://localhost:8080/rooms/create -H \"Authorization: Bearer YOUR_TOKEN\" -H \"Content-Type: application/json\" -d \"{\\\"name\\\":\\\"Seed Room\\\",\\\"description\\\":\\\"seed\\\",\\\"capacity\\\":5}\""

help:
	@echo "Available commands:"
	@echo " make up        - start app with docker"
	@echo " make down      - stop containers"
	@echo " make down-v    - stop and remove volumes"
	@echo " make test      - run tests"
	@echo " make coverage  - run tests with coverage"
	@echo " make lint      - run linter"
	@echo " make seed      - print seed instructions"