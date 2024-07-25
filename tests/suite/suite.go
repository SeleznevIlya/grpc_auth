package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/SeleznevIlya/grpc_auth/internal/config"
	grpc_authv1 "github.com/SeleznevIlya/protos/gen/go/grpc_auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T                        // Потребуется для вызова методов *testing.T внутри Suite
	Cfg        *config.Config         // Конфигурация приложения
	AuthClient grpc_authv1.AuthClient // Клиент для взаимодействия с gRPC-сервером
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	// key := "CONFIG_PATH"
	// if v := os.Getenv(key); v != "" {
	// 	return v
	// }

	cfg := config.MustLoadByPath("../config/local.yaml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})
	// cc, err := grpc.NewClient()
	cc, err := grpc.DialContext(context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials())) // Используем insecure-коннект для тестов
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: grpc_authv1.NewAuthClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
