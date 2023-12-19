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

func (a *Authentication) Register(ctx context.Context, req *pb.RegisterRequest) *pb.RegisterResponse {
	user := models.User{Email: req.Email, Password: req.Password}
	err := a.authSvc.Register(ctx, user)
	if err != nil {
		return &pb.RegisterResponse{
			Status: http.StatusConflict,
			Error:  err.Error(),
		}
	}

	return &pb.RegisterResponse{
		Status: http.StatusCreated,
	}
}

func (a *Authentication) Login(ctx context.Context, req *pb.LoginRequest) *pb.LoginResponse {
	user := models.User{Email: req.Email, Password: req.Password}
	token, err := a.authSvc.Login(ctx, user)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusConflict,
			Error:  err.Error(),
		}
	}

	return &pb.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}
}
