package server

import (
	"context"
	"github.com/bingoohuang/fastrest/fgrpc/service"
	"google.golang.org/grpc/peer"
)

type GrpcServer struct {
	service.ServiceServer
}

func (s *GrpcServer) Encrypt(ctx context.Context, r *service.EncryptRequest) (*service.EncryptResponse, error) {
	return &service.EncryptResponse{
		Data: r.PlainText,
	}, nil
}

func (s *GrpcServer) Status(ctx context.Context, r *service.StatusRequest) (*service.StatusResponse, error) {
	return &service.StatusResponse{
		Status:  200,
		Message: "成功",
	}, nil
}

// GetRPCClientIP 检查上下文以检索客户机的ip地址
func GetRPCClientIP(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok {
		return p.Addr.String()
	}
	return ""
}
