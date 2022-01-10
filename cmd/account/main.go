package main

import (
	"context"

	"github.com/reiot777/go-grpc-example/service/account"
)

func main() {
	srv := account.Server{
		Host: "0.0.0.0",
		Port: 10010,
	}

	srv.Serve(context.Background())
}
