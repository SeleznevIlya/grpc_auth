package app

import (
	"log/slog"
	"time"

	grpcApp "github.com/SeleznevIlya/grpc_auth/internal/app/grpc"
	"github.com/SeleznevIlya/grpc_auth/internal/services/auth"
	"github.com/SeleznevIlya/grpc_auth/internal/storage/postgres"
)

type App struct {
	GRPCSrv *grpcApp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := postgres.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcApp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
