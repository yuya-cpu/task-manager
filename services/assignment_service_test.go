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

func (m *mockAssignmentRepository) FindByUserID(userID uint, filter dto.AssignmentListFilter) ([]models.Assignment, int64, error) {
	var result []models.Assignment
	for _, a := range m.assignments {
		if a.UserID != userID {
			continue
		}
		if filter.Status != "" && a.Status != filter.Status {
			continue
		}
		if filter.Priority != "" && a.Priority != filter.Priority {
			continue
		}
		result = append(result, a)
	}
	total := int64(len(result))
	return result, total, nil
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

func TestAssignmentService_FindAllFilter(t *testing.T) {
	repo := &mockAssignmentRepository{
		assignments: []models.Assignment{
			{ID: 1, UserID: 1, Title: "A", Status: models.StatusTodo, Priority: models.PriorityHigh},
			{ID: 2, UserID: 1, Title: "B", Status: models.StatusDone, Priority: models.PriorityLow},
		},
	}
	service := NewAssignmentService(repo)

	filter := dto.DefaultAssignmentListFilter()
	filter.Status = models.StatusTodo

	result, err := service.FindAll(1, filter)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.Meta.Total)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, "A", result.Items[0].Title)
}

func TestAssignmentListFilter_Validate(t *testing.T) {
	filter := dto.DefaultAssignmentListFilter()
	filter.Status = "invalid"
	assert.Error(t, filter.Validate())
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
