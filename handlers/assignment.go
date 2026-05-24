package handlers

import (
	"net/http"
	"strconv"

	"task-manager/dto"
	"task-manager/services"

	"github.com/gin-gonic/gin"
)

type AssignmentHandler interface {
	FindAll(c *gin.Context)
	FindByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type assignmentHandler struct {
	service services.AssignmentService
}

func NewAssignmentHandler(service services.AssignmentService) AssignmentHandler {
	return &assignmentHandler{service: service}
}

func (h *assignmentHandler) FindAll(c *gin.Context) {
	assignments, err := h.service.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": assignments})
}

func (h *assignmentHandler) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignment id"})
		return
	}

	assignment, err := h.service.FindByID(uint(id))
	if err != nil {
		if err.Error() == "assignment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": assignment})
}

func (h *assignmentHandler) Create(c *gin.Context) {
	var input dto.CreateAssignmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assignment, err := h.service.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": assignment})
}

func (h *assignmentHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignment id"})
		return
	}

	var input dto.UpdateAssignmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assignment, err := h.service.Update(uint(id), input)
	if err != nil {
		if err.Error() == "assignment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": assignment})
}

func (h *assignmentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignment id"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if err.Error() == "assignment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}
	c.Status(http.StatusNoContent)
}
