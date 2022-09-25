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

// Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
// Login(context.Context, *LoginRequest) (*LoginResponse, error)
// Validate(context.Context, *ValidateRequest) (*ValidateResponse, error)

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

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user models.User

	if result := s.H.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error != nil {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)

	if !match {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	token, _ := s.Jwt.GenerateToken(user)

	return &pb.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	_, err := s.Jwt.ValidateToken(req.Token)

	if err != nil {
		return &pb.ValidateResponse{Status: http.StatusBadRequest, Error: "User not found"}, nil
	}

	var user models.User

	if result := s.H.DB.Where(&models.User{}).First(&user); result.Error != nil {
		return &pb.ValidateResponse{Status: http.StatusNotFound, Error: "User not found"}, nil
	}

	return &pb.ValidateResponse{Status: http.StatusCreated, UserId: user.Id}, nil
}

func (s *Server) mustEmbedUnimplementedAuthServiceServer() {

}
