package main

import (
	"context"
	"log"

	db "github.com/hard-gainer/task-tracker/internal/db/sqlc"
	"github.com/hard-gainer/task-tracker/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	connString = "postgresql://root:secret@localhost:5432/task_tracker?sslmode=disable"
)

func main() {

	connPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	store := db.NewStore(connPool)
	server := service.NewServer(store)

	err = server.Start(":8080")
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
