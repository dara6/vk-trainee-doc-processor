.PHONY: migrate-up migrate-down docker-up docker-down docker-db wait-postgres test

include .env
export

env:
	cp .env.example .env

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-db:
	docker-compose exec ${DB_SERVICE_NAME} psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}

wait-postgres:
	echo "Waiting for PostgreSQL to be ready..."
	while ! docker-compose exec db pg_isready -U postgers; do \
  		sleep 2; \
 	done

migrate-up:
	migrate \
	-path ${MIGRATION_DIR} \
	-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" \
	up

migrate-down:
	migrate \
	-path ${MIGRATION_DIR} \
	-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" \
	down --all

test:
	docker-compose up -d db
	make wait-postgres migrate-down migrate-up
	go test -v ./...
	docker-compose down
