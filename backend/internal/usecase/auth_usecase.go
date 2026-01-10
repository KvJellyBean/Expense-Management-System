package usecase

import (
	"context"
	"errors"
	"expense-management-system/internal/domain"
	"expense-management-system/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo domain.UserRepository
	cfg      *config.Config
}

func NewAuthUsecase(userRepo domain.UserRepository, cfg *config.Config) domain.AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := u.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (u *authUsecase) ValidateToken(ctx context.Context, tokenString string) (*domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(u.cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user_id in token")
	}

	user, err := u.userRepo.GetByID(ctx, int(userID))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *authUsecase) generateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * time.Duration(u.cfg.JWTExpiryHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.cfg.JWTSecret))
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
