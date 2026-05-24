package services

import (
	"errors"
	"testing"

	"task-manager/models"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockAuthRepository struct {
	users []models.User
}

func (m *mockAuthRepository) CreateUser(user *models.User) error {
	user.ID = uint(len(m.users) + 1)
	m.users = append(m.users, *user)
	return nil
}

func (m *mockAuthRepository) FindUserByEmail(email string) (*models.User, error) {
	for i := range m.users {
		if m.users[i].Email == email {
			return &m.users[i], nil
		}
	}
	return nil, errors.New("user not found")
}

func TestAuthService_SignUpAndLogin(t *testing.T) {
	t.Setenv("SECRET_KEY", "test-secret")

	repo := &mockAuthRepository{}
	service := NewAuthService(repo)

	err := service.SignUp("test@example.com", "password123")
	assert.NoError(t, err)

	token, err := service.Login("test@example.com", "password123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	user, err := service.GetUserFromToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestAuthService_LoginInvalidPassword(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	repo := &mockAuthRepository{
		users: []models.User{{Email: "test@example.com", Password: string(hashed)}},
	}
	repo.users[0].ID = 1

	service := NewAuthService(repo)
	_, err := service.Login("test@example.com", "wrongpassword")
	assert.Error(t, err)
}
