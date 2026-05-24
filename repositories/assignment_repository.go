package repositories

import (
	"errors"

	"task-manager/dto"
	"task-manager/models"

	"gorm.io/gorm"
)

type AssignmentRepository interface {
	FindByUserID(userID uint, filter dto.AssignmentListFilter) ([]models.Assignment, int64, error)
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

func (r *assignmentRepository) FindByUserID(userID uint, filter dto.AssignmentListFilter) ([]models.Assignment, int64, error) {
	query := r.db.Model(&models.Assignment{}).Where("user_id = ?", userID)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sort := filter.Sort
	if sort == "" {
		sort = dto.SortDueDateAsc
	}

	switch sort {
	case dto.SortDueDateDesc:
		query = query.Order("CASE WHEN due_date IS NULL THEN 1 ELSE 0 END, due_date DESC, id ASC")
	case dto.SortNewest:
		query = query.Order("created_at DESC, id DESC")
	default:
		query = query.Order("CASE WHEN due_date IS NULL THEN 1 ELSE 0 END, due_date ASC, id ASC")
	}

	offset := (filter.Page - 1) * filter.Limit
	var assignments []models.Assignment
	if err := query.Offset(offset).Limit(filter.Limit).Find(&assignments).Error; err != nil {
		return nil, 0, err
	}

	return assignments, total, nil
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
