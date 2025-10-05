package app

import (
	"log/slog"
	"time"

	"github.com/nhassl3/sso-app/internals/app/grpcapp"
	"github.com/nhassl3/sso-app/internals/domain/services/auth"
)

type App struct {
	GRPCServer *grpcapp.App
}

func NewApp(log *slog.Logger, gRPCPort int, storagePath string, tokenTTL time.Duration) *App {
	// TODO: initialize storage
	//storage :=

	// TODO: initialize auth service
	authObj := auth.NewAuth(log, storage.UserSaver, storage.UserProvider, storage.AppProvider, tokenTTL)

	gRPCApp := grpcapp.NewApp(log, gRPCPort, authObj)
	_ = gRPCApp

	return &App{
		GRPCServer: gRPCApp,
	}
}
