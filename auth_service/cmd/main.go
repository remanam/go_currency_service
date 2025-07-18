package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v5"
	"github.com/remanam/go_currency_service/auth_service/internal/handler"
	"github.com/remanam/go_currency_service/auth_service/internal/repository"
	"github.com/remanam/go_currency_service/auth_service/internal/service"
	pb "github.com/remanam/go_currency_service/auth_service/pb/auth_service/proto"
	"google.golang.org/grpc"
)

func main() {
	// 1. Подключение к Postgres
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:root@localhost:5432/auth_service?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer conn.Close(context.Background())

	// 2. Создание репозитория и сервиса через интерфейс
	userRepo := repository.NewPostgresUserRepo(conn) // реализует domain.UserRepository
	authService := service.NewAuthService(userRepo)

	// 3. Создание gRPC сервера
	grpcServer := grpc.NewServer()

	// 4. Создание handler-а (AuthHandler реализует pb.AuthServiceServer)
	authHandler := handler.NewAuthHandler(*authService)

	// 5. Регистрация gRPC сервиса
	pb.RegisterAuthServiceServer(grpcServer, authHandler)

	// 6. Запуск сервера
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("AuthService gRPC server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
