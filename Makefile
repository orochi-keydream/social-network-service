.PHONY: run
run:
	go run ./cmd/app/main.go

.PHONY: gen-profiles
gen-profiles:
	go run ./cmd/gen-profiles/main.go

.PHONY: compose-up
compose-up:
	docker compose -f ./docker-compose.yml up -d

.PHONY: compose-down
compose-down:
	docker compose -f ./docker-compose.yml down

.PHONY: migrate-up
migrate-up:
	goose -dir ./migrations/ postgres "host=localhost port=5432 user=postgres password=123 dbname=social_network_db" up

.PHONY: migrate-down
migrate-down:
	goose -dir ./migrations/ postgres "host=localhost port=5432 user=postgres password=123 dbname=social_network_db" down

.PHONY: gen-swag
gen-swag:
	swag init -g ./internal/app/app.go

.PHONY: add-migration
add-migration:
	goose -dir ./migrations/ create untitled sql