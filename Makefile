DB_NAME=simple_bank
DB_USER=simple_bank
DB_PASSWORD=simple_bank
DB_URL=postgresql://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable

postgres:
	docker run --network bank-network --name postgressi -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine
createdb: createuser
	docker exec postgressi psql -U root -c "CREATE DATABASE ${DB_NAME};"
	docker exec postgressi psql -U root -c "GRANT ALL PRIVILEGES ON DATABASE ${DB_NAME} TO ${DB_USER};"
	docker exec postgressi psql -U root -c "GRANT ALL ON SCHEMA public TO ${DB_USER};"
createuser:
	docker exec postgressi psql -U root -c "CREATE USER ${DB_USER} WITH SUPERUSER PASSWORD '${DB_PASSWORD}';"
	docker exec postgressi psql -U root -c "ALTER ROLE ${DB_USER} SET client_encoding TO 'utf8';"
	docker exec postgressi psql -U root -c "ALTER ROLE ${DB_USER} SET default_transaction_isolation TO 'read committed';"
	docker exec postgressi psql -U root -c "ALTER ROLE ${DB_USER} SET timezone TO 'UTC';"
dropdb:
	docker exec postgressi dropdb simple_bank
migrateup:
	migrate -path db/migration -database "${DB_URL}" -verbose up
migrateup1:
	migrate -path db/migration -database "${DB_URL}" -verbose up 1
migratedown:
	migrate -path db/migration -database "${DB_URL}" -verbose down
migratedown1:
	migrate -path db/migration -database "${DB_URL}" -verbose down 1
sqlc:
	sqlc generate
db_docs:
	dbdocs build doc/db.dbml
db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/umtdemr/simplebank/db/sqlc Store
proto:
	rm -f pb/*.go
	protoc  --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
        --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
        proto/*.proto

.PHONY: postgres createdb createuser dropdb migrateup migrateup1 migratedown migratedown1 sqlc db_docs db_schema server mock proto
