package app

import (
	"log/slog"
	"time"

	"github.com/nhassl3/sso-app/internals/app/grpcapp"
	"github.com/nhassl3/sso-app/internals/domain/services/auth"
	"github.com/nhassl3/sso-app/internals/storage/sqlite"
)

type App struct {
	GRPCServer *grpcapp.App
}

func NewApp(log *slog.Logger, gRPCPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.NewStorage(storagePath)
	if err != nil {
		panic(err)
	}

	authObj := auth.NewAuth(log, storage, storage, storage, tokenTTL)

	gRPCApp := grpcapp.NewApp(log, gRPCPort, authObj)

	return &App{
		GRPCServer: gRPCApp,
	}
}
