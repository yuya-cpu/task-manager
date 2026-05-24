package services

import (
	"errors"
	"testing"
	"time"

	"task-manager/dto"
	"task-manager/models"

	"github.com/stretchr/testify/assert"
)

type mockAssignmentRepository struct {
	assignments []models.Assignment
}

func (m *mockAssignmentRepository) FindAllByUserID(userID uint) ([]models.Assignment, error) {
	var result []models.Assignment
	for _, a := range m.assignments {
		if a.UserID == userID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (m *mockAssignmentRepository) FindByIDForUser(id, userID uint) (models.Assignment, error) {
	for _, a := range m.assignments {
		if a.ID == id && a.UserID == userID {
			return a, nil
		}
	}
	return models.Assignment{}, errors.New("assignment not found")
}

func (m *mockAssignmentRepository) Create(assignment models.Assignment) (models.Assignment, error) {
	assignment.ID = uint(len(m.assignments) + 1)
	m.assignments = append(m.assignments, assignment)
	return assignment, nil
}

func (m *mockAssignmentRepository) Update(assignment models.Assignment) (models.Assignment, error) {
	for i, a := range m.assignments {
		if a.ID == assignment.ID {
			m.assignments[i] = assignment
			return assignment, nil
		}
	}
	return models.Assignment{}, errors.New("assignment not found")
}

func (m *mockAssignmentRepository) DeleteForUser(id, userID uint) error {
	for i, a := range m.assignments {
		if a.ID == id && a.UserID == userID {
			m.assignments = append(m.assignments[:i], m.assignments[i+1:]...)
			return nil
		}
	}
	return errors.New("assignment not found")
}

func TestAssignmentService_Create(t *testing.T) {
	repo := &mockAssignmentRepository{}
	service := NewAssignmentService(repo)

	due := "2026-06-01"
	created, err := service.Create(1, dto.CreateAssignmentInput{
		Title:    "テストタスク",
		Priority: models.PriorityHigh,
		DueDate:  &due,
	})

	assert.NoError(t, err)
	assert.Equal(t, uint(1), created.UserID)
	assert.Equal(t, models.StatusTodo, created.Status)
	assert.Equal(t, models.PriorityHigh, created.Priority)
	assert.NotNil(t, created.DueDate)
}

func TestAssignmentService_UpdateStatus(t *testing.T) {
	repo := &mockAssignmentRepository{
		assignments: []models.Assignment{
			{ID: 1, UserID: 1, Title: "既存", Priority: models.PriorityLow, Status: models.StatusTodo},
		},
	}
	service := NewAssignmentService(repo)

	done := models.StatusDone
	updated, err := service.Update(1, 1, dto.UpdateAssignmentInput{Status: &done})

	assert.NoError(t, err)
	assert.Equal(t, models.StatusDone, updated.Status)
}

func TestParseDueDate(t *testing.T) {
	empty := ""
	parsed, err := parseDueDate(&empty)
	assert.NoError(t, err)
	assert.Nil(t, parsed)

	valid := "2026-12-31"
	parsed, err = parseDueDate(&valid)
	assert.NoError(t, err)
	assert.Equal(t, 2026, parsed.Year())
	assert.Equal(t, time.December, parsed.Month())
}
