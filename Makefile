include .env
export


.PHONY: fmt lint install up down stop mock-database sqlc docs migrate run test test-integration test-benchmark coverage

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
	sqlc generate -f internal/infrastructure/database/sqlc/sqlc.yaml

sqlc-vet:
	sqlc vet -f internal/infrastructure/database/sqlc/sqlc.yaml

docs:
	make lint
	swag init -g cmd/main.go -o api 

# Executa somente testes unitários
test:
	go test ./... -v -count=1 -race -cover -coverprofile=./tmp/coverage.unit.out

# Executa somente testes de integração
test-integration:
	go test -tags=integration ./... -v -count=1 -race -cover -coverprofile=./tmp/coverage.integration.out
#go test -tags=integration ./...
# Executa benchmarks
test-benchmark:
	go test ./... -v -run=^$ -bench=. -benchmem

# Gera relatório de cobertura em HTML
coverage:
	go tool cover -html=./tmp/coverage.integration.out -o ./tmp/coverage.html

# Limpa arquivos de coverage
clean:
	rm -f ./tmp/coverage.*

dev:
	make docs
	air

run:
	make docs
	go run -race cmd/main.go

migrate:
	@ migrate -source file://internal/database/migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&search_path=public" up

migration:
	@read -p "Enter migration name: " name; \
		migrate create -ext sql -dir internal/infrastructure/database/migrations -seq $$name

migration-up: 
	@ migrate -path internal/infrastructure/database/migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&search_path=public" -verbose up

migration-down: 
	@ migrate -path internal/infrastructure/database/migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&search_path=public" -verbose down

migration-fix: 
	@read -p "Enter migration version: " version; \
	migrate -path internal/infrastructure/database/migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&search_path=public" force $$version

infra:
	go run scripts/create_infra.go