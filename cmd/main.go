package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"
	"github.com/vandong9/go-grpc-auth-svc/pkg/config"
	"github.com/vandong9/go-grpc-auth-svc/pkg/db"
	"github.com/vandong9/go-grpc-auth-svc/pkg/endpoint"
	"github.com/vandong9/go-grpc-auth-svc/pkg/middleware"
	"github.com/vandong9/go-grpc-auth-svc/pkg/pb"
	"github.com/vandong9/go-grpc-auth-svc/pkg/router"
	"github.com/vandong9/go-grpc-auth-svc/pkg/services"
	"github.com/vandong9/go-grpc-auth-svc/pkg/utils"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	repository := db.Init(c.DBUrl)

	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "go-grpc-auth-svc",
		ExpirationHours: 24 * 365,
	}

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Auth Svc on", c.Port)

	s := services.GprcServer{
		H:   repository,
		Jwt: jwt,
	}

	// Create gprc Server
	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}

	// Create http server

	var (
		endpoints = endpoint.MakeEndpoints(services.HttpSever{})

		logger = logrus.New()
		h      = router.NewHandler(endpoints, logger)
		m      = middleware.NewMiddleware(logger)
		port   = c.Port

		server = http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: h.MakeHandlers(m),
		}
	)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Running HTTP server: %v", err)
	}

}
