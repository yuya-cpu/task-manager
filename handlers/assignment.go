package handlers

import (
	"net/http"
	"strconv"

	"task-manager/dto"
	"task-manager/models"
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

func (h *assignmentHandler) currentUserID(c *gin.Context) (uint, bool) {
	user, exists := c.Get("user")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return 0, false
	}
	return user.(models.User).ID, true
}

func parseAssignmentListFilter(c *gin.Context) (dto.AssignmentListFilter, error) {
	filter := dto.DefaultAssignmentListFilter()
	filter.Status = c.Query("status")
	filter.Priority = c.Query("priority")
	if sort := c.Query("sort"); sort != "" {
		filter.Sort = sort
	}
	if page := c.Query("page"); page != "" {
		p, err := strconv.Atoi(page)
		if err != nil {
			return filter, err
		}
		filter.Page = p
	}
	if limit := c.Query("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return filter, err
		}
		filter.Limit = l
	}
	return filter, filter.Validate()
}

func (h *assignmentHandler) FindAll(c *gin.Context) {
	userID, ok := h.currentUserID(c)
	if !ok {
		return
	}

	filter, err := parseAssignmentListFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.FindAll(userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result.Items, "meta": result.Meta})
}

func (h *assignmentHandler) FindByID(c *gin.Context) {
	userID, ok := h.currentUserID(c)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignment id"})
		return
	}

	assignment, err := h.service.FindByID(uint(id), userID)
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
	userID, ok := h.currentUserID(c)
	if !ok {
		return
	}

	var input dto.CreateAssignmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assignment, err := h.service.Create(userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": assignment})
}

func (h *assignmentHandler) Update(c *gin.Context) {
	userID, ok := h.currentUserID(c)
	if !ok {
		return
	}

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

	assignment, err := h.service.Update(uint(id), userID, input)
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
	userID, ok := h.currentUserID(c)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignment id"})
		return
	}

	if err := h.service.Delete(uint(id), userID); err != nil {
		if err.Error() == "assignment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}
	c.Status(http.StatusNoContent)
}
