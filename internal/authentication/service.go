package authentication

import (
	"context"
	"errors"
	"strconv"

	"github.com/tinchourteaga/go-grpc-auth-svc/internal/models"
	"github.com/tinchourteaga/go-grpc-auth-svc/pkg/utils"
)

const (
	ErrInvalidCredentials = "invalid credentials"
	ErrInvalidToken       = "invalid token"
)

type Service interface {
	Register(context.Context, models.User) error
	Login(context.Context, models.User) (string, error)
	Validate(context.Context, string) (string, error)
}

type service struct {
	repository Repository
	jwt        utils.JwtWrapper
}

func NewService(repo Repository, jwt utils.JwtWrapper) Service {
	return &service{
		repository: repo,
		jwt:        jwt,
	}
}

func (s *service) Register(ctx context.Context, user models.User) error {
	user.Email = utils.HashPassword(user.Email)
	err := s.repository.Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Login(ctx context.Context, user models.User) (string, error) {
	dbUser, err := s.repository.GetByEmail(ctx, user.Email)
	if err != nil {
		return "", errors.New(ErrInvalidCredentials)
	}

	match := utils.PasswordMatchesHash(user.Password, dbUser.Password)

	if !match {
		return "", errors.New(ErrInvalidCredentials)
	}

	token, _ := s.jwt.GenerateToken(dbUser)

	return token, nil
}

func (s *service) Validate(ctx context.Context, token string) (string, error) {
	claims, err := s.jwt.ValidateToken(token)
	if err != nil {
		return "", errors.New(ErrInvalidToken)
	}

	dbUser, err := s.repository.GetByEmail(ctx, claims.Email)
	if err != nil {
		return "", errors.New(ErrInvalidToken)
	}

	return strconv.Itoa(int(dbUser.Id)), nil
}
