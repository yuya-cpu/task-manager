package services

import (
	"time"

	"task-manager/dto"
	"task-manager/models"
	"task-manager/repositories"
)

type AssignmentService interface {
	FindAll() ([]models.Assignment, error)
	FindByID(id uint) (models.Assignment, error)
	Create(input dto.CreateAssignmentInput) (models.Assignment, error)
	Update(id uint, input dto.UpdateAssignmentInput) (models.Assignment, error)
	Delete(id uint) error
}

type assignmentService struct {
	repository repositories.AssignmentRepository
}

func NewAssignmentService(repository repositories.AssignmentRepository) AssignmentService {
	return &assignmentService{repository: repository}
}

func (s *assignmentService) FindAll() ([]models.Assignment, error) {
	return s.repository.FindAll()
}

func (s *assignmentService) FindByID(id uint) (models.Assignment, error) {
	return s.repository.FindByID(id)
}

func (s *assignmentService) Create(input dto.CreateAssignmentInput) (models.Assignment, error) {
	dueDate, err := parseDueDate(input.DueDate)
	if err != nil {
		return models.Assignment{}, err
	}

	status := input.Status
	if status == "" {
		status = models.StatusTodo
	}

	assignment := models.Assignment{
		Title:       input.Title,
		Description: input.Description,
		DueDate:     dueDate,
		Priority:    input.Priority,
		Status:      status,
	}
	return s.repository.Create(assignment)
}

func (s *assignmentService) Update(id uint, input dto.UpdateAssignmentInput) (models.Assignment, error) {
	assignment, err := s.repository.FindByID(id)
	if err != nil {
		return models.Assignment{}, err
	}

	if input.Title != nil {
		assignment.Title = *input.Title
	}
	if input.Description != nil {
		assignment.Description = *input.Description
	}
	if input.DueDate != nil {
		dueDate, err := parseDueDate(input.DueDate)
		if err != nil {
			return models.Assignment{}, err
		}
		assignment.DueDate = dueDate
	}
	if input.Priority != nil {
		assignment.Priority = *input.Priority
	}
	if input.Status != nil {
		assignment.Status = *input.Status
	}

	return s.repository.Update(assignment)
}

func (s *assignmentService) Delete(id uint) error {
	return s.repository.Delete(id)
}

func parseDueDate(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", *value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
