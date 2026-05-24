package repositories

import (
	"errors"

	"task-manager/models"

	"gorm.io/gorm"
)

type AssignmentRepository interface {
	FindAllByUserID(userID uint) ([]models.Assignment, error)
	FindByIDForUser(id, userID uint) (models.Assignment, error)
	Create(assignment models.Assignment) (models.Assignment, error)
	Update(assignment models.Assignment) (models.Assignment, error)
	DeleteForUser(id, userID uint) error
}

type assignmentRepository struct {
	db *gorm.DB
}

func NewAssignmentRepository(db *gorm.DB) AssignmentRepository {
	return &assignmentRepository{db: db}
}

func (r *assignmentRepository) FindAllByUserID(userID uint) ([]models.Assignment, error) {
	var assignments []models.Assignment
	if err := r.db.Where("user_id = ?", userID).
		Order("CASE WHEN due_date IS NULL THEN 1 ELSE 0 END, due_date ASC, id ASC").
		Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}

func (r *assignmentRepository) FindByIDForUser(id, userID uint) (models.Assignment, error) {
	var assignment models.Assignment
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&assignment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Assignment{}, errors.New("assignment not found")
		}
		return models.Assignment{}, err
	}
	return assignment, nil
}

func (r *assignmentRepository) Create(assignment models.Assignment) (models.Assignment, error) {
	if err := r.db.Create(&assignment).Error; err != nil {
		return models.Assignment{}, err
	}
	return assignment, nil
}

func (r *assignmentRepository) Update(assignment models.Assignment) (models.Assignment, error) {
	if err := r.db.Save(&assignment).Error; err != nil {
		return models.Assignment{}, err
	}
	return assignment, nil
}

func (r *assignmentRepository) DeleteForUser(id, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Assignment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("assignment not found")
	}
	return nil
}
