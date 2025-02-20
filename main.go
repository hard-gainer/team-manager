package main

import (
	"context"
	"log"

	"github.com/hard-gainer/task-tracker/internal/auth"
	db "github.com/hard-gainer/task-tracker/internal/db/sqlc"
	"github.com/hard-gainer/task-tracker/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	connString = "postgresql://root:secret@localhost:5432/task_tracker?sslmode=disable"
	authAddr   = "localhost:44044"
)

func main() {
	connPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	authConn, err := grpc.NewClient(
		authAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("cannot connect to auth service: %v", err)
	}
	defer authConn.Close()

	authClient := auth.NewAuthClient(authConn)

	store := db.NewStore(connPool)
	server := service.NewServer(store, authClient)

	err = server.Start(":8080")
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
