package authentication

import (
	"context"
	"errors"

	"github.com/tinchourteaga/go-grpc-auth-svc/internal/models"
	"github.com/tinchourteaga/go-grpc-auth-svc/pkg/db"
	"gorm.io/gorm"
)

type Repository interface {
	Create(context.Context, models.User) error
	GetByEmail(ctx context.Context, email string) (models.User, error)
	Exists(ctx context.Context, email string) bool
}

type repository struct {
	con db.Connector
}

func NewRepository(con db.Connector) Repository {
	return &repository{
		con: con,
	}
}

func (r *repository) Create(ctx context.Context, user models.User) error {
	exists := r.Exists(ctx, user.Email)
	if exists {
		return errors.New("email already exists")
	}

	r.con.DB.Create(&user)

	return nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	user := models.User{}
	result := r.con.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return models.User{}, gorm.ErrRecordNotFound
	}

	return user, nil
}

func (r *repository) Exists(ctx context.Context, email string) bool {
	result := r.con.DB.Where("email = ?", email).First(&models.User{})
	return result.Error == nil
}
