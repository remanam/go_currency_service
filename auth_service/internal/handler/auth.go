package handler

import (
	"context"
	"fmt"

	"github.com/remanam/go_currency_service/auth_service/internal/service"
	pb "github.com/remanam/go_currency_service/auth_service/pb/auth_service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Password == "" || req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "логин и пароль не могут быть путсыми")
	}

	accessToken, err := h.service.Login(ctx, req.Username, req.Password)
	if err != nil {
		// Проверяем какая конкретно ошибка у нас
		if err.Error() == "no rows in result set" || err.Error() == "username is empty" {
			return nil, status.Error(codes.NotFound, "Пользователь не найден")
		}
		return &pb.LoginResponse{}, err
	}
	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: "",
	}, nil
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	fmt.Printf("username: '%s', email: '%s', password: '%s'\n", req.Username, req.Email, req.Password)
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "Все поля обязательны")
	}

	user_id, err := h.service.Register(ctx, req.Username, req.Email, req.Password)

	if err != nil {
		return &pb.RegisterResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.RegisterResponse{UserId: user_id}, nil
}
