package account

import (
	"context"
	"strings"
	"time"

	"github.com/reiot777/go-grpc-example/packet"
	"github.com/rs/xid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	_ packet.AccountServiceServer = (*Service)(nil)
	_ packet.PingServiceServer    = (*Service)(nil)
)

type Service struct {
	Logger *zap.Logger
}

func (s *Service) CreateAccount(ctx context.Context, in *packet.CreateAccountRequest) (*packet.CreateAccountResponse, error) {
	newAccount := &packet.Account{
		Id:        xid.New().String(),
		Owner:     strings.Split(in.Email, "@")[0],
		Email:     in.Email,
		Amount:    100,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	return &packet.CreateAccountResponse{
		Data: newAccount,
	}, nil
}

func (s *Service) GetAccount(ctx context.Context, in *packet.GetAccountRequest) (*packet.GetAccountResponse, error) {
	s.Logger.Debug("GetAccount", zap.String("Request", in.String()))

	account := &packet.Account{
		Id:        "xx",
		Owner:     "sushi",
		Email:     "sushi@gmail.com",
		Amount:    100,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	return &packet.GetAccountResponse{
		Data: account,
	}, nil
}

func (s *Service) ProduceAmount(ctx context.Context, in *packet.ProduceAmountRequest) (*packet.ProduceAmountResponse, error) {
	return &packet.ProduceAmountResponse{}, nil
}

func (s *Service) ConsumeAmount(ctx context.Context, in *packet.ConsumeAmountRequest) (*packet.ConsumeAmountResponse, error) {
	return &packet.ConsumeAmountResponse{}, nil
}

func (s *Service) Ping(ctx context.Context, _ *emptypb.Empty) (*packet.PingResponse, error) {
	return &packet.PingResponse{
		Ts: time.Now().Unix(),
	}, nil
}
