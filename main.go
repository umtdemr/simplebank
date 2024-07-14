package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/umtdemr/simplebank/api"
	db "github.com/umtdemr/simplebank/db/sqlc"
	"github.com/umtdemr/simplebank/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	fmt.Println(config.DBSource, config.ServerAddress)
	conn, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create a server")
	}
	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start the server")
	}
}
