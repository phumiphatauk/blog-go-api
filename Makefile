BINARY_NAME=blog-go-api
ENVIRONMENT=development
DB_HOST=localhost
DB_PORT=5432
DB_NAME=blog?sslmode=disable
DB_USER=postgres
DB_PASSWORD=secret
DB_URL=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)
URL_LOCALHOST=http://localhost:5173
TOKEN_SYMMETRIC_KEY=12345678901234567890123456789012
MIGRATION_URL=file://db/migration
DB_DRIVER=postgres
MINIO_ENDPOINT=http://localhost:9000
MINIO_ACCESS_KEY_ID=bTUfPQYWcSX7fgXL2hQj
MINIO_SECRET_ACCESS_KEY=PACrKXIFke9J6bq28z8q7Qqg6ihLp0QFH5yQ4ajk
MINIO_USE_SSL=true
MINIO_BUCKET_NAME=blog
MINIO_URL_RESULT=http://localhost:9000/blog/
EMAIL_SENDER_NAME=sender
EMAIL_SENDER_ADDRESS=yourself@mail.com
EMAIL_SENDER_PASSWORD=yourselfpassword

## postgres: Start PostgreSQL container
postgres:
	docker compose up postgres -d

## minio: Start MinIO container
minio:
	docker compose up minio -d

## postgresdown: Stop PostgreSQL container
postgresdown:
	docker compose down postgres

## migrate: Create a new migration file
migrate:
	migrate create -ext sql -dir db/migration -seq $(name)

## migrateup: Apply all migrations
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

## migrateup_no: Apply the next migration
migrateup_no:
	migrate -path db/migration -database "$(DB_URL)" -verbose up $(no)

## migratedown: Rollback all migrations
migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

## migratedown_no: Rollback the previous migration
migratedown_no:
	migrate -path db/migration -database "$(DB_URL)" -verbose down $(no)

migrate_version:
	migrate -path db/migration -database "$(DB_URL)" -verbose version

migrate_goto:
	migrate -path db/migration -database "$(DB_URL)" -verbose goto $(no)

migrate_force:
	migrate -path db/migration -database "$(DB_URL)" -verbose force $(no)

## sqlc: Generate code from SQL queries
sqlc:
	sqlc generate

## build: Build binary
build:
	@echo "Building back end..."
	go build -o ${BINARY_NAME} .
	@echo "Binary built!"

## run: builds and runs the application
run: build
	@echo "Starting back end..."
	@env HTTP_SERVER_ADDRESS=0.0.0.0:8080 \
	ENVIRONMENT=${ENVIRONMENT} \
	DB_SOURCE=${DB_URL} \
	GIN_MODE=debug \
	URL_LOCALHOST=${URL_LOCALHOST} \
	TOKEN_SYMMETRIC_KEY=${TOKEN_SYMMETRIC_KEY} \
	MIGRATION_URL=${MIGRATION_URL} \
	DB_DRIVER=${DB_DRIVER} \
	MINIO_ENDPOINT=${MINIO_ENDPOINT} \
	MINIO_ACCESS_KEY_ID=${MINIO_ACCESS_KEY_ID} \
	MINIO_SECRET_ACCESS_KEY=${MINIO_SECRET_ACCESS_KEY} \
	MINIO_USE_SSL=${MINIO_USE_SSL} \
	MINIO_BUCKET_NAME=${MINIO_BUCKET_NAME} \
	MINIO_URL_RESULT=${MINIO_URL_RESULT} \
	EMAIL_SENDER_NAME=${EMAIL_SENDER_NAME} \
	EMAIL_SENDER_ADDRESS=${EMAIL_SENDER_ADDRESS} \
	EMAIL_SENDER_PASSWORD=${EMAIL_SENDER_PASSWORD} \
	./${BINARY_NAME} &
	@echo "Back end started!"

## stop: stops the running application
stop:
	@echo "Stopping back end..."
	@if [ -f "${BINARY_NAME}" ]; then \
		pkill -SIGTERM -f "./${BINARY_NAME}"; \
		rm ${BINARY_NAME}; \
	fi
	@echo "Stopped back end!"

## restart: stops and starts the running application
restart: stop run

swag:
	swag init

.PHONY: postgres minio postgresdown migrate migrateup migrateup_no migratedown migratedown_no migrate_version migrate_goto migrate_force sqlc build run stop restart swag
