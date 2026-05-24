package dto

import (
	"fmt"

	"task-manager/models"
)

const (
	SortDueDateAsc  = "due_date_asc"
	SortDueDateDesc = "due_date_desc"
	SortNewest      = "newest"
)

type AssignmentListFilter struct {
	Status   string
	Priority string
	Sort     string
	Page     int
	Limit    int
}

type ListMeta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

func DefaultAssignmentListFilter() AssignmentListFilter {
	return AssignmentListFilter{
		Sort:  SortDueDateAsc,
		Page:  1,
		Limit: 20,
	}
}

func (f AssignmentListFilter) Validate() error {
	if f.Status != "" && !models.ValidStatus(f.Status) {
		return fmt.Errorf("invalid status")
	}
	if f.Priority != "" && !models.ValidPriority(f.Priority) {
		return fmt.Errorf("invalid priority")
	}
	switch f.Sort {
	case "", SortDueDateAsc, SortDueDateDesc, SortNewest:
	default:
		return fmt.Errorf("invalid sort")
	}
	if f.Page < 1 {
		return fmt.Errorf("page must be >= 1")
	}
	if f.Limit < 1 || f.Limit > 100 {
		return fmt.Errorf("limit must be between 1 and 100")
	}
	return nil
}
