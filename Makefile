include .env
export


.PHONY: fmt lint install up down stop mock-database sqlc docs test test-integration migrate run test test-integration test-benchmark coverage

lint:
	@gofmt -s -w .

install:
	go mod tidy

up:
	docker compose -f deployments/docker-compose.yaml --env-file .env up -d

down:
	docker compose -f deployments/docker-compose.yaml down

stop:
	docker compose -f deployments/docker-compose.yaml stop

mock-db:
	go run scripts/populate_local_db.go

sqlc:
	sqlc generate -f configs/sqlc.yaml

sqlc-vet:
	sqlc vet -f configs/sqlc.yaml

docs:
	make lint
	swag init -g cmd/main.go -o api 

# Executa somente testes unitários
test-unit:
	go test ./... -v -count=1 -race -run "TestUnit" -cover -coverprofile=coverage.unit.out

# Executa somente testes de integração
test-integration:
	go test ./... -v -count=1 -race -run "TestIntegration" -cover -coverprofile=coverage.integration.out

# Executa benchmarks
test-benchmark:
	go test ./... -v -run=^$ -bench=. -benchmem

# Gera relatório de cobertura em HTML
coverage:
	go tool cover -html=coverage.out -o coverage.html

# Limpa arquivos de coverage
clean:
	rm -f coverage.* 

#Executa todos os testes
test:
	go test ./... -v -count=1 -race -cover -coverprofile=coverage.out

# Executa todos os testes e gera relatório de cobertura
test-all: clean test coverage

dev:
	make docs
	air 

run:
	make docs
	go run cmd/main.go

migrate:
	@ migrate -source file://internal/database/migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&search_path=public" up

migration:
	@read -p "Enter migration name: " name; \
		migrate create -ext sql -dir internal/database/migrations -seq $$name

migration-up: 
	@ migrate -path internal/database/migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&search_path=public" -verbose up

migration-down: 
	@ migrate -path internal/database/migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&search_path=public" -verbose down

migration-fix: 
	@read -p "Enter migration version: " version; \
	migrate -path internal/database/migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&search_path=public" force $$version
