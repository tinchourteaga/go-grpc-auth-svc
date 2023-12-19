package authentication

import (
	"context"
	"net/http"

	"github.com/tinchourteaga/go-grpc-auth-svc/internal/models"
	"github.com/tinchourteaga/go-grpc-auth-svc/internal/pb"
)

type Authentication struct {
	authSvc Service
}

func NewAuthentication(svc Service) *Authentication {
	return &Authentication{
		authSvc: svc,
	}
}

func (a *Authentication) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user := models.User{Email: req.Email, Password: req.Password}
	err := a.authSvc.Register(ctx, user)
	if err != nil {
		return &pb.RegisterResponse{
			Status: http.StatusConflict,
			Error:  err.Error(),
		}, nil
	}

	return &pb.RegisterResponse{
		Status: http.StatusCreated,
	}, nil
}

func (a *Authentication) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user := models.User{Email: req.Email, Password: req.Password}
	token, err := a.authSvc.Login(ctx, user)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  err.Error(),
		}, nil
	}

	return &pb.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

func (a *Authentication) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	userId, err := a.authSvc.Validate(ctx, req.Token)
	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  err.Error(),
		}, nil
	}

	return &pb.ValidateResponse{
		Status: http.StatusOK,
		UserId: userId,
	}, nil
}
