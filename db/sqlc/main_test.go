package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"testing"
)

var testQueries *Queries

const dbURL = "postgresql://simple_bank:simple_bank@localhost:5432/simple_bank?sslmode=disable"

func TestMain(m *testing.M) {
	conn, err := pgx.Connect(context.Background(), dbURL)

	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
