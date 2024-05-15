migrate create -seq -ext=.sql -dir=./db/migrations init_db_schema


migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up


migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down



migrate create -seq -ext=.sql -dir=./db/migrations add_users

migrate create -seq -ext=.sql -dir=./db/migrations add_sessions