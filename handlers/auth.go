package handlers

import (
	"net/http"

	"task-manager/dto"
	"task-manager/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	SignUp(c *gin.Context)
	Login(c *gin.Context)
}

type authHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) AuthHandler {
	return &authHandler{service: service}
}

func (h *authHandler) SignUp(c *gin.Context) {
	var input dto.SignupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.SignUp(input.Email, input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *authHandler) Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.service.Login(input.Email, input.Password)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
