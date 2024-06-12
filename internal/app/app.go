package app

import (
	"log/slog"
	"time"

	grpcApp "github.com/SeleznevIlya/grpc_auth/internal/app/grpc"
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
	// TODO: инициализировать хранилище (storage)

	// TODO: init auth service (auth)

	grpcApp := grpcApp.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
