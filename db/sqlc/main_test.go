package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDb *pgxpool.Pool

const dbURL = "postgresql://simple_bank:simple_bank@localhost:5432/simple_bank?sslmode=disable"

func TestMain(m *testing.M) {
	var err error
	testDb, err = pgxpool.New(context.Background(), dbURL)

	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
