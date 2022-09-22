package services

import (
	"context"
	"net/http"

	"github.com/vandong9/go-grpc-auth-svc/pkg/db"
	"github.com/vandong9/go-grpc-auth-svc/pkg/models"
	"github.com/vandong9/go-grpc-auth-svc/pkg/pb"
	"github.com/vandong9/go-grpc-auth-svc/pkg/utils"
)

type Server struct {
	H   db.Handler
	Jwt utils.JwtWrapper
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var user models.User

	if result := s.H.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error == nil {
		return &pb.RegisterResponse{Status: http.StatusConflict, Error: "E-mail already exist"}, nil
	}

	user.Email = req.Email
	user.Password = utils.HashPassword(req.Password)
	s.H.DB.Create(&user)

	return &pb.RegisterResponse{Status: http.StatusCreated}, nil
}
