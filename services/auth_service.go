package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"task-manager/models"
	"task-manager/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignUp(email, password string) error
	Login(email, password string) (string, error)
	GetUserFromToken(tokenString string) (models.User, error)
}

type authService struct {
	repository repositories.AuthRepository
}

func NewAuthService(repository repositories.AuthRepository) AuthService {
	return &authService{repository: repository}
}

func (s *authService) SignUp(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
	}
	return s.repository.CreateUser(&user)
}

func (s *authService) Login(email, password string) (string, error) {
	foundUser, err := s.repository.FindUserByEmail(email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	return CreateToken(foundUser.ID, foundUser.Email)
}

func CreateToken(userID uint, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey()))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *authService) GetUserFromToken(tokenString string) (models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey()), nil
	})
	if err != nil {
		return models.User{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return models.User{}, fmt.Errorf("invalid token")
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return models.User{}, jwt.ErrTokenExpired
	}

	user, err := s.repository.FindUserByEmail(claims["email"].(string))
	if err != nil {
		return models.User{}, err
	}
	return *user, nil
}

func secretKey() string {
	key := os.Getenv("SECRET_KEY")
	if key == "" {
		return "dev-secret-key-change-in-production"
	}
	return key
}
