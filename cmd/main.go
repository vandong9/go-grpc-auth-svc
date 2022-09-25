package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/vandong9/go-grpc-auth-svc/pkg/config"
	"github.com/vandong9/go-grpc-auth-svc/pkg/db"
	"github.com/vandong9/go-grpc-auth-svc/pkg/pb"
	"github.com/vandong9/go-grpc-auth-svc/pkg/services"
	"github.com/vandong9/go-grpc-auth-svc/pkg/utils"
)

func main() {
	// Load config
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("ailed at config", err)
	}

	// Init db connection
	h := db.Init(c.DBUrl)

	// Init JWT
	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "go-grpc-auth-svc",
		ExpirationHours: 24 * 365,
	}

	// start service with port
	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing", err)
	}

	fmt.Println("Auth Svc on", c.Port)

	// Init service and GRPC server
	s := services.Server{
		H:   h,
		Jwt: jwt,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
