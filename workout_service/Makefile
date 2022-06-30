include .env
export
abs_path:=$(CURDIR)
up=0
down=0

proto:
	protoc -I=. --go_out=paths=source_relative:.  --go-grpc_out=paths=source_relative:.  pkg/pb/**/*.proto \
	&& sed -i "s/,omitempty//" pkg/pb/**/*.pb.go


server:
	go run cmd/main.go

migrate-up:
	docker run --rm -v $(abs_path)/migration:/migrations --network host migrate/migrate \
						-path=/migrations/ \
						-database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@0.0.0.0:5435/${DB_DATABASE}?sslmode=disable" up

migrate:
	docker run --rm -v $(abs_path)/migration:/migrations --network host migrate/migrate \
						-path=/migrations/ \
						-database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@0.0.0.0:5435/${DB_DATABASE}?sslmode=disable" up

migrate-create:
	migrate create -ext sql -dir migration friendly_sport

elastic-migrate:
	./elastic_migrations.sh