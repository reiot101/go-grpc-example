package account

import (
	"context"
	"net"
	"strconv"

	"github.com/reiot777/go-grpc-example/packet"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	Host string
	Port int

	logger *zap.Logger
}

func (s *Server) Serve(ctx context.Context) {
	s.logger, _ = zap.NewDevelopment()
	defer s.logger.Sync()

	svc := Service{
		Logger: s.logger,
	}

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(s.Logging()),
	)

	go func() {
		defer srv.GracefulStop()
		<-ctx.Done()
	}()

	packet.RegisterPingServiceServer(srv, &svc)
	packet.RegisterAccountServiceServer(srv, &svc)
	reflection.Register(srv)

	lis, err := net.Listen("tcp", net.JoinHostPort(s.Host, strconv.Itoa(s.Port)))
	if err != nil {
		s.logger.Fatal("failed to listen ", zap.Error(err))
	}
	defer lis.Close()

	s.logger.Info("gRPC server serving", zap.String("listen", lis.Addr().String()))

	if err := srv.Serve(lis); err != nil {
		s.logger.Fatal("failed to serve ", zap.Error(err))
	}
}
