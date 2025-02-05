package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	connString = "postgresql://root:secret@localhost:5432/task_tracker?sslmode=disable"
)

var testStore Store

func TestMain(m *testing.M) {

	connPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
