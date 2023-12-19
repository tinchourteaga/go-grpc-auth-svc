package main

import (
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/tinchourteaga/go-grpc-auth-svc/internal/authentication"
	"github.com/tinchourteaga/go-grpc-auth-svc/internal/pb"
	"github.com/tinchourteaga/go-grpc-auth-svc/pkg/config"
	"github.com/tinchourteaga/go-grpc-auth-svc/pkg/db"
	"github.com/tinchourteaga/go-grpc-auth-svc/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Error().Msg("config loading: " + err.Error())
	}

	listener, err := net.Listen("tcp", viper.GetString("PORT"))
	if err != nil {
		log.Fatal().Msg("port listening failed: " + err.Error())
	}

	jwt := utils.JwtWrapper{
		SecretKey:       viper.GetString("JWT_SECRET_KEY"),
		Issuer:          "go-grpc-auth-svc",
		ExpirationHours: 24 * 365,
	}

	grpcServer := grpc.NewServer()
	connector := db.NewDatabaseConnection()
	repo := authentication.NewRepository(connector)
	service := authentication.NewService(repo, jwt)
	handler := authentication.NewAuthentication(service)

	pb.RegisterAuthServiceServer(grpcServer, handler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal().Msg("failed to serve: " + err.Error())
	}

	fmt.Println("Auth service listening on: " + viper.GetString("PORT"))
}
