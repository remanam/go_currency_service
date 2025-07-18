package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/remanam/go_currency_service/auth_service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo domain.UserRepository
}

func NewAuthService(userRepo domain.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (a *AuthService) Register(ctx context.Context, username string, email string, password string) (int32, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	//----------------- Проверяем, что пользователя с таким логином и почтой нет в базе
	_, err = a.userRepo.GetByEmail(email)
	if err != nil {
		return 0, err
	}
	if err == pgx.ErrNoRows {
		_, err = a.userRepo.GetByUsername(username)
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("пользователь с таким username или почтой уже существует")
		}
	}
	//------------------------------------------

	user := &domain.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}
	user_id, err := a.userRepo.Create(user)
	if err != nil {
		return 0, err
	}

	return user_id, nil

}

// Claims для JWT токена
type Claims struct {
	UserID   int32
	Username string
	jwt.RegisteredClaims
}

func (a *AuthService) Login(ctx context.Context, username string, password string) (string, error) {
	user, err := a.userRepo.GetByUsername(username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	// Создаем только access token с более длительным сроком
	accessToken, err := createAccessToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

const jwtSecret = "3344"

func createAccessToken(userId int32, username string) (string, error) {
	claims := &Claims{
		UserID:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
