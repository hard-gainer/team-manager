package main

import (
	"context"
	"log"

	"github.com/hard-gainer/team-manager/internal/auth"
	"github.com/hard-gainer/team-manager/internal/config"
	db "github.com/hard-gainer/team-manager/internal/db/sqlc"
	"github.com/hard-gainer/team-manager/internal/mail"
	"github.com/hard-gainer/team-manager/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.MustLoad()
	mailer := mail.NewMailer(cfg)

	connPool, err := pgxpool.New(context.Background(), cfg.DBConfig.StoragePath)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	authConn, err := grpc.NewClient(
		cfg.GRPCConfig.AuthAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("cannot connect to auth service: %v", err)
	}
	defer authConn.Close()

	authClient := auth.NewAuthClient(authConn)

	store := db.NewStore(connPool)
	server := service.NewServer(cfg, store, authClient, mailer)

	err = server.Start(":8080")
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
