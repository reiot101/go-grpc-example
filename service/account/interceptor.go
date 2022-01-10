package account

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func (s *Server) Logging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		s.logger.Debug(info.FullMethod, zap.Any("request", req))
		resp, err = handler(ctx, req)
		s.logger.Debug(info.FullMethod, zap.Any("response", resp))
		return
	}
}
