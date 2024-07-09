package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/umtdemr/simplebank/api"
	db "github.com/umtdemr/simplebank/db/sqlc"
	"log"
)

const dbURL = "postgresql://simple_bank:simple_bank@localhost:5432/simple_bank?sslmode=disable"
const serverAddress = "0.0.0.0:8080"

func main() {
	var err error
	conn, err := pgxpool.New(context.Background(), dbURL)

	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(serverAddress)

	if err != nil {
		log.Fatal("cannot start the server")
	}
}
