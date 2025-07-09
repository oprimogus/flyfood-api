include .env
export

.PHONY: lint
lint:
	@gofmt -s -w .

.PHONY: install
install:
	go mod tidy

.PHONY: up
up:
	docker compose -f deployments/docker-compose.yaml --env-file .env up -d

.PHONY: down
down:
	docker compose -f deployments/docker-compose.yaml down

.PHONY: stop
stop:
	docker compose -f deployments/docker-compose.yaml stop

.PHONY: sqlc
sqlc:
	go tool sqlc generate -f internal/infrastructure/database/sqlc/sqlc.yaml

.PHONY: docs
docs:
	make install
	make lint
	go tool swag fmt -d ./
	go tool swag init -g cmd/main.go -o api --v3.1 --parseInternal --parseDependency

# Executa somente testes unitários
.PHONY: test
test:
	go test ./internal... -v -count=1 -race -cover -coverprofile=./tmp/coverage.unit.out

# Executa somente testes de integração
.PHONY: test-integration
test-integration:
	go test -tags=integration ./internal... -v -count=1 -race -cover -coverprofile=./tmp/coverage.integration.out

.PHONY: test-ci
test-ci:
	go test ./internal/... -v -count=1 -race

.PHONY: test-integration-ci
test-integration-ci:
	go test -tags=integration ./internal... -v -count=1 -race

# Executa benchmarks
.PHONY: test-benchmark
test-benchmark:
	go test ./... -v -run=^$ -bench=. -benchmem

# Gera relatório de cobertura em HTML
.PHONY: coverage
coverage:
	go tool cover -html=./tmp/coverage.integration.out -o ./tmp/coverage.html

# Limpa arquivos de coverage
.PHONY: clean
clean:
	rm -f ./tmp/coverage.*

.PHONY: dev
dev:
	go tool air

.PHONY: run
run:
	make docs
	go run -race cmd/main.go

.PHONY: migrate
migrate:
	@ migrate -source file://internal/infrastructure/database/postgres/migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&search_path=public" up

.PHONY: migration
migration:
	@read -p "Enter migration name: " name; \
		migrate create -ext sql -dir internal/infrastructure/database/postgres/migrations -seq $$name

.PHONY: migration-up
migration-up: 
	@ migrate -path internal/infrastructure/database/postgres/migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&search_path=public" -verbose up

.PHONY: migration-down
migration-down: 
	@ migrate -path internal/infrastructure/database/postgres/migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&search_path=public" -verbose down

.PHONY: migration-fix
migration-fix: 
	@read -p "Enter migration version: " version; \
	migrate -path internal/infrastructure/database/postgres/migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&search_path=public" force $$version

.PHONY: infra
infra:
	go run scripts/create_infra/main.go

.PHONY: mock-db
mock-db:
	go run scripts/populate_local_db/main.go