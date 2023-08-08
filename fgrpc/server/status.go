package server

import (
	"context"

	"github.com/bingoohuang/fastrest/fgrpc/status"
	"google.golang.org/grpc/peer"
)

type StatusServer struct {
	status.StatusServiceServer
}

func (s *StatusServer) Status(ctx context.Context, r *status.StatusRequest) (*status.StatusResponse, error) {
	return &status.StatusResponse{
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
