DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

run:
	nodemon --exec go run main.go --signal SIGTERM


fetch-postgres-image:
	docker pull postgres:12-alpine

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

mysql:
	docker run --name mysql8 -p 3306:3306  -e MYSQL_ROOT_PASSWORD=secret -d mysql:8

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mock_db github.com/saalikmubeen/backend-masterclass-go/db/sqlc Store


# Generates a random 32 character string to be used as TOKEN_SYMMETRIC_KEY
# for PASETO token generation
token_key:
	openssl rand -hex 64 | head -c 32

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock token_key