.PHONY: migrate-up migrate-down docker-up docker-down docker-db wait-postgres test input-topic processed-topic topic-list

include .env
export

env:
	cp .env.example .env

run:
	go run ./cmd/main.go

message:
	go run ./cmd/script/write_message.go -url="example.com" -pubDate=1 -fetchTime=3 -text="some text" -firstFetchTime=2

create-topics:
	docker-compose exec kafka kafka-topics --create --topic ${KAFKA_IN_TOPIC} --bootstrap-server ${KAFKA_BROKER_HOST}:${KAFKA_PORT} --partitions 3 --replication-factor 1
	docker-compose exec kafka kafka-topics --create --topic ${KAFKA_OUT_TOPIC} --bootstrap-server ${KAFKA_BROKER_HOST}:${KAFKA_PORT} --partitions 3 --replication-factor 1

topic-list:
	docker-compose exec kafka kafka-topics --list --bootstrap-server ${KAFKA_BROKER_HOST}:${KAFKA_PORT}

docker-up:
	docker-compose up -d
	make wait-postgres migrate-up
	make create-topics

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

proto-to-golang:
	protoc \
    --go_out=. \
    --go-grpc_out=. \
    --proto_path=docs \
    docs/tdocument.proto

# Run vim ~/.bash_profile
# Add: export GO_PATH=~/go export PATH=$PATH:/$GO_PATH/bin
# Run source ~/.bash_profile
