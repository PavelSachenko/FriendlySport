include .env
export
abs_path:=$(CURDIR)
up=0
down=0

start:
	docker-compose up -d

stop:
	docker-compose down

migrate-up:
	docker run --rm -v $(abs_path)/database/migration:/migrations --network host migrate/migrate \
						-path=/migrations/ \
						-database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@0.0.0.0:5433/${DB_DATABASE}?sslmode=disable" up $(up)

migrate-down:
	docker run --rm -v $(abs_path)/database/migration:/migrations --network host migrate/migrate \
						-path=/migrations/ \
						-database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@0.0.0.0:5433/${DB_DATABASE}?sslmode=disable" down $(down)


migration:
	migrate create -ext sql -dir database/migration friendly_sport


swagger:
	swag init -g cmd/main.go

run:
	go run app/main.go

build:
	docker-compose build --no-cache

migrate:
	docker run --rm -v $(abs_path)/database/migration:/migrations --network host migrate/migrate \
						-path=/migrations/ \
						-database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@0.0.0.0:5433/${DB_DATABASE}?sslmode=disable" up
test:
	echo ${DB_DATABASE}