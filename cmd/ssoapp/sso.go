package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nhassl3/sso-app/internals/app"
	"github.com/nhassl3/sso-app/internals/config"
	"github.com/nhassl3/sso-app/internals/lib/logger/handlers/slogpretty"
)

const (
	_              = iota
	EnvLocal       // local
	EnvDevelopment // development
	EnvProduction  // production
)

func main() {
	// Configuration load
	cfg := config.MustLoad()

	// Logger load
	log := setupLogger(cfg.Env)

	// Application load
	application := app.NewApp(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCServer.MustStart()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	sign := <-stop

	application.GRPCServer.Stop()
	// Stop every service and components (databases for ex) separately

	log.Info("Application stopped", slog.String("sign", sign.String()))
}

func setupLogger(env uint8) *slog.Logger {
	var level slog.Leveler
	switch env {
	case EnvLocal, EnvDevelopment:
		level = slog.LevelDebug
	case EnvProduction:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: level,
		},
	}

	return slog.New(opts.NewPrettyHandler(os.Stdout))
}
