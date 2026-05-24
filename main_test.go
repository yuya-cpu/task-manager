package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"task-manager/data"
	"task-manager/dto"
	"task-manager/models"
	"task-manager/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	_ = os.Setenv("SECRET_KEY", "test-secret-key")
	os.Exit(m.Run())
}

func setupTestRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	db := data.SetupTestDB()

	hashed, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	assert.NoError(t, err)
	assert.NoError(t, db.Create(&models.User{Email: "test@example.com", Password: string(hashed)}).Error)

	return setupRouter(db), db
}

func TestSignupAndLogin(t *testing.T) {
	router, _ := setupTestRouter(t)

	signupBody, _ := json.Marshal(dto.SignupInput{
		Email:    "new@example.com",
		Password: "password123",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(signupBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	loginBody, _ := json.Marshal(dto.LoginInput{
		Email:    "new@example.com",
		Password: "password123",
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var loginResp map[string]string
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &loginResp))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, loginResp["token"])
}

func TestAssignmentsCRUD(t *testing.T) {
	router, db := setupTestRouter(t)

	token, err := services.CreateToken(1, "test@example.com")
	assert.NoError(t, err)

	createBody, _ := json.Marshal(dto.CreateAssignmentInput{
		Title:    "APIテスト",
		Priority: models.PriorityMedium,
		Status:   models.StatusTodo,
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/assignments", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createResp map[string]models.Assignment
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &createResp))
	assert.Equal(t, "APIテスト", createResp["data"].Title)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/assignments", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var listResp map[string][]models.Assignment
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &listResp))
	assert.Len(t, listResp["data"], 1)

	done := models.StatusDone
	updateBody, _ := json.Marshal(dto.UpdateAssignmentInput{Status: &done})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/assignments/1", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, "/assignments/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	var count int64
	db.Model(&models.Assignment{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestAssignmentsUnauthorized(t *testing.T) {
	router, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/assignments", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
