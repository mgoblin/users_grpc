package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/joho/godotenv"

	usersServer "mgoblin/users_grpc/internal/users/server"
	"mgoblin/users_grpc/internal/users/server/config"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	conf := config.New()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.UsersServer.TCPPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}
	log.Printf("Listening port %d", conf.UsersServer.TCPPort)

	grpcServer := usersServer.New(conf)

	ctx := context.Background()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			grpcServer.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
