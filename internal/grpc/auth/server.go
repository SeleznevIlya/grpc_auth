package auth

import (
	"context"

	grpc_authv1 "github.com/SeleznevIlya/protos/gen/go/grpc_auth"
	"google.golang.org/grpc"
)

type serverApi struct {
	grpc_authv1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	grpc_authv1.RegisterAuthServer(gRPC, &serverApi{})
}

func (s *serverApi) Login(ctx context.Context, req *grpc_authv1.LoginRequest) (*grpc_authv1.LoginResponse, error) {
	panic("implement me")
}

func (s *serverApi) Register(ctx context.Context, req *grpc_authv1.RegisterRequest) (*grpc_authv1.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverApi) IsAdmin(ctx context.Context, req *grpc_authv1.IsAdminRequest) (*grpc_authv1.IsAdminResponse, error) {
	panic("implement me")
}
