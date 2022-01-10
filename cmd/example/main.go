package main

import (
	"context"
	"log"

	"github.com/reiot777/go-grpc-example/pkg/gurl"
)

func main() {
	ctx := context.Background()
	if b, err := gurl.Call(ctx, "packet.PingService/Ping", "", gurl.Addr(":10010")); err != nil {
		log.Println("packet.PingService/Ping", err)
	} else {
		log.Println("packet.PingService/Ping", string(b))
	}

	if b, err := gurl.Call(ctx, "packet.AccountService/CreateAccount", `{"email":"superman@gmail.com"}`, gurl.Addr(":10010")); err != nil {
		log.Println(err)
	} else {
		log.Println("packet.AccountService", string(b))
	}

	if b, err := gurl.Call(ctx, "packet.AccountService/CreateAccount", `{"emai":"superman@gmail.com"}`, gurl.Addr(":10010")); err != nil {
		log.Println(err)
	} else {
		log.Println("packet.AccountService", string(b))
	}
}
