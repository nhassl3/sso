package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/nhassl3/sso-app/internals/domain/services/auth"
	authgrpc "github.com/nhassl3/sso-app/internals/grpc/auth"
	"google.golang.org/grpc"
)

const (
	opStart = "grpcapp.Start"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(log *slog.Logger, port int, authObj *auth.Auth) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authObj)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustStart launching gRPC server on low level protocol - TCP
func (s *App) MustStart() {
	log := s.log.With(slog.String("op", opStart), slog.Int("port", s.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		panic(fmt.Errorf("%s: %w", opStart, err))
	}

	log.Info("server started", slog.String("address", l.Addr().String()))

	if err := s.gRPCServer.Serve(l); err != nil {
		panic(fmt.Errorf("%s: %w", opStart, err))
	}
}

// Stop graceful stops gRPC server
func (s *App) Stop() {
	s.gRPCServer.GracefulStop()
}
