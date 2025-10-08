package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/nhassl3/sso-app/internals/config"
	ssov1 "github.com/nhassl3/sso-contracts/generated/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	gRPCHost = "localhost"

	EmptyAppID int32 = 0
	AppID      int32 = 2
	AppSecret        = "test-secret"

	passDefaultLen = 10
	DeltaSecond    = 1
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

func NewSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/local_tests.yaml")

	ctx, cancelContext := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelContext()
	})

	cc, err := grpc.NewClient(
		net.JoinHostPort(gRPCHost, strconv.Itoa(cfg.GRPC.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed %v", err)
	}

	return ctx, &Suite{
		t,
		cfg,
		ssov1.NewAuthClient(cc),
	}
}

func (s *Suite) NewPassword() string {
	return gofakeit.Password(true, false, true, false, false, passDefaultLen)
}

func (s *Suite) NewEmail() string {
	return gofakeit.Email()
}
