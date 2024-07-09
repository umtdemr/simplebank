DB_NAME=simple_bank
DB_USER=simple_bank
DB_PASSWORD=simple_bank

postgres:
	docker run --name postgressi -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine
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
	migrate -path db/migration -database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable" -verbose down
sqlc:
	sqlc generate

test:
	go test -v -cover ./...
server:
	go run main.go

.PHONY: postgres createdb createuser dropdb migrateup migratedown sqlc server
